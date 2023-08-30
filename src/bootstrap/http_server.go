package bootstrap

import (
	"check-price/src/common/configs"
	"check-price/src/common/log"
	"check-price/src/present/httpui/router"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func BuildHTTPServerModule() fx.Option {
	return fx.Options(
		fx.Provide(gin.New),
		fx.Invoke(router.RegisterHandler),
		fx.Invoke(router.RegisterRouters),
		fx.Invoke(NewHttpServer),
	)
}

func NewHttpServer(lc fx.Lifecycle, engine *gin.Engine) {
	logger := log.GetLogger().GetZap()
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := engine.Run(fmt.Sprintf(":%s", configs.Get().Server.Address)); err != nil {
					logger.Fatalf("Cannot start application due by error [%v]", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Infof("Stopping HTTP server")
			return nil
		},
	})
}
