/*
 *
 * The MIT License (MIT)
 *
 * Copyright (c) 2014 Juan Batiz-Benet
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 *
 * This program demonstrate a simple gossip application using p2p.
 * you can simply start a node by executing `main` and added to the network using
 * http portal `addPeers`, and using `send` for communicating.
 *
 */
package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	csms "github.com/libp2p/go-conn-security-multistream"
	bhost "github.com/libp2p/go-libp2p-blankhost"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/metrics"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/libp2p/go-libp2p-peerstore/pstoremem"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	secio "github.com/libp2p/go-libp2p-secio"
	swarm "github.com/libp2p/go-libp2p-swarm"
	tnet "github.com/libp2p/go-libp2p-testing/net"
	tptu "github.com/libp2p/go-libp2p-transport-upgrader"
	yamux "github.com/libp2p/go-libp2p-yamux"
	msmux "github.com/libp2p/go-stream-muxer-multistream"
	"github.com/libp2p/go-tcp-transport"
	"github.com/multiformats/go-multiaddr"
)

// GetHost generate a host
func GetHost(ctx context.Context, seq int) (host.Host, error) {
	// tnet will generate peers with specific order on a machine
	// a seq can help creating a unique peer id
	p := new(tnet.PeerNetParams)
	for i := 0; ; i++ {
		pp, err := tnet.RandPeerNetParams()
		if err != nil {
			panic(err)
		}
		if i == seq {
			p = pp
			break
		}
	}

	ps := pstoremem.NewPeerstore()
	ps.AddPubKey(p.ID, p.PubKey)
	ps.AddPrivKey(p.ID, p.PrivKey)

	sw := swarm.NewSwarm(ctx, p.ID, ps, metrics.NewBandwidthCounter())

	// build upgrader
	id := sw.LocalPeer()
	pk := sw.Peerstore().PrivKey(id)
	secMuxer := new(csms.SSMuxer)
	secMuxer.AddTransport(secio.ID, &secio.Transport{
		LocalID:    id,
		PrivateKey: pk,
	})

	stMuxer := msmux.NewBlankTransport()
	stMuxer.AddTransport("/yamux/1.0.0", yamux.DefaultTransport)

	upgrader := &tptu.Upgrader{
		Secure: secMuxer,
		Muxer:  stMuxer,
	}
	// @unuse
	// upgrader.ConnGater = cfg.ConnGater
	tcpTransport := tcp.NewTCPTransport(upgrader)
	// @unuse
	// tcpTransport.DisableReuseport = cfg.DisableReuseport

	if err := sw.AddTransport(tcpTransport); err != nil {
		panic(err)
	}

	// dail
	if err := sw.Listen(p.Addr); err != nil {
		panic(err)
	}
	sw.Peerstore().AddAddrs(p.ID, sw.ListenAddresses(), peerstore.PermanentAddrTTL)

	h := bhost.NewBlankHost(sw)

	return h, nil
}

func main() {
	// config
	help := flag.Bool("help", false, "Display Help")
	cfg := parseFlags()

	if *help {
		fmt.Printf("Start a gossip peer.")
		fmt.Printf("Usage: \n Run ./main -hp [httpPort] -s [sequence]")
		os.Exit(0)
	}

	ctx := context.Background()

	host, err := GetHost(ctx, cfg.seq)
	if err != nil {
		panic(err)
	}

	//
	ps, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		panic(err)
	}

	// subsciption topic
	subs, err := ps.Subscribe(cfg.Topic)
	if err != nil {
		panic(err)
	}

	var port string
	for _, la := range host.Network().ListenAddresses() {
		if p, err := la.ValueForProtocol(multiaddr.P_TCP); err == nil {
			port = p
			break
		}
	}

	thisAddr := fmt.Sprintf("/ip4/127.0.0.1/tcp/%v/p2p/%s", port, host.ID().Pretty())
	fmt.Println("this multiaddr is : ", thisAddr)

	// curl "http://127.0.0.1:8002/addPeers?dest=/ip4/127.0.0.1/tcp/61071/p2p/12D3KooWE4qDcRrueTuRYWUdQZgcy7APZqBngVeXRt4Y6ytHizKV"
	// curl "http://127.0.0.1:8003/addPeers?dest=/ip4/127.0.0.1/tcp/49683/p2p/12D3KooWB1b3qZxWJanuhtseF3DmPggHCtG36KZ9ixkqHtdKH9fh"
	// curl "http://127.0.0.1:8003/send?msg=hellokatty"
	// 建立http连接
	go func() {
		http.HandleFunc("/addPeers", func(w http.ResponseWriter, req *http.Request) {
			req.ParseForm()
			dests := req.FormValue("dest")
			for _, dest := range strings.Split(dests, ",") {
				// 建立连接，如果http请求包含地址，则创建连接
				maddr, err := multiaddr.NewMultiaddr(dest)
				if err != nil {
					panic(err)
				}

				pi, err := peer.AddrInfoFromP2pAddr(maddr)
				if err != nil {
					panic(err)
				}

				err = host.Connect(ctx, *pi)
				if err != nil {
					panic(err)
				}
			}

			time.Sleep(2 * time.Second)
			fmt.Println("peer connected!")
		})
		http.HandleFunc("/send", func(w http.ResponseWriter, req *http.Request) {
			req.ParseForm()
			s := req.FormValue("msg")
			fmt.Println("http request is : ", s)
			// msg <- s
			ps.Publish(cfg.Topic, []byte(s))
		})
		http.ListenAndServe(fmt.Sprintf(":%v", cfg.httpServerPort), nil)
	}()

	// 接收消息，并阻塞主协程
	for {
		msg1, err := subs.Next(ctx)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(msg1.Data))
	}
}
