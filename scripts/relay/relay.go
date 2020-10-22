package main

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p"
	circuit "github.com/libp2p/go-libp2p-circuit"
)

// can build a relay node
func main() {
	conf := parseFlags()
	h2, err := libp2p.New(
		context.Background(),
		libp2p.EnableRelay(circuit.OptHop),
		libp2p.ListenAddrStrings(conf.listenAddress),
	)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v/p2p/%s\n", h2.Addrs()[0], h2.ID().Pretty())

	select {}
}
