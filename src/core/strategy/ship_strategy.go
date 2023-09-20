package strategy

import (
	"context"
	"go.uber.org/fx"
)

type ShipStrategy interface {
	Code() string
	GetMultiplePriceV3(ctx context.Context)
}

func ProvideEInvoiceStrategyFilterStrategy(constructor interface{}) fx.Option {
	return fx.Provide(fx.Annotated{
		Group:  "ship_strategy",
		Target: constructor,
	})
}

type ShipStrategyFilterResolverIn struct {
	fx.In
	Strategies []ShipStrategy `group:"ship_strategy"`
}

type ShipStrategyFilterResolver struct {
	MapStrategies map[string]ShipStrategy
}

func NewShipStrategyFilterResolver(in ShipStrategyFilterResolverIn) *ShipStrategyFilterResolver {
	mapStrategies := make(map[string]ShipStrategy)
	for _, strategy := range in.Strategies {
		mapStrategies[strategy.Code()] = strategy
	}
	return &ShipStrategyFilterResolver{MapStrategies: mapStrategies}
}

func (s ShipStrategyFilterResolver) Resolve(code string) (ShipStrategy, bool) {
	if _, exist := s.MapStrategies[code]; !exist {
		return nil, false
	}
	return s.MapStrategies[code], true
}
