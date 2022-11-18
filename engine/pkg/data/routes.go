package data

import (
	"github.com/gin-gonic/gin"
	"github.com/tide/engine/pkg/auth"
	"github.com/tide/engine/pkg/config"
	"github.com/tide/engine/pkg/data/routes"
)

func RegisterRoutes(r *gin.Engine, c *config.Config, authSvc *auth.ServiceClient) {
	a := auth.InitAuthMiddleware(authSvc)

	svc := &ServiceClient{
		Client: InitServiceClient(c),
	}

	routes := r.Group("/data")
	routes.Use(a.AuthRequired)
	routes.POST("/get_balance", svc.GetBalance)
	routes.POST("/get_volume", svc.GetVolume)
	routes.POST("/get_orders", svc.GetOrders)
}

func (svc *ServiceClient) GetBalance(ctx *gin.Context) {
	routes.GetBalance(ctx, svc.Client)
}

func (svc *ServiceClient) GetVolume(ctx *gin.Context) {
	routes.GetVolume(ctx, svc.Client)
}

func (svc *ServiceClient) GetOrders(ctx *gin.Context) {
	routes.GetOrders(ctx, svc.Client)
}
