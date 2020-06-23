package server

import (
	"encoding/json"
	"fmt"
	"github.com/triasteam/go-streamnet/types"
	"net/http"
)

func Start() {
	// http server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, "hello, streamnet-go") })
	http.HandleFunc("/save", SaveHandle)
	http.HandleFunc("/get", GetHandle)

	http.ListenAndServe(":14700", nil)
}

func SaveHandle(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var params types.StoreData //map[string]string

	err := decoder.Decode(&params)
	if err != nil {
		fmt.Println("Save error: %v.", err)
		return
	}

	fmt.Printf("POST json: Attester=%s, Attestee=%s\n", params.Attester, params.Attestee)

	fmt.Fprintf(w, `{"code":0, "hash": }`)
}

func GetHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello, get")
	query := r.URL.Query()

	value := query.Get("hash")

	fmt.Printf("GET: value=%s\n", value)

	fmt.Fprintf(w, `{"code":0, "value": %s}`, value)
}
