package server

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/triasteam/go-streamnet/config"

	"github.com/triasteam/go-streamnet/abci/proto"

	"github.com/triasteam/go-streamnet/types"
	"google.golang.org/grpc"
)

func callApp(data string) string {
	// create connection
	conn, err := grpc.Dial(address+":"+rpcPort, grpc.WithInsecure())
	if nil != conn {
		defer conn.Close()
	}

	if nil != err {
		fmt.Printf("Connect to grpc server failed: %s\n", err)
		return ""
	}

	client := proto.NewStreamnetServiceClient(conn)

	req := &proto.RequestStoreBlock{
		BlockInfo: data,
	}

	result, _ := client.StoreBlock(context.Background(), req)
	if nil != result {
		fmt.Printf("%s \n", result.Data)
	} else {
		fmt.Println("response is nil")
	}

	return result.Data
}

func StoreMessage(message *types.StoreData) ([]byte, error) {
	// Tipselection
	txsToApprove := sn.Tips.GetTransactionsToApprove(15, types.NilHash)
	if txsToApprove.Index(0) == types.NilHash || txsToApprove.Index(1) == types.NilHash {
		// Using genesis.
		txsToApprove.RemoveAtIndex(0)
		txsToApprove.RemoveAtIndex(0)
		txsToApprove.Append(config.GenesisTrunk)
		txsToApprove.Append(config.GenesisBranch)
	}

	// Transaction
	tx := types.Transaction{}
	tx.Init(txsToApprove)

	// grpc
	grpcResult := callApp(message.String())
	h := types.NewHashHex(grpcResult)
	tx.DataHash = h
	log.Printf("Grpc result: %s\n", h)

	// todo: POW

	// timestamp
	tx.Timestamp = time.Now()

	// tx hash
	txBytes, err := tx.Bytes()
	if err != nil {
		log.Printf("Transaction to bytes failed: %s\n", err)
		return nil, err
	}
	txHash := types.Sha256(txBytes)

	// Save to dag
	err = sn.Dag.Add(txHash, &tx)
	if err != nil {
		log.Printf("Dag add tx failed: %s\n", err)
		return nil, err
	}

	// Save to db
	k, err := sn.Store.SaveValue(txBytes)
	if err != nil {
		log.Printf("Save data to database failed: %v\n", err)
	}
	log.Printf("Store to database successed!\n")

	// todo: broadcast to neighbors.

	return k, err
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
	key, err := StoreMessage(&params)
	if err != nil {
		fmt.Printf("Save message error: %v.", err)
		return
	}

	// hex encode
	key_hex := make([]byte, hex.EncodedLen(len(key)))
	hex.Encode(key_hex, key)

	// return
	store_reply := types.StoreReply{
		Code: 0,
		Hash: fmt.Sprintf("0x%s", key_hex),
	}
	reply, _ := json.Marshal(store_reply)
	w.Write(reply)
}
