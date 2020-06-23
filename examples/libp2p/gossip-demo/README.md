## 前言
基于go-libp2p-pubsub开发的一个gossip实现例子。

## 实现步骤
当待接入节点与gossip网络中任意节点建立连接后即可加入该网络。
### 1 添加节点
libp2p使用multiaddr拨号其他节点，因此连接前需要得到被接入网络节点的multiaddr。pubsub采用发布订阅的方式进行gossip，当节点成功连接到另一个节点时，该节点就成为一个发布订阅者（发布+订阅）。主要实现代码如下：
```
maddr, err := multiaddr.NewMultiaddr(dest)
if err != nil {
  panic(err)
}

pi, err := peer.AddrInfoFromP2pAddr(maddr)
if err != nil {
  panic(err)
}

err = host.Connect(ctx, *pi)
if err != nil {
  panic(err)
}
```

### 2 广播消息
广播消息实现较为简单，当节点发布消息的时gossip网络中的其他节点则会收到消息。
```
# 发布消息
ps.Publish(cfg.Topic, []byte(s))

# 订阅消息, subs.Next() 是阻塞方法
msg, err := subs.Next(ctx)
if err != nil {
  panic(err)
}
```

### 3 交互方式
为了简化测试时的交互，节点启动gossip后会同时启动一个http服务，客户端使用http服务而非交互命令行的方式与gossip网络进行交互。
```
http.HandleFunc("/addPeers", func(w http.ResponseWriter, req *http.Request) {
  # ignore other code
  ...
  err = host.Connect(ctx, *pi)
  ...
})
http.HandleFunc("/send", func(w http.ResponseWriter, req *http.Request) {
   # ignore other 
   ...
   ps.Publish(cfg.Topic, []byte(s))
   ...
})
http.ListenAndServe(fmt.Sprintf(":%v", cfg.httpServerPort), nil)
```

### 4 注意事项
该演示demo只使用了pubsub实现的gossip方案，不包含节点发现、路由等功能因此需要指定现有gossip网络中任意一节点的multiaddr，连接成功后才算是接入了gossip网络。也不包含NAT穿透功能因此不能与广域网节点组件gossip网络。

## 部署方法
```
$ git clone https://github.com/triasteam/go-streamnet.git
$ pwd
~/go-streamnet
$ cd gossip-demo
$ go build -o main
$ ./main -hp 8001 -s 1
this multiaddr is :  /ip4/127.0.0.1/tcp/51847/p2p/12D3KooWE4qDcRrueTuRYWUdQZgcy7APZqBngVeXRt4Y6ytHizKV

# 多开窗口执行main程序，http端口是用来交互的，必须重新自定义
# -s 参数只有当测试机全部在同一台物理机时需要指定且不能重复
$ ./main -hp 8002 -s 2
this multiaddr is :  /ip4/127.0.0.1/tcp/51864/p2p/12D3KooWB1b3qZxWJanuhtseF3DmPggHCtG36KZ9ixkqHtdKH9fh

# 再打开一个窗口，此时如果只是简单的将两台机器组成gossip网络则直接向其中一台机器发送另一台机器的multi address即可
$ curl "http://127.0.0.1:8001/addPeers?dest=/ip4/127.0.0.1/tcp/51864/p2p/12D3KooWB1b3qZxWJanuhtseF3DmPggHCtG36KZ9ixkqHtdKH9fh"

# 此时8001“机器”将打印 “peer connected!”
$ ./gossip-demo -hp 8001 -s 1
this multiaddr is :  /ip4/127.0.0.1/tcp/51847/p2p/12D3KooWE4qDcRrueTuRYWUdQZgcy7APZqBngVeXRt4Y6ytHizKV
peer connected!

# 发送消息，选择任一节点执行http命令即可
# 请求格式：curl "http://127.0.0.1:8002/send?msg=hellokatty"
$ curl "http://127.0.0.1:8002/send?msg=hellokatty"
<另外两个客户端将同时打印结果>

```
