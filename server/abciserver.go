package server

import (
	"reflect"

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
		var dest = &pb.NodeRankTeescore{}
		copyStruct(scoreUnit, dest)
		respTeescore[i] = dest
	}

	for i, ctxUnit := range teectx {
		var dest = &pb.NodeRankTeectx{}
		copyStruct(ctxUnit, dest)
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

func copyStruct(src, dst interface{}) { // struct copy
	sval := reflect.ValueOf(src).Elem()
	dval := reflect.ValueOf(dst).Elem()

	for i := 0; i < sval.NumField(); i++ {
		value := sval.Field(i)
		name := sval.Type().Field(i).Name

		dvalue := dval.FieldByName(name)
		if dvalue.IsValid() == false {
			continue
		}
		dvalue.Set(value)
	}
}
