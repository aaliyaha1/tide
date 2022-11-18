package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/globalsign/mgo/bson"
	"github.com/tide/data/pkg/db/mongo_v0.1/mongo"
	"github.com/tide/data/pkg/models"
	"github.com/tide/data/pkg/pb"
)

type Server struct {
	D       *mongo.MongoDB
	DbTable string
	pb.UnimplementedDataServiceServer
}

// GetBalance 获取账号余额
func (s *Server) GetBalance(ctx context.Context, req *pb.GetBalanceRequest) (*pb.GetBalanceResponse, error) {
	funcName := "[ERROR] [Data] [GetBalance]"
	// 获取当前用户信息
	account := req.Account
	if account == "" {
		errRes := fmt.Errorf(funcName+" %s", "account is empty")
		return &pb.GetBalanceResponse{
			Status: 110101,
			Error:  errRes.Error(),
		}, nil
	}

	var balances []models.BalanceData
	err := s.D.FindAll(s.DbTable, "Account_Balance", "_id", bson.M{"account": account}, nil, &balances)
	if err != nil {
		errRes := fmt.Errorf(funcName+" %s", err)
		return &pb.GetBalanceResponse{
			Status: 110003,
			Error:  errRes.Error(),
		}, nil
	}

	// 204
	if len(balances) == 0 {
		errRes := fmt.Errorf(funcName+" %s", "balances is empty")
		return &pb.GetBalanceResponse{
			Status: 110002,
			Error:  errRes.Error(),
		}, nil
	}

	balance, _ := json.Marshal(balances[len(balances)-1])
	return &pb.GetBalanceResponse{
		Status: 110001,
		Res:    balance,
	}, nil
}

// GetVolume 获取账户刷量情况
func (s *Server) GetVolume(ctx context.Context, req *pb.GetVolumeRequest) (*pb.GetVolumeResponse, error) {
	funcName := "[ERROR] [Data] [GetVolume]"
	account := req.Account
	if account == "" {
		errRes := fmt.Errorf(funcName+" %s", "account is empty")
		return &pb.GetVolumeResponse{
			Status: 110101,
			Error:  errRes.Error(),
		}, nil
	}

	var volumes []models.VolumeData
	err := s.D.FindAll(s.DbTable, "Account_Volume", "_id", bson.M{"account": account}, nil, &volumes)
	if err != nil {
		errRes := fmt.Errorf(funcName+" %s", err)
		return &pb.GetVolumeResponse{
			Status: 110003,
			Error:  errRes.Error(),
		}, nil
	}

	if len(volumes) == 0 {
		errRes := fmt.Errorf(funcName+" %s", "res is empty")
		return &pb.GetVolumeResponse{
			Status: 110002,
			Error:  errRes.Error(),
		}, nil
	}
	volume, _ := json.Marshal(volumes[len(volumes)-1])
	return &pb.GetVolumeResponse{
		Status: 110001,
		Res:    volume,
	}, nil
}

func (s *Server) GetOrders(ctx context.Context, req *pb.GetOrdersRequest) (*pb.GetOrdersResponse, error) {
	funcName := "[ERROR] [Data] [GetOrders]"
	account := req.Account
	if account == "" {
		errRes := fmt.Errorf(funcName+" %s", "account is empty")
		return &pb.GetOrdersResponse{
			Status: 110101,
			Error:  errRes.Error(),
		}, nil
	}

	var orders []models.OrderData
	err := s.D.FindLimit("cctx_db_test", "test_orders", 100, nil, nil, &orders)
	if err != nil {
		return &pb.GetOrdersResponse{
			Status: 110003,
			Error:  err.Error(),
		}, nil
	}

	if len(orders) == 0 {
		errRes := fmt.Errorf(funcName+" %s", "res is empty")
		return &pb.GetOrdersResponse{
			Status: 110002,
			Error:  errRes.Error(),
		}, nil
	}

	os, _ := json.Marshal(orders)
	return &pb.GetOrdersResponse{
		Status: 110001,
		Res:    os,
	}, nil
}
