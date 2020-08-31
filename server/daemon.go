// Package server contain all thing of http server
package server

import (
	"log"
	"net/http"

	"net"

	pb "github.com/triasteam/go-streamnet/abci/proto"
	streamnet_conf "github.com/triasteam/go-streamnet/config"
	"github.com/triasteam/go-streamnet/streamnet"
	"google.golang.org/grpc"
)

var (
	address = "localhost"
	rpcPort string
	server  *http.Server
	sn      *streamnet.StreamNet
)

func Start(stream *streamnet.StreamNet) {

	go startWeb(stream)
	go startGrpc()

	select {}

}

func startWeb(stream *streamnet.StreamNet) {
	if server != nil {
		log.Printf("Server already started.\n")
		return
	}

	log.Printf("Go-StreamNet web-server is starting...\n")

	// set global data
	sn = stream

	// http server
	mux := http.NewServeMux()
	mux.Handle("/", &gsnHandler{})
	mux.HandleFunc("/save", SaveHandle)
	mux.HandleFunc("/get", GetHandle)
	//mux.HandleFunc("/QueryNodes", QueryNodesHandle)

	server = &http.Server{
		Addr:    streamnet_conf.EnvConfig.Port,
		Handler: mux,
		//WriteTimeout: time.Second * 3,
	}

	log.Fatal(server.ListenAndServe())
}

func startGrpc() {
	log.Printf("Go-StreamNet grpc-server is starting...\n")
	rpcPort = streamnet_conf.EnvConfig.GRPC.Port
	lis, err := net.Listen("tcp", rpcPort)
	if err != nil {
		log.Fatalf("failed to listen : %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterStreamnetServiceServer(s, NewAbciServer())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// Stop the server
func Stop() {
	log.Printf("Go-StreamNet server is closing...\n")
	err := server.Shutdown(nil)
	if err != nil {
		log.Printf("!!! Failed to close Go-StreamNet: %v\n", err)
	}
}

type gsnHandler struct{}

func (*gsnHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello, go-streamnet.\n"))
}
