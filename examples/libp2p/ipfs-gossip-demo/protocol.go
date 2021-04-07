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
 * This program demonstrate a gossip application using p2p pubsub protocol and
 * AutoRelay protocol.
 * With pubsub protocol nodes could be connected and tranlating, with the AutoRelay
 * Protocol, node behind NAT can also join the network.
 *
 * this file describes functions of the demo
 */
package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"fmt"
	"os"
	"time"

	"strings"

	"github.com/libp2p/go-libp2p-core/host"
	peer "github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

func sendMessage(ps *pubsub.PubSub, msg string) {
	msgId := make([]byte, 10)
	_, err := rand.Read(msgId)
	defer func() {
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()
	if err != nil {
		return
	}
	now := time.Now().Unix()
	req := &Request{
		Type: Request_SEND_MESSAGE.Enum(),
		SendMessage: &SendMessage{
			Id:      msgId,
			Data:    []byte(msg),
			Created: &now,
		},
	}
	msgBytes, err := req.Marshal()
	if err != nil {
		return
	}
	err = ps.Publish(pubsubTopic, msgBytes)
}

func updatePeer(ps *pubsub.PubSub, id peer.ID, handle string) {
	oldHandle, ok := handles[id.String()]
	if !ok {
		oldHandle = id.ShortString()
	}
	handles[id.String()] = handle

	req := &Request{
		Type: Request_UPDATE_PEER.Enum(),
		UpdatePeer: &UpdatePeer{
			UserHandle: []byte(handle),
		},
	}
	reqBytes, err := req.Marshal()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	err = ps.Publish(pubsubTopic, reqBytes)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	fmt.Printf("%s -> %s\n", oldHandle, handle)
}

func chatInputLoop(ctx context.Context, h host.Host, ps *pubsub.PubSub, donec chan struct{}) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		msg := scanner.Text()
		if strings.HasPrefix(msg, "/name ") {
			newHandle := strings.TrimPrefix(msg, "/name ")
			newHandle = strings.TrimSpace(newHandle)
			updatePeer(ps, h.ID(), newHandle)
		} else {
			sendMessage(ps, msg)
		}
	}
	donec <- struct{}{}
}
