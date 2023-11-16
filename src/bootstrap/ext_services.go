package bootstrap

import (
	"check-price/src/infra/external"
	"check-price/src/infra/external/ahamove"
	"check-price/src/infra/external/ghtk"
	"check-price/src/infra/external/voucher"
	"go.uber.org/fx"
)

func BuildExtServicesModules() fx.Option {
	return fx.Options(
		fx.Provide(external.NewBaseClient),
		fx.Provide(ghtk.NewGHTKExtService),
		fx.Provide(voucher.NewVoucherExtService),
		fx.Provide(ahamove.NewAhaMoveExtService),
	)
}
