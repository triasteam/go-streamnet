package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/triasteam/go-streamnet/types"
)

func StoreMessage(message *types.StoreData) error {
	// Tipselection
	// Check genesis
	txToApprove := sn.Tips.GetTransactionsToApprove(15, types.NilHash)

	// Transaction
	tx := types.Transaction{}
	tx.Init(txToApprove)

	// Hash

	// POW

	// Save to dag

	// Save to db
	/*k, err := db.SaveValue([]byte(message.String()))
	if err != nil {
		log.Printf("Save data to database failed: %v\n", err)
		fmt.Fprintf(w, `{"code":-1, "hash": }`)
		return
	}*/
	return nil
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
