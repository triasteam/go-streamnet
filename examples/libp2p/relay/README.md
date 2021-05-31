## 运行demo
### 1. 选择三台及其分别作为中继节点和两台普通节点
### 2. 更新代码并编译
```
## 如果core版本不是0.6.1需要手动更新该版本核心包： go get -v "github.com/libp2p/go-libp2p-core"@v0.6.1
## 如果编译失败请将该relay项目移动到 $GOPATH/src/github.com/目录下
$ go build . -o demo
```

### 3. 启动中继节点
```
## 日志中 multiaddr 的ip地址在实际使用时需要替换为公网ip，同时要确保端口是可访问的，可以使用telnet测试下。
## 本例中中继节点的外网ip是154.8.160.48
## TODO 中继节点使用固定的端口而不是随机生成的。
$ ./demo
n2 p2p address is  /ip4/172.21.32.6/tcp/45759/p2p/QmSAewBQ4UgXsuPm8XSjqdzxi4j8SgPy34ia9YvC1Xid7n
```

### 4. 启动第一个节点
```
$ ./demo -n 3 -address /ip4/154.8.160.48/tcp/45759/p2p/QmSAewBQ4UgXsuPm8XSjqdzxi4j8SgPy34ia9YvC1Xid7n
n3 p2p address is  /ip4/127.0.0.1/tcp/43221/p2p/QmUUXBfnnuiaPhFBvi1F7TEQHH8bYSzQBfwP2ziHGuY9Zh
```

### 5. 启动第二个节点，通过中继访问第一个节点
```
$ ./demo -n 1 -address /ip4/154.8.160.48/tcp/45759/p2p/QmSAewBQ4UgXsuPm8XSjqdzxi4j8SgPy34ia9YvC1Xid7n -peer QmUUXBfnnuiaPhFBvi1F7TEQHH8bYSzQBfwP2ziHGuY9Zh
n1 address is  /ip4/172.21.32.12/tcp/44029/p2p/QmVhqy8oXu1CBqRoxwGoD1xHqQdvdAFAedjvk6RMsZEDHj
## 该日志没有实际作用，只是说明连接成功了。同时第一个节点也会打印连接成功日志。
```
