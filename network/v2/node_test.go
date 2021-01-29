package networkv2

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p"
	secio "github.com/libp2p/go-libp2p-secio"
	"github.com/libp2p/go-libp2p/p2p/host/relay"
	localConfig "github.com/triasteam/go-streamnet/config"
)

func TestNewNode(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//************ init param  ***********
	relay.AdvertiseBootDelay = 1 * time.Second

	var result []byte
	OnReceived := func(b []byte) error {
		result = b
		return nil
	}

	err := os.Remove("priv.pem")
	if err != nil {
		t.Logf("priv.pem not exists, ignore this.")
	}
	//************ init param  ***********

	// 1. Create a seed node, only responsible for bootstrap initing a network
	security := libp2p.Security(secio.ID, secio.New)
	seed, err := libp2p.New(ctx, security)
	seedAddr := fmt.Sprintf("%s/p2p/%s", seed.Addrs()[1].String(), seed.ID().Pretty())

	time.Sleep(1 * time.Second)

	// 2. Create a relay hop
	hopCfg := &localConfig.Config{}
	hopCfg.Seed = seedAddr
	hopCfg.Port = "8800"
	hopCfg.RelayType = "hop"

	_, err = NewNode(ctx, hopCfg, OnReceived)
	if err != nil {
		t.Fatal(err)
	}
	os.Remove("priv.pem")

	time.Sleep(1 * time.Second)

	// 3. Create a normal node need relay to join the network
	needRelayConf := &localConfig.Config{}
	needRelayConf.Seed = seedAddr
	needRelayConf.Port = "8801"
	needRelayConf.RelayType = "autorelay"

	needRelayNode, err := NewNode(ctx, needRelayConf, OnReceived)
	if err != nil {
		t.Fatal(err)
	}
	os.Remove("priv.pem")

	time.Sleep(1 * time.Second)

	// 4. Create a normal node need relay to join the network
	anotherNeedRelayConf := &localConfig.Config{}
	anotherNeedRelayConf.Seed = seedAddr
	anotherNeedRelayConf.Port = "8802"
	anotherNeedRelayConf.RelayType = "autorelay"

	_, err = NewNode(ctx, anotherNeedRelayConf, OnReceived)
	if err != nil {
		t.Fatal(err)
	}
	os.Remove("priv.pem")

	time.Sleep(1 * time.Second)

	message := []byte("hello")
	needRelayNode.Broadcast(message)

	// seed与new node建立连接，其addr book长度等于4
	if len(seed.Peerstore().Peers()) != 4 {
		t.Fatalf("expected 4, actural %d", len(seed.Network().Peers()))
	}

	time.Sleep(10 * time.Second)

	if reflect.DeepEqual(result, message) {
		t.Log("successed!")
	} else {
		t.Fatalf("expected %s, actural %s", message, result)
	}
	ctx.Done()
	return
}
