package network

import (
	"bufio"
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	swarm "github.com/libp2p/go-libp2p-swarm"
	"github.com/multiformats/go-multiaddr"
)

// Node hold send chan and store func
type Node struct {
	// sendChan contains sendData used for broadcast message
	SendChan chan []byte
	// StoreFunc can be invoked when receiving broadcast message from neigbors
	Receive func(message string) error
}

// Init SendChan and on receive method
func (node *Node) Init(_receive func(data string) error) {
	node.SendChan = make(chan []byte, 1)
	node.Receive = _receive
	node.NewNetwork()
}

// Broadcast message to other node
func (node *Node) Broadcast(data string) bool {
	node.SendChan <- []byte(data)
	// close(node.SendChan)
	return true
}

func (node *Node) handleStream(stream network.Stream) {
	fmt.Println("Got a new stream!")

	// Create a buffer stream for non blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	go node.readData(rw)
	go node.writeData(rw)

	// 'stream' will stay open until you close it (or the other side closes it).
}

func (node *Node) readData(rw *bufio.ReadWriter) {
	for {
		str, err := rw.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from buffer")
			panic(err)
		}

		if str == "" {
			return
		}
		if str != "\n" {
			// Green console colour: 	\x1b[32m
			// Reset console colour: 	\x1b[0m
			fmt.Printf("\x1b[32m%s\x1b[0m> ", str)
			node.Receive(str)
		}

	}
}

func (node *Node) writeData(rw *bufio.ReadWriter) {
	for {
		fmt.Print("waiting input... ")
		sendData := <-node.SendChan
		fmt.Print("send data is ", string(sendData))
		_, err := rw.WriteString(fmt.Sprintf("%s\n", sendData))
		if err != nil {
			fmt.Println("Error writing to buffer")
			panic(err)
		}
		err = rw.Flush()
		if err != nil {
			fmt.Println("Error flushing buffer")
			panic(err)
		}
	}
}

// NewNetwork create a new p2p network server
func (node *Node) NewNetwork() {
	// config
	help := flag.Bool("help", false, "Display Help")
	cfg := parseFlags()

	if *help {
		fmt.Printf("Start a gossip peer.")
		fmt.Printf("Usage: \n Run ./main -sp [port] -d [destination multiaddr string]")
		os.Exit(0)
	}

	var r io.Reader

	r = rand.Reader

	// Creates a new RSA key pair for this host
	privKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		panic(err)
	}

	fmt.Printf("========= sp : %v, relay: %s, d: %s \n", cfg.sp, cfg.relayAddress, cfg.d)

	// 0.0.0.0 will listen on any interface device.
	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", cfg.sp))

	// libp2p.New constructs a new libp2p Host.
	// Other options can be added here.
	host, err := libp2p.New(
		context.Background(),
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(privKey),
		libp2p.EnableRelay(),
	)
	if err != nil {
		panic(err)
	}

	// relay info,
	relayInfo := &peer.AddrInfo{}
	if cfg.relayAddress != "" {
		fmt.Printf("relay address is : %s \n", cfg.relayAddress)
		relayAddress, err := multiaddr.NewMultiaddr(cfg.relayAddress)
		if err != nil {
			panic(err)
		}

		relayInfo, err = peer.AddrInfoFromP2pAddr(relayAddress)
		if err != nil {
			panic(err)
		}
	}

	if cfg.d == "" {

		// Let's get the actual TCP port from our listen multiaddr, in case we're using 0 (default; random available port).
		var port string
		for _, la := range host.Network().ListenAddresses() {
			if p, err := la.ValueForProtocol(multiaddr.P_TCP); err == nil {
				port = p
				break
			}
		}

		if port == "" {
			panic("was not able to find actual local port")
		}

		// If relay exists, then connect it
		if cfg.relayAddress != "" {
			if err = host.Connect(context.Background(), *relayInfo); err != nil {
				panic(err)
			}
		}

		fmt.Printf("Run './main -d /ip4/127.0.0.1/tcp/%v/p2p/%s' on another console.\n", port, host.ID().Pretty())
		fmt.Println("You can replace 127.0.0.1 with public IP as well.")
		fmt.Printf("\nWaiting for incoming connection\n\n")

		// Set a function as stream handler.
		// This function is called when a peer connects, and starts a stream with this protocol.
		// Only applies on the receiving side.
		host.SetStreamHandler("/chat/1.0.0", node.handleStream)
		// Hang forever
		// <-make(chan struct{})
	} else {
		fmt.Println("This node's multiaddresses:")
		for _, la := range host.Addrs() {
			fmt.Printf(" - %v\n", la)
		}
		fmt.Println()

		// Turn the destination into a multiaddr.
		maddr, err := multiaddr.NewMultiaddr(cfg.d)
		if err != nil {
			log.Fatalln(err)
		}

		// Extract the peer ID from the multiaddr.
		info, err := peer.AddrInfoFromP2pAddr(maddr)
		if err != nil {
			log.Fatalln(err)
		}

		// Add the destination's peer multiaddress in the peerstore.
		// This will be used during connection and stream creation by libp2p.
		host.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)

		// Start a stream with the destination.
		// Multiaddress of the destination peer is fetched from the peerstore using 'peerId'.
		s, err := host.NewStream(context.Background(), info.ID, "/chat/1.0.0")
		if err != nil {
			// If panic connect err, then try to use relay
			// First, host must connect to relay successfully
			if err = host.Connect(context.Background(), *relayInfo); err != nil {
				panic(err)
			}

			// Second, build a multiaddress use relay
			destAddrWithRelay, err := multiaddr.NewMultiaddr("/p2p/" + relayInfo.ID.Pretty() + "/p2p-circuit/p2p/" + info.ID.Pretty())
			if err != nil {
				panic(err)
			}

			host.Network().(*swarm.Swarm).Backoff().Clear(relayInfo.ID)

			// Third, Connect to dest addr with relay
			destAddrInfoWithRelay := &peer.AddrInfo{
				ID:    info.ID,
				Addrs: []multiaddr.Multiaddr{destAddrWithRelay},
			}

			if err = host.Connect(context.Background(), *destAddrInfoWithRelay); err != nil {
				panic(err)
			}

			// Forth, try to connect again
			s, err = host.NewStream(context.Background(), info.ID, "/chat/1.0.0")
			if err != nil {
				panic(err)
			}
		}

		// Create a buffered stream so that read and writes are non blocking.
		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

		// Create a thread to read and write data.
		go node.writeData(rw)
		go node.readData(rw)

		// Hang forever.
		// select {}
	}
}
