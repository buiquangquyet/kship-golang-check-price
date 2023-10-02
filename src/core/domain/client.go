package domain

import (
	"check-price/src/common"
	"context"
	"time"
)

type Client struct {
	Id               int64     `json:"id"`
	Name             string    `json:"name"`
	Code             string    `json:"code"`
	Status           int       `json:"status"`
	Image            string    `json:"image"`
	ShopsPossible    string    `json:"shops_possible"`
	LinkTracking     string    `json:"link_tracking"`
	OnBoardingStatus int       `json:"on_boarding_status"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	DeletedAt        time.Time `json:"deleted_at"`
}

type ClientRepo interface {
	GetByCode(ctx context.Context, clientCode string) (*Client, *common.Error)
}

func (Client) TableName() string {
	return "clients"
}
