package server

import (
	"fmt"
	"log"

	"github.com/triasteam/go-streamnet/types"
)

func StoreMessage(message *types.StoreData) error {
	// Tipselection
	txToApprove := GetTransactionToApproveTips(15)

	// Transaction

	// Hash

	// POW

	// Save to dag

	// Save to db
	k, err := db.SaveValue([]byte(message.String()))
	if err != nil {
		log.Printf("Save data to database failed: %v\n", err)
		fmt.Fprintf(w, `{"code":-1, "hash": }`)
		return
	}
}
