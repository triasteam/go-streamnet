package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/triasteam/go-streamnet/abci/proto"
	"github.com/triasteam/go-streamnet/types"
	"google.golang.org/grpc"
)

func getRank(data []string, peroid uint32, numRank uint32) *types.Message {
	// create connection
	conn, err := grpc.Dial(address+":"+rpcPort, grpc.WithInsecure())
	if nil != conn {
		defer conn.Close()
	}

	if nil != err {
		fmt.Printf("Connect to grpc server failed: %s\n", err)
		return nil
	}

	client := proto.NewStreamnetServiceClient(conn)

	req := &proto.RequestGetNoderank{
		BlockHash: data,
		Period:    peroid,
		NumRank:   numRank,
	}

	message := types.Message{}
	message.Code = 0
	message.Timestamp = time.Now().Unix()
	message.Message = "Query node data successfully"

	result, err := client.GetNoderank(context.Background(), req)
	if err != nil {
		dt := types.DataTee{}
		teeScores := result.GetTeescore()
		newTeeScoreArray := make([]types.TeeScore, len(teeScores))
		for _, ts := range teeScores {
			newTs := types.TeeScore{}
			newTs.Attestee = ts.GetAttestee()
			newTs.Score = float64(ts.GetScore())
			newTeeScoreArray = append(newTeeScoreArray, newTs)
		}

		teeCtxes := result.GetTeectx()
		newTeeCtxArr := make([]types.TeeCtx, len(teeCtxes))
		for _, tc := range teeCtxes {
			newTc := types.TeeCtx{}
			newTc.Attestee = tc.GetAttestee()
			newTc.Attester = tc.GetAttester()
			if tc.GetScore() != "" {
				newTc.Score, _ = strconv.ParseFloat(tc.GetScore(), 64)
			}
			newTeeCtxArr = append(newTeeCtxArr, newTc)
		}
		dt.Teectx = newTeeCtxArr
		dt.Teescore = newTeeScoreArray
		message.Data = dt
	} else {
		message.Code = 1
		message.Message = "response is nil"
		fmt.Println("response is nil")
	}
	return &message
}

// QueryNodesHandle process the 'get' request.
func QueryNodesHandle(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	// params
	var params types.QueryNodeReq
	err := decoder.Decode(&params)
	if err != nil {
		fmt.Printf("Get error: %v.", err)
		return
	}
	log.Printf("POST json: %s\n", params)

	// get data from dag
	value := sn.Dag.GetTotalOrder()

	input := make([]string, 0)
	for _, b := range value {
		input = append(input, b.String())
	}

	response := getRank(input, params.Period, params.NumRank)

	message, err := json.Marshal(response)
	if err != nil {
		panic("query page rank error.")
	}

	// return
	// reply, _ := json.Marshal(message)
	w.Write(message)
}
