## 使用方法
```
git clone https://github.com/triasteam/go-streamnet

cd go-streamnet/examples/libp2p/gossipsub-with-relay

go build .

// 默认监听45759端口，注意种子节点的45759需要开放到公网
// 初次启动会创建私钥文件，该私钥用于生成固定的peer id
./main

```

## 启动参数说明

  - ```-seed``` 指定种子地址，完整格式为 ```/ip4/xxx.xxx.xxx.xxx/tcp/45759/ipfs/qmxxx```
  - ```-port``` 如果开放的外网端口不是45759需要通过该参数指定，一般情况下保持默认即可
  - ```-relaytype``` auto relay 类型，默认不指定表示该节点为非auto relay节点。可选取值为hop/autorelay
  - ```-public``` 公网IP，如果该节点被配置为了HOP则需要指定该节点的公网IP

## 构造星型网络
- Build image  
docker build -t ${image_name}:${tag} .
- Build the network  
docker network create -d bridge ${network_name}
- Run a container  
docker run -itd --name ${container_name} -p ${hostPort}:${containerPort}  ${image_name}:${tag}
- Connect to specified network  
docker network connect  ${network_name}  ${container_name}
- Disconnect the bridge network of all containers  
docker network disconnect bridge   ${container_name}  

1. 启动种子，种子节点兼autorelay hop。注意默认状态下该节点启动15分钟才会对外发布autorelay地址
   ./main -relaytype hop
2. 启动七个docker：
   ./main -seed /ip4/x.x.x.x/tcp/xxxx/ipfs/Qm -relaytype autorelay


## 功能说明
该demo用于演示如何在复杂网络环境下使用libp2p组建gossip网络。对等节点如果处于NAT环境中则可以采用中继的方式组建网络。

关键配置代码如下：
```
    // 首先，需要对本地host进行配置
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

    // 配置中继，如果参数未指定中继则不进行配置。
	// 对等节点只可以是 普通节点/中继hop/自动中继 其中的一种
	relayOption := func() config.Option {
		if cfg.RelayType == "hop" {
			return libp2p.ChainOptions(libp2p.EnableAutoRelay(), libp2p.EnableRelay(circuit.OptHop), libp2p.AddrsFactory(func(addrs []multiaddr.Multiaddr) []multiaddr.Multiaddr {
				for i, addr0 := range addrs {
					saddr := addr0.String()
					if strings.HasPrefix(saddr, "/ip4/127.0.0.1") {
						addrNoIP := strings.TrimPrefix(saddr, "/ip4/127.0.0.1")
						fmt.Printf("result : %d, public: %s \n", len(cfg.PublicAddr), cfg.PublicAddr)
						if len(cfg.PublicAddr) == 0 {
							addrs[i] = multiaddr.StringCast("/dns4/localhost" + addrNoIP)
						} else {
							addrs[i] = multiaddr.StringCast(fmt.Sprintf("/ip4/%s", cfg.PublicAddr) + addrNoIP)
						}
					}
				}
				return addrs
			}))
		} else if cfg.RelayType == "autorelay" {
			return libp2p.ChainOptions(libp2p.EnableAutoRelay())
		}
		return func(cfg *config.Config) error { return nil }

	}	

    // 创建新的host
	host, err := libp2p.New(
		ctx,
		transports,
		listenAddrs,
		muxers,
		security,
		routing,
		libp2p.Identity(priv),
		relayOption,
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
