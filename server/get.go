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

// func getRank(data []string) *proto.ResponseGetNoderank {
// 	// create connection
// 	conn, err := grpc.Dial(address+":"+rpcPort, grpc.WithInsecure())
// 	if nil != conn {
// 		defer conn.Close()
// 	}

// 	if nil != err {
// 		fmt.Printf("Connect to grpc server failed: %s\n", err)
// 		return nil
// 	}

// 	client := proto.NewStreamnetServiceClient(conn)

// 	req := &proto.RequestGetNoderank{
// 		BlockHash: data,
// 	}

// 	result, _ := client.GetNoderank(context.Background(), req)
// 	if nil != result {
// 		fmt.Printf("%s \n", result.GetTeectx)
// 	} else {
// 		fmt.Println("response is nil")
// 	}
// 	return result
// }

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
	// if err != nil {
	// 	log.Printf("Get error: %v.", err)
	// 	return
	// }
	// log.Printf("Value = '%s'\n", value)

	// get data from dag
	// value := sn.Dag.GetTotalOrder()

	// input := make([]string, len(value))
	// for _, b := range value {
	// 	b.String()
	// 	input = append(input, b.String())
	// }

	// response := getRank(input)

	message, err := json.Marshal(value)
	if err != nil {
		panic("query page rank error.")
	}
	// return
	get_reply := types.GetReply{
		Value: string(message),
	}
	reply, _ := json.Marshal(get_reply)
	w.Write(reply)
}
