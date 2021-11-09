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


## 启动方法
```
git clone https://github.com/triasteam/go-streamnet

cd go-streamnet/examples/libp2p/ipfs-gossip-demo

go build main.go flags.go protocol.go pubsub.go chat.pb.go

// 默认监听45759端口，注意种子节点的45759需要开放到公网
// 初次启动会创建私钥文件，该私钥用于生成固定的peer id
./main

```

## 启动参数说明

  - ```-seed``` 指定种子地址，完整格式为 ```/ip4/xxx.xxx.xxx.xxx/tcp/45759/ipfs/qmxxx```
  - ```-port``` 如果开放的外网端口不是45759需要通过该参数指定，一般情况下保持默认即可

## 单机多服务启动
需要分别启动三个终端
### 启动第一个"节点"
```shell
$ ./main
my peer is /ip4/127.0.0.1/tcp/45759/ipfs/12D3KooWJ7RPvnEPM6dE9EVHNEJ4vcpdAKAB7uud7SNDJk6AY5VN
Listening on /ip4/127.0.0.1/tcp/45759
Listening on /ip4/192.168.50.229/tcp/45759
connect to host error: failed to dial QmWjz6xb8v9K4KnYEwP5Yk75k5mMBCehzWFLCvvQpYxF3d: no good addresses
```
### 启动第二个"节点"
```shell
# 私有文件绑定生成peerId，因此需要删除或者重命名之前已经存在的密钥文件
$ mv priv.pem priv.pem.1
# 指定种子节点和端口
$ ./main -seed /ip4/127.0.0.1/tcp/45759/ipfs/12D3KooWJ7RPvnEPM6dE9EVHNEJ4vcpdAKAB7uud7SNDJk6AY5VN -port 45760
```
### 启动第三个"节点"
```shell
$ mv priv.pem priv.pem.2
# 此处的种子不一定是第一个节点的多播地址，也可以是其他已经加入网络的节点地址
$ ./main -seed /ip4/127.0.0.1/tcp/45759/ipfs/12D3KooWJ7RPvnEPM6dE9EVHNEJ4vcpdAKAB7uud7SNDJk6AY5VN -port 45761
```
### 通信
在任意一个输入消息后回车确定，其他终端会同步收到消息。


