package server

import (
	"encoding/json"
	"fmt"

	pb "github.com/triasteam/go-streamnet/abci/proto"
	streamnet_service "github.com/triasteam/go-streamnet/service"
	golang_context "golang.org/x/net/context"
)

type abciServer struct{}

func NewAbciServer() *abciServer {
	return &abciServer{}
}

func (s *abciServer) StoreBlock(ctx golang_context.Context, req *pb.RequestStoreBlock) (*pb.ResponseStoreBlock, error) {
	service := streamnet_service.NewTransServer()
	res := service.StoreDagData(req.BlockInfo)
	response := &pb.ResponseStoreBlock{
		Code: 1,
		Log:  "success",
		Data: res,
	}
	return response, nil
}

func (s *abciServer) GetNoderank(ctx golang_context.Context, req *pb.RequestGetNoderank) (*pb.ResponseGetNoderank, error) {
	service := streamnet_service.NewTransServer()
	teescore, teectx, err := service.GetNodeRank(req.BlockHash, int(req.Duration), int(req.Period), int(req.NumRank))
	if err != nil {
		return nil, nil
	}

	respTeescore := make([]*pb.NodeRankTeescore, len(teescore))
	respTeectx := make([]*pb.NodeRankTeectx, len(teectx))

	for i, scoreUnit := range teescore {
		var dest pb.NodeRankTeescore
		str, err := json.Marshal(scoreUnit)
		if err != nil {
			fmt.Println(err)
			continue
		}
		json.Unmarshal(str, &dest)
		respTeescore[i] = &dest
	}

	for i, ctxUnit := range teectx {
		var dest = &pb.NodeRankTeectx{}
		str, err := json.Marshal(ctxUnit)
		if err != nil {
			fmt.Println(err)
			continue
		}
		json.Unmarshal(str, &dest)
		respTeectx[i] = dest
	}

	response := &pb.ResponseGetNoderank{
		Code:     1,
		Log:      "success",
		Teescore: respTeescore,
		Teectx:   respTeectx,
	}
	return response, nil
}
