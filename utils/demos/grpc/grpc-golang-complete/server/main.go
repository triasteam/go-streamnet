package main
 
import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	proto "triasteam/StreamNet-go/utils/demos/grpc/grpc-golang-complete/example"
	"net"
	"runtime"
)
 
const (
	port = "41005"
)
 
type Data struct{}
 
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	//起服务
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	data:=&Data{}
	proto.RegisterSayHelloServiceServer(s,data)
	log.Printf("grpc server in: %s", port)
	s.Serve(lis)
 
}
 
func (t *Data) SayHello(ctx context.Context,in *proto.SayHelloRequest) (result *proto.SayHelloResponse, err error){
	return &proto.SayHelloResponse{
		Result:[] byte("hello :"+string(in.Name)),
	},nil
}