package main
 
import (
	"fmt"
	proto "triasteam/StreamNet-go/examples/grpc/grpc-golang-complete/example"
//	"log"
	"runtime"
	"sync"
//	"time"
 
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)
 
var (
	wg sync.WaitGroup
)
 
const (
	networkType = "tcp"
	server      =  "127.0.0.1" //"172.0.16.111" "172.0.16.105"//"127.0.0.1"
	port        = "41005"
	parallel    = 1        //连接并行度
	times       = 1    //每连接请求次数
)
 
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
//	currTime := time.Now()
 
	//并行请求
	for i := 0; i < int(parallel); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			exe()
		}()
	}
	wg.Wait()
 
//	log.Printf("time taken: %.2f s ", time.Now().Sub(currTime).Seconds())
}
 
func exe() {
	//建立连接
	conn, err := grpc.Dial(server + ":" + port,grpc.WithInsecure())
	if nil!=conn{
		defer conn.Close()
	}
 
	if (nil!=err){
		fmt.Printf("创建连接失败!%s\n",err)
		return
	}
 
	client2 := proto.NewSayHelloServiceClient(conn)
 
	for i := 0; i < int(times); i++ {
		testSayHello(client2)
	}
}
 
 
 
func testSayHello(client proto.SayHelloServiceClient){
	req:=&proto.SayHelloRequest{
		Name:[]byte("Tony"),
	}
 
	result,_:=client.SayHello(context.Background(),req)
	if nil!=result{
		fmt.Printf("%s \n",string(result.Result))
	}else{
		fmt.Println("response is nil")
	}
}