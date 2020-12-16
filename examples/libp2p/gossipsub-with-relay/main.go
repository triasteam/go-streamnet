package main

import (
	"context"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/libp2p/go-libp2p"
	circuit "github.com/libp2p/go-libp2p-circuit"
	"github.com/libp2p/go-libp2p-core/crypto"
	event "github.com/libp2p/go-libp2p-core/event"
	"github.com/libp2p/go-libp2p-core/host"
	network "github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/routing"
	kaddht "github.com/libp2p/go-libp2p-kad-dht"
	mplex "github.com/libp2p/go-libp2p-mplex"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	secio "github.com/libp2p/go-libp2p-secio"
	yamux "github.com/libp2p/go-libp2p-yamux"
	"github.com/libp2p/go-libp2p/config"
	"github.com/libp2p/go-libp2p/p2p/discovery"
	tcp "github.com/libp2p/go-tcp-transport"
	ws "github.com/libp2p/go-ws-transport"
	"github.com/multiformats/go-multiaddr"
)

const (
	privateName = "priv.pem"
)

type mdnsNotifee struct {
	h   host.Host
	ctx context.Context
}

func (m *mdnsNotifee) HandlePeerFound(pi peer.AddrInfo) {
	m.h.Connect(m.ctx, pi)
}

func loadFromPem() (crypto.PrivKey, error) {
	pwd, _ := os.Getwd()
	filePath := fmt.Sprintf("%s/%s", pwd, privateName)
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(bytes)
	priv, _ := crypto.UnmarshalPrivateKey(block.Bytes)
	return priv, nil
}

func exportToPem(priv crypto.PrivKey) error {
	privBytes, err := crypto.MarshalPrivateKey(priv)
	if err != nil {
		fmt.Printf("marshal private key error: %s \n", err)
		return err
	}
	privPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "Ed25519 PRIVATE KEY",
			Bytes: privBytes,
		},
	)
	fmt.Println("private pem: ", privPem)
	ioutil.WriteFile(privateName, []byte(privPem), 0644)
	return nil
}

func getOrGeneratePrivateKey() crypto.PrivKey {
	priv, err := loadFromPem()
	if err != nil {
		priv, _, err = crypto.GenerateKeyPair(crypto.Ed25519, -1)
		if err != nil {
			panic(err)
		}
		err = exportToPem(priv)
		if err != nil {
			fmt.Printf("export private key error: %s \n", err)
		}
	}
	return priv
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, _ := ParseFlags()

	transports := libp2p.ChainOptions(
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.Transport(ws.New),
	)

	muxers := libp2p.ChainOptions(
		libp2p.Muxer("/yamux/1.0.0", yamux.DefaultTransport),
		libp2p.Muxer("/mplex/6.7.0", mplex.DefaultTransport),
	)

	security := libp2p.Security(secio.ID, secio.New)

	listenAddrs := libp2p.ListenAddrStrings(
		fmt.Sprintf("/ip4/0.0.0.0/tcp/%s", cfg.Port),
		fmt.Sprintf("/ip4/0.0.0.0/tcp/%s/ws", cfg.Port),
	)

	var dht *kaddht.IpfsDHT
	newDHT := func(h host.Host) (routing.PeerRouting, error) {
		var err error
		dht, err = kaddht.New(ctx, h, kaddht.Mode(kaddht.ModeServer))
		return dht, err
	}
	routing := libp2p.Routing(newDHT)

	priv := getOrGeneratePrivateKey()

	relayOption := func() config.Option {

		if cfg.RelayType == "hop" {
			return libp2p.ChainOptions(libp2p.EnableAutoRelay(), libp2p.EnableRelay(circuit.OptHop), libp2p.AddrsFactory(func(addrs []multiaddr.Multiaddr) []multiaddr.Multiaddr {
				for i, addr0 := range addrs {
					saddr := addr0.String()
					if strings.HasPrefix(saddr, "/ip4/127.0.0.1") {
						addrNoIP := strings.TrimPrefix(saddr, "/ip4/127.0.0.1")
						fmt.Printf("result : %d, public: %s \n", len(cfg.PublicAddr), cfg.PublicAddr)
						if len(cfg.PublicAddr) == 0 {
							addrs[i] = multiaddr.StringCast("/dns4/localhost" + addrNoIP)
						} else {
							addrs[i] = multiaddr.StringCast(fmt.Sprintf("/ip4/%s", cfg.PublicAddr) + addrNoIP)
						}
					}
				}
				return addrs
			}))
		} else if cfg.RelayType == "autorelay" {
			return libp2p.ChainOptions(libp2p.EnableAutoRelay())
		}
		return func(cfg *config.Config) error { return nil }

	}

	host, err := libp2p.New(
		ctx,
		transports,
		listenAddrs,
		muxers,
		security,
		routing,
		libp2p.Identity(priv),
		relayOption(),
	)
	if err != nil {
		panic(err)
	}
	fmt.Printf("my peer is /ip4/127.0.0.1/tcp/%s/ipfs/%s \n", cfg.Port, host.ID().Pretty())

	// 触发搜索autorelay
	if cfg.RelayType == "autorelay" {
		go func() {
			ticker := time.NewTicker(time.Second * 5)
			for {
				select {
				case <-ticker.C:
					privEmitter, _ := host.EventBus().Emitter(new(event.EvtLocalReachabilityChanged))
					privEmitter.Emit(event.EvtLocalReachabilityChanged{Reachability: network.ReachabilityPrivate})
				}
			}
		}()
	}

	ps, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		panic(err)
	}
	sub, err := ps.Subscribe(pubsubTopic)
	if err != nil {
		panic(err)
	}
	go pubsubHandler(ctx, sub)

	for _, addr := range host.Addrs() {
		fmt.Println("Listening on", addr)
	}

	targetAddr, err := multiaddr.NewMultiaddr(cfg.Seed)
	if err != nil {
		panic(err)
	}

	targetInfo, err := peer.AddrInfoFromP2pAddr(targetAddr)
	if err != nil {
		panic(err)
	}

	err = host.Connect(ctx, *targetInfo)
	if err != nil {
		fmt.Printf("connect to host error: %s \n", err)
	} else {
		fmt.Println("Connected to", targetInfo.ID)
	}

	// TEST: 每隔10秒钟打印一次对等方的地址
	go func() {
		var printer = func() {
			if len(host.Peerstore().Peers()) < 1 {
				fmt.Println("i have no peer.")
			}
			for _, p := range host.Peerstore().Peers() {
				addrs := host.Peerstore().Addrs(p)
				fmt.Printf("i have peer[%s], it's addrs is: %s \n", p.Pretty(), addrs)
			}
		}
		ticker := time.NewTicker(time.Second * 10)
		for {
			select {
			case <-ticker.C:
				printer()
			}
		}
	}()

	mdns, err := discovery.NewMdnsService(ctx, host, time.Second*10, "")
	if err != nil {
		panic(err)
	}
	mdns.RegisterNotifee(&mdnsNotifee{h: host, ctx: ctx})

	err = dht.Bootstrap(ctx)
	if err != nil {
		panic(err)
	}

	donec := make(chan struct{}, 1)
	go chatInputLoop(ctx, host, ps, donec)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT)

	select {
	case <-stop:
		host.Close()
		os.Exit(0)
	case <-donec:
		host.Close()
	}
}
