package router

import (
	"check-price/src/common/configs"
	"check-price/src/present/httpui/controllers"
	"check-price/src/present/httpui/middlewares"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type RoutersIn struct {
	fx.In
	Engine          *gin.Engine
	PriceController *controllers.PriceController
}

func RegisterRouters(in RoutersIn) {
	group := in.Engine.Group(configs.Get().Server.Prefix)
	group.GET("/ping", middlewares.HealthCheckEndpoint)
	group.Use(middlewares.Authenticate())
	{
		group.POST("/check-price/:client-code", in.PriceController.Get)
	}
}
