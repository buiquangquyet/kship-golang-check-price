package bootstrap

import (
	"check-price/src/present/httpui/middlewares"
	"go.uber.org/fx"
)

func BuildAuthModules() fx.Option {
	return fx.Options(
		fx.Provide(middlewares.NewAuthMiddleware),
	)
}
