package router

import (
	"check-price/src/common/configs"
	"check-price/src/present/httpui/middlewares"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func RegisterHandler(engine *gin.Engine) {
	// recovery
	engine.Use(middlewares.Recovery())
	//tracer
	engine.Use(otelgin.Middleware(configs.Get().Server.Name))
	engine.Use(middlewares.Tracer())
	// log middleware
	engine.Use(middlewares.Log())

	engine.Use(cors.AllowAll())
}
