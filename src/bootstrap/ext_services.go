package bootstrap

import (
	"check-price/src/infra/external"
	"check-price/src/infra/external/ahamove"
	"check-price/src/infra/external/aieliminating"
	"check-price/src/infra/external/ghtk"
	"check-price/src/infra/external/voucher"
	"go.uber.org/fx"
)

func BuildExtServicesModules() fx.Option {
	return fx.Options(
		fx.Provide(external.NewBaseClient),

		fx.Provide(voucherext.NewService),
		fx.Provide(aieliminating.NewService),

		fx.Provide(ahamoveext.NewService),
		fx.Provide(ghtkext.NewService),
	)
}
