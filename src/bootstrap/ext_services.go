package bootstrap

import (
	"check-price/src/infra/external"
	"check-price/src/infra/external/ghtk"
	"go.uber.org/fx"
)

func BuildExtServicesModules() fx.Option {
	return fx.Options(
		fx.Provide(external.NewBaseClient),
		fx.Provide(ghtk.NewGHTKExtService),
	)
}
