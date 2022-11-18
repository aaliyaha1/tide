package order

import (
	"fmt"
	"github.com/tide/engine/pkg/config"
	"github.com/tide/engine/pkg/order/pb"
	"google.golang.org/grpc"
)

type ServiceClient struct {
	Client pb.OrderServiceClient
}

func InitServiceClient(c *config.Config) pb.OrderServiceClient {
	cc, err := grpc.Dial(c.OrderSvcUrl, grpc.WithInsecure())

	if err != nil {
		fmt.Println("Could not connect:", err)
	}

	return pb.NewOrderServiceClient(cc)
}
