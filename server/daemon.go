// Package server contain all thing of http server
package server

import (
	"log"
	"net/http"

	"github.com/triasteam/go-streamnet/streamnet"

	streamnet_conf "github.com/triasteam/go-streamnet/config"
)

var (
	server *http.Server
	sn     *streamnet.StreamNet
)

// Start a http server
func Start(stream *streamnet.StreamNet) {
	//TODO: find a better way to check whether server has started.
	if server != nil {
		log.Printf("Server already started.\n")
		return
	}

	log.Printf("Go-StreamNet server is starting...\n")

	// set global data
	sn = stream

	// http server
	mux := http.NewServeMux()
	mux.Handle("/", &gsnHandler{})
	mux.HandleFunc("/save", SaveHandle)
	mux.HandleFunc("/get", GetHandle)

	server = &http.Server{
		Addr:    streamnet_conf.EnvConfig.Port,
		Handler: mux,
		//WriteTimeout: time.Second * 3,
	}

	log.Fatal(server.ListenAndServe())
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
