# relay-demo examples and turorials
relay-demo是用于演示中继功能的例子，是指在B节点的中继下，使得原来不能直接通信的A节点与C节点通过B节点可以通信，且后续所有通信流量都经过B节点。
A <---> B <---> C

## Step 0
准备Go语言环境，安装相关包。
通过build命令来编译源码文件，完成后按照顺序启动B，A，C三个节点。
`go build main.go`
如果条件允许，可以搭建docker环境，利用bridge可以做到A和C两个节点网络隔离。

## Step 1
先使用-c 2参数启动中继节点B，命令如下：
`$ ./main -c 2`

此时，中继节点B已启动，
```
This node's multiaddresses:
 - /ip4/127.0.0.1/tcp/51140
 - /ip4/192.168.199.115/tcp/51140
 - /ip4/192.168.192.160/tcp/51140

Run './main -c 1 -d /ip4/{ip}/tcp/51140/p2p/Qmbp1hxgcaLpBQAhZw4965jWKPtSMjrMHUWci5v6e8aeZF' on another console.

I am relay
```
启动后提示信息如上所述，给出了绑定的所有网卡的地址信息，以及A节点的启动命令
`./main -c 1 -d /ip4/{ip}/tcp/51140/p2p/Qmbp1hxgcaLpBQAhZw4965jWKPtSMjrMHUWci5v6e8aeZF`，后续根据需要选择合适的网卡地址信息以替换上述命令的`{ip}`部分。

## Step 2
启动A节点，这里直接使用127.0.0.1地址测试，替换`{ip}`部分后命令如下，`./main -c 1 -d /ip4/127.0.0.1/tcp/51140/p2p/Qmbp1hxgcaLpBQAhZw4965jWKPtSMjrMHUWci5v6e8aeZF`

此时，中继节点A已启动，
```
This node's multiaddresses:
 - /ip4/127.0.0.1/tcp/51226
 - /ip4/192.168.199.115/tcp/51226
 - /ip4/192.168.192.160/tcp/51226

Run './main -c 3 -d /ip4/127.0.0.1/tcp/51140/p2p/Qmbp1hxgcaLpBQAhZw4965jWKPtSMjrMHUWci5v6e8aeZF/p2p-circuit/p2p/QmZKZoyNsomMyDmeBcZJLC6pX2MhVahWyndGZydoKrCFJC' on another console.
```
启动后提示信息如上所述，给出了绑定的所有网卡的地址信息，这部分不用关心，
注意C节点的启动命令
`./main -c 3 -d /ip4/127.0.0.1/tcp/51140/p2p/Qmbp1hxgcaLpBQAhZw4965jWKPtSMjrMHUWci5v6e8aeZF/p2p-circuit/p2p/QmZKZoyNsomMyDmeBcZJLC6pX2MhVahWyndGZydoKrCFJC`，后续可能需要选择合适的网卡地址信息以替换上述命令的`127.0.0.1`部分，具体以C能连接到B的相应网卡地址信息为准。


## Step 3
启动C节点，这里直接使用127.0.0.1地址测试，命令如下，
`./main -c 3 -d /ip4/127.0.0.1/tcp/51140/p2p/Qmbp1hxgcaLpBQAhZw4965jWKPtSMjrMHUWci5v6e8aeZF/p2p-circuit/p2p/QmZKZoyNsomMyDmeBcZJLC6pX2MhVahWyndGZydoKrCFJC`，

至此，A与C的连接建立，双方可以互相在console中通信。
