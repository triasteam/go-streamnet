package main

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p"
	circuit "github.com/libp2p/go-libp2p-circuit"
)

const RELAY_PORT string = "14701"

// can build a relay node
func main() {
	// fmt.Println("Opt hop will start")
	h2, err := libp2p.New(context.Background(),
		libp2p.ListenAddrStrings("/ip4/127.0.0.1/tcp/"+RELAY_PORT),
		libp2p.EnableRelay(circuit.OptHop))
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v/p2p/%s\n", h2.Addrs(), h2.ID().Pretty())

	select {}
}
