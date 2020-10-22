package main

import (
	"flag"
)

type config struct {
	listenAddress string
}

func parseFlags() *config {
	conf := &config{}

	flag.StringVar(&conf.listenAddress, "address", "/ip4/127.0.0.1/tcp/45759", "config relay address")

	flag.Parse()
	return conf
}
