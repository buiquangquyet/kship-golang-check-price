package domain

import (
	"check-price/src/common"
	"context"
	"time"
)

type ShopLevel struct {
	Id           int64     `json:"id"`
	Code         string    `json:"code"`
	Name         string    `json:"name"`
	Status       int       `json:"status"`
	BestUser     string    `json:"best_user"`
	BestPassword string    `json:"best_password"`
	GhnMarkup    float64   `json:"ghn_markup"`
	BestMarkup   float64   `json:"best_markup"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    time.Time `json:"deleted_at"`
}

type ShopLevelRepo interface {
	GetById(ctx context.Context, id int64) (*ShopLevel, *common.Error)
}

func (ShopLevel) TableName() string {
	return "shop_levels"
}
