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
go build .
```

__This will generate a binary  file ,The binary's name is '**go-streamnet**', which is standing for Go-StreamNet.__

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
go build .
```
__This will generate a binary  file ,The binary's name is '**go-streamnet**', which is standing for Go-StreamNet.__

## Run
- Start on the same intranet segment  
   One of the nodes run
    ```bash
    ./go-streamnet
    ```
    Other nodes
    ```bash
    ./go-streamnet -d /ip4/ipaddress/tcp/port/p2p/peerid
    ```
    ipaddresss is the first run  nodes's intranet ipaddress,port and peerid is Automatically generated in the first node run.
- Start by relay  
    One of the nodes run
    ```bash
    cd scripts/relay
    go build .
    ./relay  -address /ip4/relay-ipaddress/tcp/relay-port
    ```
    the second nodes
    ```bash
    ./go-streamnet -relay /ip4/relay-ipaddress/tcp/relay-port/p2p/relay-peerid
    ```
    other nodes
    ```bash
    ./go-streamnet -d /ip4/127.0.0.1/tcp/second-port/p2p/second-peerid  -relay  /ip4/relay-ipaddress/tcp/relay-port/p2p/relay-peerid
    ```
    relay-ipaddresss is the first run  nodes's extranet ipaddress,relay-port and relay-peerid is Automatically generated in the first node run.second-port and second-peerid is Automatically generated in the second node run.
    
***Note***:
    Now you should start another terminal to input commands, or you can start the binary background with '&'.

## Client save & get
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
