package main

import (
	//"bytes"
	"encoding/json"
	"fmt"
	"os"

	//"io"
	"net/http"
	//cmd "github.com/triasteam/StreamNet-go/commands"
	"github.com/triasteam/go-streamnet/store"
)

type StoreData struct {
	Attester string
	Attestee string
	Score    float64
}

func SaveHandle(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var params StoreData //map[string]string

	err := decoder.Decode(&params)
	if err != nil {
		fmt.Println("Save error: %v.", err)
		return
	}

	fmt.Printf("POST json: username=%s, password=%s\n", params.Attester, params.Attestee)

	fmt.Fprintf(w, `{"code":0, "hash": }`)
}

func GetHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello, get")
	query := r.URL.Query()

	value := query.Get("hash")

	fmt.Printf("GET: value=%s\n", value)

	fmt.Fprintf(w, `{"code":0, "value": %s}`, value)
}

type StreamNet struct {
	store *store.Storage
}

var GlobalData StreamNet

func main() {
	// open DB
	store, err := store.Init("./db")
	if err != nil {
		fmt.Printf("Open database failed!")
		os.Exit(-1)
	}

	GlobalData.store = store

	// http server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, "hello, streamnet-go") })
	http.HandleFunc("/save", SaveHandle)
	http.HandleFunc("/get", GetHandle)

	http.ListenAndServe(":14700", nil)

	/*
		rootCmd := cmd.RootCmd
		rootCmd.AddCommand(
			cmd.InitFilesCmd,
		)

		// parse config.  examples: sng --mwm 1 -p 14700 &>  sng.log &
		// other parameters like '--enable-streaming-graph' '--entrypoint-selector-algorithm "KATZ"' '--tip-sel-algo "CONFLUX"' '--weight-calculation-algorithm "IN_MEM"'
		// will be considered later.

		// start server
		// Create & start node
		rootCmd.AddCommand(cmd.NewRunNodeCmd())

	*/
}
