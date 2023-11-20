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
	AuthMiddleware  *middlewares.AuthMiddleware
}

func RegisterRouters(in RoutersIn) {
	group := in.Engine.Group(configs.Get().Server.Prefix)
	group.GET("/ping", middlewares.HealthCheckEndpoint)
	group.Use(in.AuthMiddleware.Authenticate())
	{
		group.POST("/:client", in.PriceController.GetPrice)
	}
}
