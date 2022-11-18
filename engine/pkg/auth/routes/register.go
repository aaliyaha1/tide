package routes

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegisterRequestBody struct {
	Account  string `json:"account"`
	Password string `json:"password"`
	Symbol   string `json:"symbol"`
	Data     string `json:"data"`
}

func Register(ctx *gin.Context, c pb.AuthServiceClient) {
	body := RegisterRequestBody{}

	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	res, err := c.Register(context.Background(), &pb.RegisterRequest{
		Account:  body.Account,
		Password: body.Password,
		Symbol:   body.Symbol,
		Data:     body.Data,
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(int(res.Status), &res)
}
