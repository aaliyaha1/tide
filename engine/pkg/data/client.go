package data

import (
	"fmt"
	"github.com/tide/engine/pkg/config"
	"github.com/tide/engine/pkg/data/pb"
	"google.golang.org/grpc"
)

type ServiceClient struct {
	Client pb.DataServiceClient
}

func InitServiceClient(c *config.Config) pb.DataServiceClient {
	// using WithInsecure() because no SSL running
	cc, err := grpc.Dial(c.DataSvcUrl, grpc.WithInsecure())

	if err != nil {
		fmt.Println("Could not connect:", err)
	}

	return pb.NewDataServiceClient(cc)
}
