package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/triasteam/go-streamnet/store"

	"github.com/triasteam/go-streamnet/types"
)

var (
	server *http.Server
	db     *store.Storage
)

func Start(store *store.Storage) {
	//TODO: find a better way to check whether server has started.
	if server != nil {
		log.Printf("Server already started.\n")
		return
	}

	// set db
	db = store

	// http server
	mux := http.NewServeMux()
	mux.Handle("/", &gsnHandler{})
	mux.HandleFunc("/save", SaveHandle)
	mux.HandleFunc("/get", GetHandle)

	server = &http.Server{
		Addr:    ":14700",
		Handler: mux,
		//WriteTimeout: time.Second * 3,
	}

	log.Fatal(server.ListenAndServe())
}

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

func SaveHandle(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var params types.StoreData //map[string]string

	err := decoder.Decode(&params)
	if err != nil {
		fmt.Printf("Save error: %v.", err)
		return
	}

	log.Printf("POST json: Attester=%s, Attestee=%s\n", params.Attester, params.Attestee)

	k, err := db.SaveValue([]byte(params.String()))
	if err != nil {
		log.Printf("Save data to database failed: %v\n", err)
		fmt.Fprintf(w, `{"code":-1, "hash": }`)
		return
	}

	fmt.Fprintf(w, `{"code":0, "hash": %v}`, k)
}

func GetHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello, get")
	query := r.URL.Query()

	value := query.Get("hash")

	fmt.Printf("GET: value=%s\n", value)

	fmt.Fprintf(w, `{"code":0, "value": %s}`, value)
}
