package request

type Page struct {
	Page  int `form:"page" json:"page" validate:"required,gte=1"`
	Limit int `form:"limit" json:"limit" validate:"required,gte=1"`
}

type PageMerchant struct {
	*Page                   `validate:"dive"`
	*MerchantMessageGateway `validate:"dive"`
}
