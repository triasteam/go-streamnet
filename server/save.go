package server

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/triasteam/go-streamnet/abci/proto"

	"github.com/triasteam/go-streamnet/streamnet"
	"github.com/triasteam/go-streamnet/types"
	"google.golang.org/grpc"
)

// GlobalData used for translating send data
var GlobalData streamnet.StreamNet

func callApp(data string) string {
	// create connection
	conn, err := grpc.Dial(address+rpcPort, grpc.WithInsecure())
	if nil != conn {
		defer conn.Close()
	}

	if err != nil {
		log.Printf("Connect to grpc server failed: %s\n", err)
		return ""
	}

	client := proto.NewStreamnetServiceClient(conn)

	req := &proto.RequestStoreBlock{
		BlockInfo: data,
	}

	result, err := client.StoreBlock(context.Background(), req)
	if err != nil {
		log.Println("Response is nil!")
	}

	return result.Data
}

// StoreMessage ...
func StoreMessage(message *types.StoreData) ([]byte, error) {
	// Tipselection
	txsToApprove := sn.Tips.GetTransactionsToApprove(15, types.NilHash)

	// grpc
	grpcResult := callApp(message.String())
	h := types.NewHashHex(grpcResult)

	log.Printf("\n Grpc result: %s\n", h)

	// Transaction
	tx := types.Transaction{}
	tx.Init(txsToApprove, h)

	// todo: POW

	// tx hash
	txBytes, err := tx.Bytes()
	if err != nil {
		log.Printf("Transaction to bytes failed: %s\n", err)
		return nil, err
	}
	txHash := types.Sha256(txBytes)
	hashBytes := txHash.Bytes()

	// Save to dag
	err = sn.Dag.Add(txHash, &tx)
	if err != nil {
		log.Printf("Dag add tx failed: %s\n", err)
		return nil, err
	}

	// Save to db
	err = sn.Store.Save(hashBytes, txBytes)
	if err != nil {
		log.Printf("Save data to database failed: %v\n", err)
	}
	log.Printf("Store to database successed!\n")

	// broadcast to neigbors
	sendData := &types.SendData{}
	// convert to json string
	msg, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}
	sendData.Data = string(msg)
	sendData.Parent = txsToApprove.Index(0).String()
	sendData.Reference = txsToApprove.Index(1).String()
	sendData.Timestamp = tx.Timestamp
	msg, err = json.Marshal(sendData)
	if err == nil {
		broadcast(string(msg))
	}

	return hashBytes, err
}

// SaveHandle process the 'save' request.
func SaveHandle(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	// params
	var params types.StoreData
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Decode error: %v.", err)
		return
	}
	log.Printf("POST json: Attester=%s, Attestee=%s\n", params.Attester, params.Attestee)

	// save data to dag & db
	key, err := StoreMessage(&params)
	if err != nil {
		log.Printf("Save message error: %v.", err)
		return
	}

	// hex encode
	keyHex := make([]byte, hex.EncodedLen(len(key)))
	hex.Encode(keyHex, key)

	// return
	storeReply := types.StoreReply{
		Code: 0,
		Hash: fmt.Sprintf("0x%s", string(keyHex)),
	}
	reply, err := json.Marshal(storeReply)
	if err == nil {
		w.Write(reply)
	}
}

// OnReceived means after received message from neigbors, first will poll parent and reference,
// then getting origin message and store it to local.
// "h" means local storage key, it's different from any other neigbors.
func OnReceived(message string) error {
	var data types.SendData
	err := json.Unmarshal([]byte(message), &data)
	if err != nil {
		panic(err)
	}
	txsToApprove := types.List{}
	txsToApprove.Append(types.NewHashHex(data.Parent))
	txsToApprove.Append(types.NewHashHex(data.Reference))

	grpcResult := callApp(fmt.Sprintf("%v", data.Data))
	h := types.NewHashHex(grpcResult)
	log.Printf("\n Grpc result: %s\n", h)

	// Transaction
	tx := types.Transaction{}
	tx.Trunk = txsToApprove.Index(0)
	tx.Branch = txsToApprove.Index(1)
	// timestamp
	tx.Timestamp = data.Timestamp
	tx.DataHash = h

	// todo: POW

	// tx hash
	txBytes, err := tx.Bytes()
	if err != nil {
		log.Printf("Transaction to bytes failed: %s\n", err)
		return err
	}

	txHash := types.Sha256(txBytes)
	hashBytes := txHash.Bytes()

	if sn.Dag.Contains(txHash) {
		fmt.Printf("tx [%s] already exist.", txHash.String())
		return nil
	}

	// Save to dag
	err = sn.Dag.Add(txHash, &tx)
	if err != nil {
		log.Printf("Dag add tx failed: %s\n", err)
		return err
	}

	// Save to db
	err = sn.Store.Save(hashBytes, txBytes)
	if err != nil {
		log.Printf("Save data to database failed: %v\n", err)
	}
	log.Printf("Store to database successed!\n")

	broadcast(message)

	return nil
}

func broadcast(message string) error {
	// broadcast to other nodes
	sn.Network.Broadcast(message)
	return nil
}
