package bootstrap

import (
	"check-price/src/core/service"
	"check-price/src/core/strategy"
	"check-price/src/core/strategy/ahamove"
	"check-price/src/core/strategy/ghtk"
	"go.uber.org/fx"
)

func BuildServiceModule() fx.Option {
	return fx.Options(
		fx.Provide(service.NewBaseService),
		fx.Provide(service.NewValidateService),
		fx.Provide(service.NewPriceService),
		fx.Provide(service.NewExtraService),
		fx.Provide(service.NewVoucherService),
		fx.Provide(strategy.NewBaseStrategy),
		fx.Provide(strategy.NewShipStrategyFilterResolver),
		strategy.ProvideShipStrategyFilterStrategy(ghtk.NewStrategy),
		strategy.ProvideShipStrategyFilterStrategy(ahamove.NewStrategy),
	)
}
