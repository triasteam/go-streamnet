## 功能说明
该demo用于演示如何在公网环境下使用libp2p组建gossip网络。
首先，需要对本地host进行配置，如下：
```
    // 使用tcp和ws的连接
	transports := libp2p.ChainOptions(
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.Transport(ws.New),
	)

    // 设置多路复用协议
	muxers := libp2p.ChainOptions(
		libp2p.Muxer("/yamux/1.0.0", yamux.DefaultTransport),
		libp2p.Muxer("/mplex/6.7.0", mplex.DefaultTransport),
	)

    // 支持TLS协议
	security := libp2p.Security(secio.ID, secio.New)

    // 监听指定端口
	listenAddrs := libp2p.ListenAddrStrings(
		fmt.Sprintf("/ip4/0.0.0.0/tcp/%s", cfg.Port),
		fmt.Sprintf("/ip4/0.0.0.0/tcp/%s/ws", cfg.Port),
	)
    
    // 设置该节点可以使用DHT发现其他节点
	var dht *kaddht.IpfsDHT
	newDHT := func(h host.Host) (routing.PeerRouting, error) {
		var err error
		dht, err = kaddht.New(ctx, h)
		return dht, err
	}
	routing := libp2p.Routing(newDHT)

    // 设置私钥，固定的peer ID也是由此生成
	priv := getOrGeneratePrivateKey()

    // 创建新的host
	host, err := libp2p.New(
		ctx,
		transports,
		listenAddrs,
		muxers,
		security,
		routing,
		libp2p.Identity(priv),
	)
```

ipfs使用的gossip协议是libp2p的 gossipsub 协议，接下来是该协议的配置方式。
```
    // 创建协议只需要指定host即可，gossip协议自身不关心如何查找节点，由discovery服务
    // 查找到节点之后写入到host的address book供其使用
	ps, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		panic(err)
	}
    // 发送和接收消息是通过topic实现的，订阅了同一个topic的节点相互之间消息共享
	sub, err := ps.Subscribe(pubsubTopic)
	if err != nil {
		panic(err)
	}
	go pubsubHandler(ctx, sub)
```

基于DHT的节点发现服务需要配置种子节点，其他节点在加入网络时需要指定种子节点。这与mdns服务有所区别，后者会
自行查找某一网段的所有节点。不考虑NAT场景，节点具有公网ip和端口就可以成为种子节点。
```
	targetAddr, err := multiaddr.NewMultiaddr(cfg.Seed)
	if err != nil {
		panic(err)
	}

	targetInfo, err := peer.AddrInfoFromP2pAddr(targetAddr)
	if err != nil {
		panic(err)
	}

    // 第一个种子节点启动时会报错，属于正常情况
	err = host.Connect(ctx, *targetInfo)
	if err != nil {
		fmt.Printf("connect to host error: %s \n", err)
	}
```