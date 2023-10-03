package domain

import (
	"check-price/src/common"
	"context"
	"time"
)

type SettingShop struct {
	Id         int64     `json:"id"`
	RetailerId int64     `json:"retailer_id"`
	ModelType  string    `json:"model_type"`
	ModelId    int       `json:"model_id"`
	Value      string    `json:"value"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type SettingShopRepo interface {
	GetByRetailerId(ctx context.Context, retailerId int64) ([]int64, *common.Error)
	GetEnableShopByRetailerId(ctx context.Context, retailerId int64) ([]int64, *common.Error)
	GetServiceExtraEnableShop(ctx context.Context, retailerId int64) (bool, *common.Error)
}

func (SettingShop) TableName() string {
	return "kship_setting_shops"
}
