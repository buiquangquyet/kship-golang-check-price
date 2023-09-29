package domain

import (
	"check-price/src/common"
	"context"
)

type SettingShop struct {
}

type SettingShopRepo interface {
	GetByRetailerId(ctx context.Context, retailerId int64) ([]int64, *common.Error)
	GetEnableShopByRetailerId(ctx context.Context, retailerId int64) ([]int64, *common.Error)
	GetServiceExtraEnableShop(ctx context.Context, retailerId int64) (bool, *common.Error)
}

func (SettingShop) TableName() string {
	return "kship_setting_shops"
}
