## 前言
基于go-libp2p-pubsub开发的一个gossip实现例子。
## 部署方法
```
$ git clone https://github.com/triasteam/StreamNet-go.git
$ pwd
~/StreamNet-go
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
