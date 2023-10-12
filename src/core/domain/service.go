package domain

import (
	"check-price/src/common"
	"check-price/src/core/enums"
	"context"
	"time"
)

type Service struct {
	Id               int64     `json:"id"`
	ClientId         int64     `json:"client_id"`
	Code             string    `json:"code"`
	Name             string    `json:"name"`
	Type             int       `json:"type"`
	Tag              string    `json:"tag"`
	Status           int       `json:"status"`
	Description      string    `json:"description"`
	GroupId          int       `json:"group_id"`
	ViewType         string    `json:"view_type"`
	Value            string    `json:"value"`
	Order            int       `json:"order"`
	ClientsPossible  string    `json:"clients_possible"`
	IsDefault        int       `json:"is_default"`
	ServicesPossible string    `json:"services_possible"`
	CitiesPossible   string    `json:"cities_possible"`
	OrderFacebook    int       `json:"order_facebook"`
	OrderMobile      int       `json:"order_mobile"`
	ClientsWhiteList string    `json:"clients_white_list"`
	OnBoardingStatus int       `json:"on_boarding_status"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	DeletedAt        time.Time `json:"deleted_at"`
}

type ServiceRepo interface {
	GetByClientId(ctx context.Context, typeService enums.TypeService, status int, clientId int64) ([]*Service, *common.Error)
	GetByClientCode(ctx context.Context, typeService enums.TypeService, status int, clientCode string) ([]*Service, *common.Error)
}

func (Service) TableName() string {
	return "services"
}
