package server

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/triasteam/go-streamnet/types"
)

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
	value, err := sn.Store.Get([]byte(hash))
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
