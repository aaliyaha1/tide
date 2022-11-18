package routes

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tide/engine/pkg/data/pb"
	"net/http"
)

func GetBalance(ctx *gin.Context, c pb.DataServiceClient) {
	account, _ := ctx.Get("account")
	// rpc 请求
	balance, err := c.GetBalance(ctx, &pb.GetBalanceRequest{
		Account: account.(string),
	})
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}
	var data BalanceData
	_ = json.Unmarshal(balance.Res, &data)
	ctx.JSON(http.StatusCreated, map[string]interface{}{
		"status": balance.Status,
		"data":   data,
		"err":    balance.Error,
	})
}

func GetVolume(ctx *gin.Context, c pb.DataServiceClient) {
	account, _ := ctx.Get("account")
	v, err := c.GetVolume(ctx, &pb.GetVolumeRequest{
		Account: account.(string),
	})
	fmt.Println(v, err)
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}
	var b VolumeData
	json.Unmarshal(v.Res, &b)
	ctx.JSON(http.StatusCreated, map[string]interface{}{
		"status": v.Status,
		"data":   b,
	})
}

func GetOrders(ctx *gin.Context, c pb.DataServiceClient) {
	account, _ := ctx.Get("account")
	os, err := c.GetOrders(ctx, &pb.GetOrdersRequest{
		Account: account.(string),
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}
	var o []OrderData
	json.Unmarshal(os.Res, &o)
	ctx.JSON(http.StatusCreated, map[string]interface{}{
		"status": os.Status,
		"data":   o,
	})
}
