package domain

import (
	"check-price/src/common"
	"check-price/src/core/enums"
	"context"
	"time"
)

type SettingShop struct {
	Id         int64     `json:"id"`
	RetailerId int64     `json:"retailer_id"`
	ModelType  string    `json:"model_type"`
	ModelId    int64     `json:"model_id"`
	Value      string    `json:"value"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type SettingShopRepo interface {
	GetByRetailerId(ctx context.Context, modelType enums.ModelTypeSettingShop, retailerId int64) ([]*SettingShop, *common.Error)
	GetByRetailerIdAndModelId(ctx context.Context, modelType enums.ModelTypeSettingShop, retailerId int64, modelId int64) ([]*SettingShop, *common.Error)
}

func (SettingShop) TableName() string {
	return "kship_setting_shops"
}
