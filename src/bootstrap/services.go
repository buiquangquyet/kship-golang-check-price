package bootstrap

import (
	"check-price/src/core/service"
	"check-price/src/core/strategy"
	"go.uber.org/fx"
)

func BuildServiceModule() fx.Option {
	return fx.Options(
		fx.Provide(service.NewBaseService),
		fx.Provide(service.NewPriceService),

		fx.Provide(strategy.NewShipStrategyFilterResolver),
		strategy.ProvideShipStrategyFilterStrategy(strategy.NewGHTKStrategy),
	)
}
