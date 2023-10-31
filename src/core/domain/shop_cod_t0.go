package domain

import (
	"check-price/src/common"
	"context"
	"time"
)

type ShopCodT0 struct {
	Id        int64     `json:"id"`
	ShopId    int64     `json:"shop_id"`
	TimeStart time.Time `json:"time_start"`
	TimeEnd   time.Time `json:"time_end"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ShopCodT0Repo interface {
	GetByShopId(ctx context.Context, shopId int64) (*ShopCodT0, *common.Error)
}

func (ShopCodT0) TableName() string {
	return "shop_use_cod_t0_histories"
}
