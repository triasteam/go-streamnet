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
 * this file describes how it works
 */
package p2pnode

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p"
	circuit "github.com/libp2p/go-libp2p-circuit"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	swarm "github.com/libp2p/go-libp2p-swarm"
	"github.com/multiformats/go-multiaddr"
)

// Node ...
type Node struct {
	MultiAddress string
	PeerID       peer.ID
}

// NewH1Node ...
func NewH1Node(h2Dest string, h3Id peer.ID) *Node {
	n := &Node{}
	multiaddrChan := make(chan string)
	peerIDChan := make(chan peer.ID)
	go func() {
		// Tell the host to monitor for relays.
		h1, err := libp2p.New(context.Background(), libp2p.EnableRelay())
		if err != nil {
			panic(err)
		}

		maddr, err := multiaddr.NewMultiaddr(h2Dest)
		if err != nil {
			panic(err)
		}

		h2Info, err := peer.AddrInfoFromP2pAddr(maddr)
		if err != nil {
			panic(err)
		}
		fmt.Println("a h2 peer info was created. ")

		if err = h1.Connect(context.Background(), *h2Info); err != nil {
			panic(err)
		}

		fmt.Println("h3Id.Pretty() ", h3Id.Pretty())
		relayAddr, err := multiaddr.NewMultiaddr("/p2p/" + h2Info.ID.Pretty() + "/p2p-circuit/p2p/" + h3Id.Pretty())
		if err != nil {
			panic(err)
		}
		fmt.Println("a relay addr was created. ")

		h1.Network().(*swarm.Swarm).Backoff().Clear(h2Info.ID)

		h3RelayInfo := peer.AddrInfo{
			ID:    h3Id,
			Addrs: []multiaddr.Multiaddr{relayAddr},
		}

		fmt.Printf("try to connect %+v \n", h3RelayInfo)

		if err = h1.Connect(context.Background(), h3RelayInfo); err != nil {
			fmt.Println("try to connect failed. ", err)
			panic(err)
		}

		fmt.Println("connect")

		s, err := h1.NewStream(context.Background(), h3Id, "/cats1")
		if err != nil {
			panic(err)
		}

		multiaddrChan <- fmt.Sprintf("%+v/p2p/%s", h1.Addrs(), h1.ID().Pretty())
		peerIDChan <- h1.ID()

		s.Read(make([]byte, 1))

		fmt.Println("finished....")
	}()
	n.MultiAddress = <-multiaddrChan
	n.PeerID = <-peerIDChan
	return n
}

// NewH2Node ...
func NewH2Node() *Node {
	n := &Node{}
	multiaddrChan := make(chan string)
	peerIDChan := make(chan peer.ID)
	go func() {
		// fmt.Println("Opt hop will start")
		h2, err := libp2p.New(context.Background(), libp2p.EnableRelay(circuit.OptHop))
		if err != nil {
			panic(err)
		}

		// fmt.Printf("%+v/p2p/%s", h2.Addrs()[1], h2.ID().Pretty())
		multiaddrChan <- fmt.Sprintf("%+v/p2p/%s", h2.Addrs()[1], h2.ID().Pretty())
		peerIDChan <- h2.ID()
		select {}
	}()
	n.MultiAddress = <-multiaddrChan
	n.PeerID = <-peerIDChan
	return n
}

// NewH3Node ...
func NewH3Node(h2Dest string) *Node {
	n := &Node{}
	multiaddrChan := make(chan string)
	peerIDChan := make(chan peer.ID)

	go func() {
		h3, err := libp2p.New(context.Background(), libp2p.ListenAddrs(), libp2p.EnableRelay())
		if err != nil {
			panic(err)
		}

		multiaddr, err := multiaddr.NewMultiaddr(h2Dest)
		if err != nil {
			panic(err)
		}

		h2Info, err := peer.AddrInfoFromP2pAddr(multiaddr)

		if err = h3.Connect(context.Background(), *h2Info); err != nil {
			panic(err)
		}

		// Now, to test things, let's set up a protocol handler on h3
		h3.SetStreamHandler("/cats1", func(s network.Stream) {
			fmt.Println("Meow! It worked!")
			s.Close()
		})

		multiaddrChan <- fmt.Sprintf("%+v/p2p/%s", h3.Addrs()[1], h3.ID().Pretty())
		peerIDChan <- h3.ID()
		select {}
	}()
	n.MultiAddress = <-multiaddrChan
	n.PeerID = <-peerIDChan
	return n
}
