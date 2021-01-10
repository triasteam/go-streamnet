// 中继可用于解决不方便暴露到公网的节点通过某一个公网节点建立网络连接。
// Libp2p提供了一个/libp2p/autorelay的内容发现协议可以将某一个中继发布给其他节点
// 中继流程：（普通中继演示了由一个位于公网的中继节点和两个位于私网节点构成的简单网络）
// 1. 启动中继节点；
// 2. 启动私网节点1并连接至中继，获取该私网节点的multiaddress，形如：/ip4/{中继IP地址}/tcp/{中继端口}/p2p/{中继PeerId}/p2p-circuit/p2p/{私网节点PeerId}
// 3. 启动私网节点2，指定连接到步骤2中节点地址。
// 4. 以上就是普通中继建立的连接的步骤。
// 自动中继流程：(网络拓扑同上)
// 1. 启动自动中继节点；
// 2. 启动私网节点1并连接至自动中继节点，不需要额外配置。由于自动中继推送服务在15分钟之后启动且只会推动一次，因此该节点需要在15分钟之类启动并连接至自动中继。
// 3. 启动私网节点2并连接至自动中继节点, 同样不需要额外配置；
// 3. 自动中继节点推送消息 /libp2p/autorelay；
// 4. 私网节点接收到 /libp2p/autorelay 消息，将创建该消息的节点的节点作为中继，创建基于该中继的地址。形如“中继流程”第2步所示；
// 5. 私网节点通过Identify协议将自身地址推送给邻居节点。邻居节点也会将该地址广播给网络中的其他节点，任何节点都可以使用该地址与私网节点通信。

package main

import (
	"context"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/routing"
	kaddht "github.com/libp2p/go-libp2p-kad-dht"
	mplex "github.com/libp2p/go-libp2p-mplex"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	secio "github.com/libp2p/go-libp2p-secio"
	yamux "github.com/libp2p/go-libp2p-yamux"
	"github.com/libp2p/go-libp2p/p2p/discovery"
	tcp "github.com/libp2p/go-tcp-transport"
	ws "github.com/libp2p/go-ws-transport"
	"github.com/multiformats/go-multiaddr"
)

const (
	privateName = "priv.pem"
)

type mdnsNotifee struct {
	h   host.Host
	ctx context.Context
}

// 实现节点连接
func (m *mdnsNotifee) HandlePeerFound(pi peer.AddrInfo) {
	m.h.Connect(m.ctx, pi)
}

// 用于实现生成该节点的PeerId
func loadFromPem() (crypto.PrivKey, error) {
	pwd, _ := os.Getwd()
	filePath := fmt.Sprintf("%s/%s", pwd, privateName)
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(bytes)
	priv, _ := crypto.UnmarshalPrivateKey(block.Bytes)
	return priv, nil
}

func exportToPem(priv crypto.PrivKey) error {
	privBytes, err := crypto.MarshalPrivateKey(priv)
	if err != nil {
		fmt.Printf("marshal private key error: %s \n", err)
		return err
	}
	privPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "Ed25519 PRIVATE KEY",
			Bytes: privBytes,
		},
	)
	fmt.Println("private pem: ", privPem)
	ioutil.WriteFile(privateName, []byte(privPem), 0644)
	return nil
}

func getOrGeneratePrivateKey() crypto.PrivKey {
	priv, err := loadFromPem()
	if err != nil {
		priv, _, err = crypto.GenerateKeyPair(crypto.Ed25519, -1)
		if err != nil {
			panic(err)
		}
		err = exportToPem(priv)
		if err != nil {
			fmt.Printf("export private key error: %s \n", err)
		}
	}
	return priv
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, _ := ParseFlags()

	transports := libp2p.ChainOptions(
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.Transport(ws.New),
	)

	muxers := libp2p.ChainOptions(
		libp2p.Muxer("/yamux/1.0.0", yamux.DefaultTransport),
		libp2p.Muxer("/mplex/6.7.0", mplex.DefaultTransport),
	)

	security := libp2p.Security(secio.ID, secio.New)

	listenAddrs := libp2p.ListenAddrStrings(
		fmt.Sprintf("/ip4/0.0.0.0/tcp/%s", cfg.Port),
		fmt.Sprintf("/ip4/0.0.0.0/tcp/%s/ws", cfg.Port),
	)

	var dht *kaddht.IpfsDHT
	newDHT := func(h host.Host) (routing.PeerRouting, error) {
		var err error
		dht, err = kaddht.New(ctx, h)
		return dht, err
	}
	routing := libp2p.Routing(newDHT)

	priv := getOrGeneratePrivateKey()

	host, err := libp2p.New(
		ctx,
		transports,
		listenAddrs,
		muxers,
		security,
		routing,
		libp2p.Identity(priv),
	)
	if err != nil {
		panic(err)
	}
	fmt.Printf("my peer is /ip4/127.0.0.1/tcp/%s/ipfs/%s \n", cfg.Port, host.ID().Pretty())

	ps, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		panic(err)
	}
	sub, err := ps.Subscribe(pubsubTopic)
	if err != nil {
		panic(err)
	}
	go pubsubHandler(ctx, sub)

	for _, addr := range host.Addrs() {
		fmt.Println("Listening on", addr)
	}

	targetAddr, err := multiaddr.NewMultiaddr(cfg.Seed)
	if err != nil {
		panic(err)
	}

	targetInfo, err := peer.AddrInfoFromP2pAddr(targetAddr)
	if err != nil {
		panic(err)
	}

	err = host.Connect(ctx, *targetInfo)
	if err != nil {
		fmt.Printf("connect to host error: %s \n", err)
	} else {
		fmt.Println("Connected to", targetInfo.ID)
	}

	mdns, err := discovery.NewMdnsService(ctx, host, time.Second*10, "")
	if err != nil {
		panic(err)
	}
	mdns.RegisterNotifee(&mdnsNotifee{h: host, ctx: ctx})

	err = dht.Bootstrap(ctx)
	if err != nil {
		panic(err)
	}

	donec := make(chan struct{}, 1)
	go chatInputLoop(ctx, host, ps, donec)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT)

	select {
	case <-stop:
		host.Close()
		os.Exit(0)
	case <-donec:
		host.Close()
	}
}
