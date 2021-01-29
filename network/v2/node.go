// Copyright 2017 The GoReporter Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package networkv2 is an upgraded version of the networkv1, and provides basic
// network layer components.
// Node is responsible for building a computer node with complete network functions.
package networkv2

import (
	"context"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/libp2p/go-libp2p"
	circuit "github.com/libp2p/go-libp2p-circuit"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/event"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
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
	localConfig "github.com/triasteam/go-streamnet/config"
)

const (
	privateName = "priv.pem"
)

type mdnsNotifee struct {
	h   host.Host
	ctx context.Context
}

// Node respect a peer
type Node struct {
	SendMessageChan chan []byte
	Receive         func(msg []byte) error
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
		log.Printf("marshal private key error: %s \n", err)
		return err
	}
	privPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "Ed25519 PRIVATE KEY",
			Bytes: privBytes,
		},
	)
	log.Println("private pem: ", privPem)
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
			log.Printf("export private key error: %s \n", err)
		}
	}
	return priv
}

// NewNode build a network node, you can add autorelay, gossip to it
func NewNode(ctx context.Context, cfg *localConfig.Config, receive func(msg []byte) error) (*Node, error) {
	node := &Node{
		Receive:         receive,
		SendMessageChan: make(chan []byte),
	}

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
		dht, err = kaddht.New(ctx, h, kaddht.Mode(kaddht.ModeAutoServer))
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
						log.Printf("result : %d, public: %s \n", len(cfg.PublicAddr), cfg.PublicAddr)
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
	log.Printf("my peer is /ip4/127.0.0.1/tcp/%s/ipfs/%s \n", cfg.Port, host.ID().Pretty())

	// 触发搜索autorelay
	if cfg.RelayType == "autorelay" {
		go func() {
			ticker := time.NewTicker(time.Second * 5)

			for {
				relayHop := make(map[string]struct{})
				for _, p := range host.Peerstore().Peers() {
					addrs := host.Peerstore().Addrs(p)
					for _, addr := range addrs {
						if match, _ := regexp.Match("p2p-circuit", []byte(addr.String())); match {
							peerID := ParseRelayPeerID(addr)
							relayHop[peerID] = struct{}{}
						}
					}
				}
				if len(relayHop) > 0 {
					log.Printf("Find Relay %v!! \n", relayHop)
				}
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
	psh := &PubSubHandler{
		node: node,
	}
	go psh.pubsubHandler(ctx, sub)

	for _, addr := range host.Addrs() {
		log.Println("Listening on", addr)
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
		log.Printf("connect to host error: %s \n", err)
	} else {
		log.Println("Connected to", targetInfo.ID)
	}

	mdns, err := discovery.NewMdnsService(ctx, host, time.Second*10, "")
	if err != nil {
		panic(err)
	}
	mdns.RegisterNotifee(&mdnsNotifee{h: host, ctx: ctx})

	err = dht.Bootstrap(ctx)
	if err != nil {
		panic(err)
	}

	protocol := &StreamNetProtocol{
		node: node,
	}

	donec := make(chan struct{}, 1)
	go protocol.chatInputLoop(ctx, host, ps, donec)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT)
	return node, nil
}

// Broadcast store msg to SendMessageChan channel
func (node *Node) Broadcast(msg []byte) {
	node.SendMessageChan <- msg
}
