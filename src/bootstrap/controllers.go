package bootstrap

import (
	"check-price/src/present/httpui/controllers"
	"check-price/src/present/httpui/validator"
	"go.uber.org/fx"
)

func BuildControllerModule() fx.Option {
	return fx.Options(
		fx.Provide(controllers.NewBaseController),
		fx.Provide(controllers.NewPriceController),
		fx.Provide(controllers.NewLogController),
	)
}

func BuildValidator() fx.Option {
	return fx.Options(
		fx.Provide(validator.NewValidator),
		fx.Invoke(validator.RegisterValidations),
	)
}
