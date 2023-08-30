package constant

import "check-price/src/common/configs"

const (
	AppEnvDev  = "dev"
	AppEnvProd = "prod"

	TraceIdName = "trace_id"
	Merchant    = "merchant"
)

func IsProdEnv() bool {
	return configs.Get().Mode == AppEnvProd
}

func IsDevEnv() bool {
	return configs.Get().Mode == AppEnvDev
}
