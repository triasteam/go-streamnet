package network

import (
	"flag"
)

type config struct {
	sp           int
	d            string
	relayAddress string
}

func parseFlags() *config {
	conf := &config{}

	flag.IntVar(&conf.sp, "sp", 0, "Source port number")
	flag.StringVar(&conf.d, "d", "", "destination multiaddr string \n")
	flag.StringVar(&conf.relayAddress, "relay", "", "relay multi address")

	flag.Parse()
	return conf
}
