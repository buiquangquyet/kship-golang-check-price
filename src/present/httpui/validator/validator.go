package validator

import (
	"check-price/src/common/log"
	"github.com/go-playground/validator/v10"
)

func NewValidator() *validator.Validate {
	return validator.New()

}

func registerValidation(validator *validator.Validate, tag string, fn validator.Func) {
	if err := validator.RegisterValidation(tag, fn); err != nil {
		log.GetLogger().GetZap().Fatalf("Register custom validation %s failed with error: %s", tag, err.Error())
	}
	return
}

func registerStructValidation(validator *validator.Validate, fn validator.StructLevelFunc, in ...interface{}) {
	validator.RegisterStructValidation(fn, in...)
}

func RegisterValidations(validator *validator.Validate) {
}
