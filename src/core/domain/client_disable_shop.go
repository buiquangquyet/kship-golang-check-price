package domain

import (
	"check-price/src/common"
	"context"
)

type ClientDisableShop struct {
}

type ClientDisableShopRepo interface {
	GetByRetailerId(ctx context.Context, retailerId int64) ([]int64, *common.Error)
}
