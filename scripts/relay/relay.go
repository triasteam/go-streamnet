package main

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p"
	circuit "github.com/libp2p/go-libp2p-circuit"
)

// can build a relay node
func main() {
	// fmt.Println("Opt hop will start")
	h2, err := libp2p.New(context.Background(), libp2p.EnableRelay(circuit.OptHop))
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v/p2p/%s\n", h2.Addrs()[1], h2.ID().Pretty())

	select {}
}
