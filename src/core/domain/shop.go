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
	GHTKUsername   string    `json:"ghtk_username"`
	GHTKPassword   string    `json:"ghtk_password"`
	JtCustomerId   string    `json:"jt_customer_id"`
	GHNGWShopId    string    `json:"ghnfw_shop_id"`
	GHNFWPhone     string    `json:"ghnfw_phone"`
	UsernameBestFw string    `json:"username_bestfw"`
	PasswordBestFw string    `json:"password_bestfw"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	DeletedAt      time.Time `json:"deleted_at"`
}

type ShopRepo interface {
	GetByRetailerId(ctx context.Context, retailerId string) (*Shop, *common.Error)
	GetByCode(ctx context.Context, shopCode string) (*Shop, *common.Error)
}
