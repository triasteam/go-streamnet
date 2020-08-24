// Package server contain all thing of http server
package server

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	streamnet_conf "github.com/triasteam/go-streamnet/config"
	"github.com/triasteam/go-streamnet/store"
	"github.com/triasteam/go-streamnet/types"
)

var (
	server *http.Server
	db     *store.Storage
)

// Start a http server
func Start(store *store.Storage) {
	//TODO: find a better way to check whether server has started.
	if server != nil {
		log.Printf("Server already started.\n")
		return
	}

	log.Printf("Go-StreamNet server is starting...\n")

	// set db
	db = store

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

// SaveHandle process the 'save' request.
func SaveHandle(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	// params
	var params types.StoreData
	err := decoder.Decode(&params)
	if err != nil {
		fmt.Printf("Decode error: %v.", err)
		return
	}
	log.Printf("POST json: Attester=%s, Attestee=%s\n", params.Attester, params.Attestee)

	// save data to dag & db
	err = StoreMessage(&params)
	if err != nil {
		fmt.Printf("Save message error: %v.", err)
		return
	}

	/*// hex encode
	key_hex := make([]byte, hex.EncodedLen(len(k)))
	hex.Encode(key_hex, k)

	// return
	store_reply := types.StoreReply{
		Code: 0,
		Hash: fmt.Sprintf("0x%s", key_hex),
	}
	reply, _ := json.Marshal(store_reply)
	w.Write(reply)*/
}

// GetHandle process the 'get' request.
func GetHandle(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	// params
	var params types.GetReq
	err := decoder.Decode(&params)
	if err != nil {
		fmt.Printf("Get error: %v.", err)
		return
	}
	log.Printf("POST json: Key=%s\n", params.Key)

	// hex decode
	k := strings.TrimPrefix(params.Key, "0x")
	hash := make([]byte, hex.DecodedLen(len(k)))
	_, err = hex.Decode(hash, []byte(k))
	if err != nil {
		log.Fatal(err)
	}

	// get data from db
	value, err := db.Get([]byte(hash))
	if err != nil {
		log.Printf("Get error: %v.", err)
		return
	}
	log.Printf("Value = '%s'\n", value)

	// return
	get_reply := types.GetReply{
		Value: string(value),
	}
	reply, _ := json.Marshal(get_reply)
	w.Write(reply)
}
