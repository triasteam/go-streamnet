package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/triasteam/go-streamnet/abci/proto"
	"github.com/triasteam/go-streamnet/types"
	"github.com/triasteam/go-streamnet/utils"
	"google.golang.org/grpc"
)

func getRank(data []string, peroid uint32, numRank uint32) *types.Message {
	// create connection
	conn, err := grpc.Dial(address+rpcPort, grpc.WithInsecure())
	if nil != conn {
		defer conn.Close()
	}

	if nil != err {
		log.Printf("Connect to grpc server failed: %s\n", err)
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
	if err == nil {
		dt := types.DataTee{}
		teeScores := result.GetTeescore()
		newTeeScoreArray := make([]types.TeeScore, 0)
		for _, ts := range teeScores {
			newTs := types.TeeScore{}
			newTs.Attestee = ts.GetAttestee()
			newTs.Score = float64(ts.GetScore())
			newTeeScoreArray = append(newTeeScoreArray, newTs)
		}

		teeCtxes := result.GetTeectx()
		newTeeCtxArr := make([]types.TeeCtx, 0)
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
		log.Println("response is nil", err)
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
		log.Printf("Get error: %v.", err)
		return
	}
	log.Printf("POST json: %s\n", params)

	// get data from dag
	value := sn.Dag.GetTotalOrder()

	// FIXME while change genesis forward, this method would be changed.
	// Get peroid using paging algorithm

	data4Period := utils.Paging(value, params.Period, int(params.NumRank))
	// get tx's dataHash from db
	input := make([]string, 0)
	for _, hash := range data4Period {
		txBytes, err := sn.Store.Get(hash.Bytes())
		if err != nil {
			log.Panicln("Get data from database failed!", err)
		}
		tx := types.TransactionFromBytes(txBytes)
		input = append(input, tx.DataHash.String())
	}

	response := getRank(input, uint32(params.Period), params.NumRank)

	message, err := json.Marshal(response)
	if err != nil {
		panic("query page rank error.")
	}

	// return
	// reply, _ := json.Marshal(message)
	w.Write(message)
}
