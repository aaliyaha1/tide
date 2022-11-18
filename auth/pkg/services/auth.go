package services

import (
	"context"
	"fmt"
	"github.com/tide/auth/pkg/db"
	"github.com/tide/auth/pkg/models"
	"github.com/tide/auth/pkg/pb"
	"github.com/tide/auth/pkg/utils"
)

type Server struct {
	H   db.Handler
	Jwt utils.JwtWrapper
	pb.UnimplementedAuthServiceServer
}

func (s *Server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	funcName := "[ERROR] [Auth] [Register]"
	var user models.User

	// TODO
	result := s.H.DB.Where(&models.User{Account: req.Account}).First(&user)
	if result.Error == nil {
		errRes := fmt.Errorf(funcName+" %s", "account already exists")
		return &pb.RegisterResponse{
			Status: 120003,
			Error:  errRes.Error(),
		}, nil
	}

	user.Account = req.Account
	user.Password = utils.HashPassword(req.Password)

	s.H.DB.Create(&user)

	return &pb.RegisterResponse{
		Status: 120001,
	}, nil
}

func (s *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	funcName := "[ERROR] [Auth] [Login]"
	var user models.User

	if result := s.H.DB.Where(&models.User{Account: req.Account}).First(&user); result.Error != nil {
		errRes := fmt.Errorf(funcName+" %s", "user not found")
		return &pb.LoginResponse{
			Status: 120003,
			Error:  errRes.Error(),
		}, nil
	}

	match := utils.CheckPasswordHash(req.Password, user.Password)

	if !match {
		errRes := fmt.Errorf(funcName+" %s", "user not found")
		return &pb.LoginResponse{
			Status: 120004,
			Error:  errRes.Error(),
		}, nil
	}

	token, _ := s.Jwt.GenerateToken(user)

	return &pb.LoginResponse{
		Status: 120001,
		Token:  token,
	}, nil
}

func (s *Server) Validate(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	funcName := "[ERROR] [Auth] [Validate]"
	claims, err := s.Jwt.ValidateToken(req.Token)

	if err != nil {
		errRes := fmt.Errorf(funcName+" %s", err)
		return &pb.ValidateResponse{
			Status: 120003,
			Error:  errRes.Error(),
		}, nil
	}

	var user models.User

	if result := s.H.DB.Where(&models.User{Account: claims.Account}).First(&user); result.Error != nil {
		errRes := fmt.Errorf(funcName+" %s", "user not found")
		return &pb.ValidateResponse{
			Status: 120004,
			Error:  errRes.Error(),
		}, nil
	}

	return &pb.ValidateResponse{
		Status:  120001,
		UserId:  user.Id,
		Account: user.Account,
	}, nil
}
