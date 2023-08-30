package domain

import (
	"check-price/src/common"
	"context"
	"time"
)

type Shop struct {
	Id             int64     `json:"id"`
	VtpUsername    string    `json:"vtp_username"`
	VtpPassword    string    `json:"vtp_password"`
	VnpCmsCode     string    `json:"vnp_cms_code"`
	GhtkUsername   string    `json:"ghtk_username"`
	GhtkPassword   string    `json:"ghtk_password"`
	JtCustomerId   string    `json:"jt_customer_id"`
	GhnfwShopId    string    `json:"ghnfw_shop_id"`
	GhnfwPhone     string    `json:"ghnfw_phone"`
	UsernameBestfw string    `json:"username_bestfw"`
	PasswordBestfw string    `json:"password_bestfw"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	DeletedAt      time.Time `json:"deleted_at"`
}

type ShopRepo interface {
	GetByRetailerId(ctx context.Context, retailerId string) (*Shop, *common.Error)
	GetByCode(ctx context.Context, shopCode string) (*Shop, *common.Error)
}
