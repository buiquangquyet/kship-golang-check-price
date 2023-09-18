package domain

import (
	"check-price/src/common"
	"context"
)

type ClientSettingShop struct {
	Id int64 `json:"id"`
}

type ClientSettingShopRepo interface {
	GetEnableShopByRetailerId(ctx context.Context, retailerId int64) ([]int64, *common.Error)
}
