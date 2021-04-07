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
 * This program demonstrate a application using p2p relay protocol.
 * With relay Protocol, node behind NAT can also join the network.
 *
 * this file describes main of the demo
 */
package main

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/relay/p2pnode"

	circuit "github.com/libp2p/go-libp2p-circuit"
	swarm "github.com/libp2p/go-libp2p-swarm"
	ma "github.com/multiformats/go-multiaddr"
)

func main() {
	// fmt.Println("start three node below:")
	// n2 := p2pnode.NewH2Node()
	// fmt.Println("\n n2 p2p address is ", n2.MultiAddress)

	// n3 := p2pnode.NewH3Node(n2.MultiAddress)

	// fmt.Println("\n  n3 p2p address is ", n3.MultiAddress, ", peer Id is ", n3.PeerID)

	// n1 := p2pnode.NewH1Node(n2.MultiAddress, n3.PeerID)

	// fmt.Println("\n n1 address is ", n1.MultiAddress)
	// select {}
	// mainDeprecated()

	cfg := parseFlags()
	if cfg.n == 1 {
		pid, err := peer.IDB58Decode(cfg.peerID)
		if err != nil {
			panic(err)
		}

		n1 := p2pnode.NewH1Node(cfg.address, pid)
		fmt.Println("\n n1 address is ", n1.MultiAddress)
	} else if cfg.n == 2 {
		n2 := p2pnode.NewH2Node()
		fmt.Println("\n  n2 p2p address is ", n2.MultiAddress)
	} else if cfg.n == 3 {
		n3 := p2pnode.NewH3Node(cfg.address)
		fmt.Println("\n  n3 p2p address is ", n3.MultiAddress)
	} else {
		n2 := p2pnode.NewH2Node()
		fmt.Println("\n  n2 p2p address is ", n2.MultiAddress)
	}
	select {}
}

// main_deprecated ...
func mainDeprecated() {
	// Create three libp2p hosts, enable relay client capabilities on all
	// of them.

	// Tell the host to monitor for relays.
	h1, err := libp2p.New(context.Background(), libp2p.EnableRelay())
	if err != nil {
		panic(err)
	}

	// Tell the host to relay connections for other peers (The ability to *use*
	// a relay vs the ability to *be* a relay)
	h2, err := libp2p.New(context.Background(), libp2p.EnableRelay(circuit.OptHop))
	if err != nil {
		panic(err)
	}

	// Zero out the listen addresses for the host, so it can only communicate
	// via p2p-circuit for our example
	h3, err := libp2p.New(context.Background(), libp2p.ListenAddrs(), libp2p.EnableRelay())
	if err != nil {
		panic(err)
	}

	h2info := peer.AddrInfo{
		ID:    h2.ID(),
		Addrs: h2.Addrs(),
	}

	// Connect both h1 and h3 to h2, but not to each other
	if err := h1.Connect(context.Background(), h2info); err != nil {
		panic(err)
	}
	if err := h3.Connect(context.Background(), h2info); err != nil {
		panic(err)
	}

	// Now, to test things, let's set up a protocol handler on h3
	h3.SetStreamHandler("/cats", func(s network.Stream) {
		fmt.Println("Meow! It worked!")
		s.Close()
	})

	_, err = h1.NewStream(context.Background(), h3.ID(), "/cats")
	if err == nil {
		fmt.Println("Didnt actually expect to get a stream here. What happened?")
		return
	}
	fmt.Println("Okay, no connection from h1 to h3: ", err)
	fmt.Println("Just as we suspected")

	// Creates a relay address
	relayaddr, err := ma.NewMultiaddr("/p2p/" + h2.ID().Pretty() + "/p2p-circuit/p2p/" + h3.ID().Pretty())
	if err != nil {
		panic(err)
	}

	// Since we just tried and failed to dial, the dialer system will, by default
	// prevent us from redialing again so quickly. Since we know what we're doing, we
	// can use this ugly hack (it's on our TODO list to make it a little cleaner)
	// to tell the dialer "no, its okay, let's try this again"
	h1.Network().(*swarm.Swarm).Backoff().Clear(h3.ID())

	h3relayInfo := peer.AddrInfo{
		ID:    h3.ID(),
		Addrs: []ma.Multiaddr{relayaddr},
	}

	fmt.Printf("h3relayInfo is %+v", h3relayInfo)

	// if err := h1.Connect(context.Background(), h3relayInfo); err != nil {
	// 	panic(err)
	// }
	h1.Peerstore().AddAddr(h3relayInfo.ID, h3relayInfo.Addrs[0], peerstore.PermanentAddrTTL)

	// Woohoo! we're connected!
	s, err := h1.NewStream(context.Background(), h3.ID(), "/cats")
	if err != nil {
		fmt.Println("huh, this should have worked: ", err)
		return
	}

	s.Read(make([]byte, 1)) // block until the handler closes the stream
}
