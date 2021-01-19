package networkv2

import (
	"context"
	"fmt"
	"testing"

	"github.com/libp2p/go-libp2p"
	secio "github.com/libp2p/go-libp2p-secio"
	localConfig "github.com/triasteam/go-streamnet/config"
)

func TestNewNodeWithNoParam(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, _ := localConfig.ParseFlags()

	var receive func(b []byte) error
	_, err := NewNode(ctx, &cfg, receive)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewNodeWithSeed(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	security := libp2p.Security(secio.ID, secio.New)
	seed, err := libp2p.New(ctx, security)
	seedAddr := fmt.Sprintf("%s/p2p/%s", seed.Addrs()[1].String(), seed.ID().Pretty())

	cfg := &localConfig.Config{}
	cfg.Seed = seedAddr
	cfg.Port = "8800"

	var receive func(b []byte) error
	_, err = NewNode(ctx, cfg, receive)
	if err != nil {
		t.Fatal(err)
	}

	// seed与new node建立连接，其addr book长度等于2
	if len(seed.Peerstore().Peers()) != 2 {
		t.Fatal("new node didn't connect to seed.")
	}
}

func TestNewNodeWithAutoRelay(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	seed, err := libp2p.New(ctx)
	seedAddr := fmt.Sprintf("%s/p2p/%s", seed.Addrs()[0].String(), seed.ID().Pretty())

	cfg, err := localConfig.ParseFlags()
	cfg.Seed = seedAddr

	// define relay hop
	cfg.RelayType = "hop"
	var receive func(b []byte) error
	_, err = NewNode(ctx, &cfg, receive)
	if err != nil {
		t.Fatal(err)
	}

	// define relay
	cfg.RelayType = "autorelay"
	_, err = NewNode(ctx, &cfg, receive)
	if err != nil {
		t.Fatal(err)
	}

}

func TestBroadcast(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	security := libp2p.Security(secio.ID, secio.New)
	seed, err := libp2p.New(ctx, security)
	seedAddr := fmt.Sprintf("%s/p2p/%s", seed.Addrs()[0].String(), seed.ID().Pretty())

	cfg, err := localConfig.ParseFlags()
	cfg.Seed = seedAddr

	// var haveReceived bool
	// define relay hop
	receive := func(b []byte) error {
		// haveReceived = true
		if string(b) != "hello" {
			t.Fatal("not hello !")
		}
		return nil
	}

	var fun func([]byte) error
	_, err = NewNode(ctx, &cfg, fun)
	if err != nil {
		t.Fatal(err)
	}

	cfg.Port = "45758"
	_, err = NewNode(ctx, &cfg, receive)
	if err != nil {
		t.Fatal(err)
	}

	// time.Sleep(1 * time.Second)

	// sender.Broadcast([]byte("hello"))

	// time.Sleep(1 * time.Second)

	// if haveReceived == false {
	// 	t.Fatal("haven't received message !")
	// }

	// time.Sleep(1 * time.Second)
}
