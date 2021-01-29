// Copyright 2017 The GoReporter Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package network is now deprecated, more to networkv2
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
