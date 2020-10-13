package main

import (
	"flag"
)

type config struct {
	n       int
	address string
	peerID  string
}

func parseFlags() *config {
	conf := &config{}

	flag.IntVar(&conf.n, "n", 2, "node sequence")
	flag.StringVar(&conf.address, "address", "", "node sequence")
	flag.StringVar(&conf.peerID, "peer", "", "node sequence")
	flag.Parse()
	return conf
}
