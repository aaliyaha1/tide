package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tide/engine/pkg/auth"
	"github.com/tide/engine/pkg/config"
	"github.com/tide/engine/pkg/data"
	"github.com/tide/engine/pkg/order"
	"github.com/tide/engine/pkg/tools/log"
)

func init() {
	var param []byte
	logger, err := log.NewLoggerByLogrus(param)
	if err != nil {
		fmt.Printf("[ERROR] [Main] init log by Logrus failed , err: %s\n", err)
		return
	}
	log.InitGlobal(logger)
}

func main() {
	c, err := config.LoadConfig()

	if err != nil {
		log.Error(err)
		return
	}

	fmt.Println("===== [ TIDE WEB SVC Running! ] ======")
	fmt.Printf("===== [      Port on%s     ] ======\n", c.Port)

	r := gin.Default()
	r.Use(auth.Cors())

	authSvc := *auth.RegisterRoutes(r, &c)
	order.RegisterRoutes(r, &c, &authSvc)
	data.RegisterRoutes(r, &c, &authSvc)

	go func() {
		r.Run(":3001")
	}()
	select {}
}
