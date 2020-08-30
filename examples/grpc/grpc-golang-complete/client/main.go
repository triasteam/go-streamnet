package main

import (
	"fmt"

	proto "github.com/triasteam/go-streamnet/abci/proto"

	//	"log"

	//	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	networkType = "tcp"
	server      = "127.0.0.1" //"172.0.16.111" "172.0.16.105"//"127.0.0.1"
	port        = "50051"
)

func main() {
	// execStoreBlock()
	execGetNodeRank()
	//	log.Printf("time taken: %.2f s ", time.Now().Sub(currTime).Seconds())
}

func execGetNodeRank() {
	conn, err := grpc.Dial(server+":"+port, grpc.WithInsecure())
	if nil != conn {
		defer conn.Close()
	}

	if nil != err {
		fmt.Printf("创建连接失败!%s\n", err)
		return
	}

	client2 := proto.NewStreamnetServiceClient(conn)

	req := &proto.RequestGetNoderank{
		BlockHash: []string{"205cbf5ee8795fa1e298aa72c0340e0037c7def2da84e84f3ccdbef224b36e56",
			"b127780b656504457356b7c8d10b71f2305bf3f8afdc7b4f37941b7f2db0c903",
			"a4f5cb909fc7f2dbb8d2799b7d580f1a1eaa3a2ddbf72020e26a5151c69a309b",
			"375d34bfa6cd272e99980d642eb93773561cd8696bf819ad8e59ef121c9657e3",
			"8b533f950877ac4024429b27828722972b06c0e31f6497134a6730fce9b5d6c5",
			"2f07b3e7732c1fd2a6caa92917c8d68529da3b24bb19d4ef92840a5842826594",
			"53143f252a5ab11db86ec6d235b783494ac389656095f86937c9cdb0cfaa5510",
			"989043b061b9723022caa9de60098354600463fc61dbf13560d253034a698570"},
		Duration: 100,
		Period:   0,
		NumRank:  10,
	}

	result, _ := client2.GetNoderank(context.Background(), req)
	if nil != result {
		fmt.Println(result)
	} else {
		fmt.Println("response is nil")
	}
}

func execStoreBlock() {
	//建立连接
	conn, err := grpc.Dial(server+":"+port, grpc.WithInsecure())
	if nil != conn {
		defer conn.Close()
	}

	if nil != err {
		fmt.Printf("创建连接失败!%s\n", err)
		return
	}

	client2 := proto.NewStreamnetServiceClient(conn)

	req := &proto.RequestStoreBlock{
		BlockInfo: "this is the test blockInfo",
	}

	result, _ := client2.StoreBlock(context.Background(), req)
	if nil != result {
		fmt.Printf("%s \n", string(result.Data))
	} else {
		fmt.Println("response is nil")
	}
}
