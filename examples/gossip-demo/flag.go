package main

import (
	"flag"
)

type config struct {
	Topic          string
	httpServerPort int
	seq            int
}

func parseFlags() *config {
	conf := &config{}

	flag.StringVar(&conf.Topic, "topic", "meetme", "share with other peer let them to connect you \n")
	flag.IntVar(&conf.httpServerPort, "hp", 8001, "http listen port")
	flag.IntVar(&conf.seq, "s", 0, "local test, peer id seq")

	flag.Parse()
	return conf
}
