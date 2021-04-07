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
 * this file handle all argument for input
 *
 */
package main

import (
	"flag"
)

// Config contains ...
type Config struct {
	Seed       string
	Port       string
	RelayType  string
	PublicAddr string
}

// ParseFlags parsing arguments for app
func ParseFlags() (Config, error) {
	config := Config{}
	flag.StringVar(&config.Seed, "seed", "/ip4/127.0.0.1/tcp/45759/ipfs/QmWjz6xb8v9K4KnYEwP5Yk75k5mMBCehzWFLCvvQpYxF3d", "while starting you peer, a seed should be specify")
	flag.StringVar(&config.Port, "port", "45759", "listening port")
	flag.StringVar(&config.RelayType, "relaytype", "", "hop/autorelay")
	flag.StringVar(&config.PublicAddr, "public", "", "public address, as \"154.8.160.48\"")
	flag.Parse()

	return config, nil
}
