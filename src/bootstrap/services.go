package bootstrap

import (
	"check-price/src/core/service"
	"go.uber.org/fx"
)

func BuildServiceModule() fx.Option {
	return fx.Options(
		fx.Provide(service.NewBaseService),
		fx.Provide(service.NewPriceService),
	)
}
