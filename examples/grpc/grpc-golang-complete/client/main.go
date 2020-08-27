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
	execStoreBlock()
	//	log.Printf("time taken: %.2f s ", time.Now().Sub(currTime).Seconds())
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
