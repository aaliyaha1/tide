package main

import (
	"fmt"
	"github.com/tide/data/pkg/config"
	"github.com/tide/data/pkg/db/mongo_v0.1/mongo"
	"github.com/tide/data/pkg/pb"
	"github.com/tide/data/pkg/services"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	c, err := config.LoadConfig()

	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	h, _ := mongo.NewMongoDBWithoutBytes(c.DBUrl, c.DbAuthdb, c.DbUser, c.DbPassword, c.DbTable, "Account_Balance")

	lis, err := net.Listen("tcp", c.Port)

	if err != nil {
		log.Fatalln("Failed to listing:", err)
	}

	fmt.Println("===== [ TIDE Data SVC Running! ] ======")
	fmt.Printf("===== [     Port on%s      ] ======", c.Port)

	s := services.Server{
		D: h,
	}

	grpcServer := grpc.NewServer()

	pb.RegisterDataServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}
}
