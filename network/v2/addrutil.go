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

// Package networkv2 is an upgraded version of the networkv1, and provides basic
// network layer components.
package networkv2

import (
	fmt "fmt"
	"log"

	regexp "github.com/dlclark/regexp2"

	"github.com/multiformats/go-multiaddr"
)

// ParseRelayPeerID used to parse peer id from multi address which containing
// "p2p-circuit" flag.
func ParseRelayPeerID(addr multiaddr.Multiaddr) (peerID string) {

	var matchResult bool
	var err error

	regx := `p2p/(.*?)/p2p-circuit`
	regexpstr := regexp.MustCompile(regx, 0)
	addrStr := addr.String()

	if matchResult, err = regexpstr.MatchString(addrStr); err != nil {
		log.Fatalf("cant parse from : %s", addr.String())
		return
	}

	if !matchResult {
		log.Printf("failed parse from : %s", addrStr)
		return
	}
	m, err := regexpstr.FindStringMatch(addrStr)
	if err != nil {
		panic(err)
	}
	peerID = m.Groups()[1].Capture.String()

	fmt.Println(peerID)
	return

}
