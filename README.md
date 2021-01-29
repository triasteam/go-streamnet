# go-streamnet
StreamNet implemented in Golang.

## Modules

* config: parsing the configuration file.
* dag: the dag structure and methods.
* streamwork:  the algorithm to 'drag' DAG to chain.
* tipselection: selecting two tips from dag.
* server: the http server.
* store: the storage module.
* forward: the genesis forward module.
* network: the network module, for broadcasting and synchronization.
* docs: documents.
* examples: all the tests and examples.
* types: all the common types.
* utils: all common libs.
* commands: parsing all subcommands.


## Install
### From Binary
(not supported yet)

### From Source

#### prerequisite 

* Golang version >= v1.14.2, [installed](https://golang.org/doc/install) 

#### Ubuntu
You'll need  [[gorocksdb]](https://github.com/triasteam/go-streamnet/blob/master/docs/software/gorocksdb%20%E5%AE%89%E8%A3%85.md)[[redis]](https://redis.io/download) installed first.

##### Get Source Code

```bash
mkdir -p $GOPATH/src/github.com/triasteam
cd $GOPATH/src/github.com/triasteam
git clone https://github.com/triasteam/go-streamnet.git
cd go-streamnet
```

##### Compile

```bash
go build -o main .
```

__This will generate a executable file, named by ```-o``` param. ```.``` stands for current directory, therefore, you should execute the command above in ```go-streamnet``` project root directory.__

#### MacOS

##### dependencies
- RocksDb
  ```
  $ brew install rocksdb
  ```

##### Get Source Code

```bash
mkdir -p $GOPATH/src/github.com/triasteam
cd $GOPATH/src/github.com/triasteam
git clone https://github.com/triasteam/go-streamnet.git
cd go-streamnet
```

##### Compile

```bash
go build -o main .
```
__This will generate a executable file, named by ```-o``` param. ```.``` stands for current directory, therefore, you should execute the command above in ```go-streamnet``` project root directory.__

## Run
- first, let us see while parameters can be used to execute this command.
```
**NAME**
  main -- go-stream project name specified by -o param, while can be to any symbol, change whatever you want.
**SYNOPSIS**
  ./main [-seed] [-port port] [-relaytype autorelay/hop] [-paddr publicAddress]
**DESCRIPTION**
  The main command is used to startup a go-streamnet node. It can be a single host, or a relay hop providing public access channel for nodes hiding the NAT. Multiple nodes defining by this way can form a gossip network. 
  -seed: Specify a node as seed node. If you want to join a network, this param is required.

  -port: Default 45759 for node rcp service, can be changed to any other number.

  -relaytype: If autorelay service is neccessary, this argument is required. Node can choose a identify of autorelay and hop, and only one of them can be seleced. Note that a hop must have a public address.(If deploy you app using docker, /dns4/localhost is used as a hop address.

  -paddr Means public address for node while it launching in internet environment.

```

- Start on the same intranet segment  
   One of the nodes run
    ```bash
    ./main
    ```
    Other nodes
    ```bash
    ./main -seed /ip4/ipaddress/tcp/port/p2p/peerid
    ```
    ipaddresss is the first run  nodes's intranet ipaddress,port and peerid is Automatically generated in the first node run. The address can be found in startup log, and you should replace 127.0.0.1 to actural ip/dns address. 
- Start by relay  
    first, a hop should be started, which can be specified by ```-relaytype hot```. A hop node is not suggested launching as a seed, for the advising server limited. The hop node will advise it's autorelay protocol only once after start of few minutes. If there is no connection, advising will failed and no one knowns it's a relay hop.
    ```bash
    ./main  -seed /ip4/192.168.50.102/tcp/45759/p2p/Qmxxx -relaytype hop
    ```
    Second, start a node need relay.
    ```bash
    ./main  -seed /ip4/192.168.50.102/tcp/45759/p2p/Qmxxx -relaytype autorelay
    ```
    other situations, rcp port and public address can be specefied, but it is not required in LAN.
    ```bash
    ./main  -seed /ip4/192.168.50.102/tcp/45759/p2p/Qmxxx -relaytype autorelay -port 10001 -paddr 10.1.12.12
    ```
    
***Note***:
    Now you should start another terminal to input commands, or you can start the binary background with '&'.

## Client save & get
A go-streamnet provides http server by fixed port 14700.
### save
```bash
curl -X POST -d '{"Attester": "192.168.1.1", "Attestee": "192.168.1.2", "Score": "1"}' http://127.0.0.1:14700/save
```

### query
```bash
curl -s -X POST http://127.0.0.1:14700/QueryNodes -d "{\"period\":1,\"numRank\":100}"
```

### get
```bash
curl -X POST -d '{"key": "0xXXXXXXXXXXX"}' http://127.0.0.1:14700/get
```

For better display, you can use `jq` to format the output.
