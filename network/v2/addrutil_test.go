package networkv2

import (
	"testing"

	"github.com/multiformats/go-multiaddr"
)

func TestParseRelayMultiAddr(t *testing.T) {
	relayAddr, _ := multiaddr.NewMultiaddr(`/dns4/localhost/tcp/45759/p2p/12D3KooWJV9tYPfLodAqYSCgu2smcCUfTWpahUVponD4U1uAYbbk/p2p-circuit`)

	noRelayAddr, _ := multiaddr.NewMultiaddr(`/dns4/localhost/tcp/45759/p2p/12D3KooWJV9tYPfLodAqYSCgu2smcCUfTWpahUVponD4U1uAYbbk`)

	relayAddrStr := ParseRelayPeerID(relayAddr)
	if relayAddrStr == "" {
		t.Fatal("error")
	}

	relayAddrStr = ParseRelayPeerID(noRelayAddr)
	if relayAddrStr != "" {
		t.Fatal("error")
	}
}
