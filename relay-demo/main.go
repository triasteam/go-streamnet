package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"

	circuit "github.com/libp2p/go-libp2p-circuit"
)

func handleStream(s network.Stream) {
	log.Println("Got a new stream!")

	// Create a buffer stream for non blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	go readData(rw)
	go writeData(rw)

	// stream 's' will stay open until you close it (or the other side closes it).
}
func readData(rw *bufio.ReadWriter) {
	for {
		str, _ := rw.ReadString('\n')

		if str == "" {
			return
		}
		if str != "\n" {
			// Green console colour: 	\x1b[32m
			// Reset console colour: 	\x1b[0m
			fmt.Printf("\x1b[32m%s\x1b[0m> ", str)
		}

	}
}
func writeData(rw *bufio.ReadWriter) {
	stdReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		sendData, err := stdReader.ReadString('\n')

		if err != nil {
			panic(err)
		}

		rw.WriteString(fmt.Sprintf("%s\n", sendData))
		rw.Flush()
	}

}

func main() {

	choice := flag.Int("c", 1, "relay example part of choice")
	dest := flag.String("d", "", "Destination multiaddr string")

	flag.Parse()

	r := rand.Reader

	// Creates a new RSA key pair for this host.
	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		panic(err)
	}

	// 0.0.0.0 will listen on any interface device.
	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", 0))


	if *choice == 1 {
		// Tell the host to monitor for relays.
		host1, err := libp2p.New(
			context.Background(),
			libp2p.ListenAddrs(sourceMultiAddr),
			libp2p.Identity(prvKey),
			libp2p.EnableRelay(circuit.OptDiscovery),
		)
		if err != nil {
			panic(err)
		}

		fmt.Println("This node's multiaddresses:")
		for _, la := range host1.Addrs() {
			fmt.Printf(" - %v\n", la)
		}
		fmt.Println()

		host1.SetStreamHandler("/chat/1.0.0", handleStream)

		maddr, err := multiaddr.NewMultiaddr(*dest)
		if err != nil {
			log.Fatalln(err)
		}

		// Extract the peer ID from the multiaddr.
		host2, err := peer.AddrInfoFromP2pAddr(maddr)
		if err != nil {
			log.Fatalln(err)
		}

		h2info := peer.AddrInfo{
			ID:    host2.ID,
			Addrs: host2.Addrs,
		}

		if err := host1.Connect(context.Background(), h2info); err != nil {
			panic(err)
		}

		fmt.Printf("Run './main -c 3 -d %s/p2p-circuit/p2p/%s' on another console.\n", *dest, host1.ID().Pretty())

		// Hang forever
		<-make(chan struct{})

	} else if *choice == 2 {
		// Tell the host to relay connections for other peers (The ability to *use*
		// a relay vs the ability to *be* a relay)
		host2, err := libp2p.New(
			context.Background(),
			libp2p.ListenAddrs(sourceMultiAddr),
			libp2p.Identity(prvKey),
			libp2p.EnableRelay(circuit.OptHop),
		)
		if err != nil {
			panic(err)
		}

		fmt.Println("This node's multiaddresses:")
		for _, la := range host2.Addrs() {
			fmt.Printf(" - %v\n", la)
		}
		fmt.Println()

		// Let's get the actual TCP port from our listen multiaddr, in case we're using 0 (default; random available port).
		var port string
		for _, la := range host2.Network().ListenAddresses() {
			if p, err := la.ValueForProtocol(multiaddr.P_TCP); err == nil {
				port = p
				break
			}
		}

		if port == "" {
			panic("was not able to find actual local port")
		}

		fmt.Printf("Run './main -c 1 -d /ip4/{ip}/tcp/%v/p2p/%s' on another console.\n", port, host2.ID().Pretty())
		fmt.Printf("\nI am relay\n\n")

		if *dest != "" {
			maddr, err := multiaddr.NewMultiaddr(*dest)
			if err != nil {
				log.Fatalln(err)
			}

			// Extract the peer ID from the multiaddr.
			host, err := peer.AddrInfoFromP2pAddr(maddr)
			if err != nil {
				log.Fatalln(err)
			}

			hoinfo := peer.AddrInfo{
				ID:    host.ID,
				Addrs: host.Addrs,
			}

			if err := host2.Connect(context.Background(), hoinfo); err != nil {
				panic(err)
			}
			fmt.Printf("\nConnect to relay\n\n")
		}

		// Hang forever
		<-make(chan struct{})
	} else if *choice == 3 {
		// Zero out the listen addresses for the host, so it can only communicate
		// via p2p-circuit for our example
		host3, err := libp2p.New(
			context.Background(),
			libp2p.ListenAddrs(sourceMultiAddr),
			libp2p.Identity(prvKey),
			libp2p.EnableRelay(),
		)
		if err != nil {
			panic(err)
		}

		fmt.Println("This node's multiaddresses:")
		for _, la := range host3.Addrs() {
			fmt.Printf(" - %v\n", la)
		}
		fmt.Println()

		maddr, err := multiaddr.NewMultiaddr(*dest)
		if err != nil {
			log.Fatalln(err)
		}

		// Extract the peer ID from the multiaddr.
		host1, err := peer.AddrInfoFromP2pAddr(maddr)
		if err != nil {
			log.Fatalln(err)
		}

		h1info := peer.AddrInfo{
			ID:    host1.ID,
			Addrs: host1.Addrs,
		}
		if err := host3.Connect(context.Background(), h1info); err != nil {
			panic(err)
		}

		// Start a stream with the destination.
		// Multiaddress of the destination peer is fetched from the peerstore using 'peerId'.
		s, err := host3.NewStream(context.Background(), h1info.ID, "/chat/1.0.0")
		if err != nil {
			panic(err)
		}

		// Create a buffered stream so that read and writes are non blocking.
		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

		// Create a thread to read and write data.
		go writeData(rw)
		go readData(rw)

		// Hang forever.
		select {}
	}

}