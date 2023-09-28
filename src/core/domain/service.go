package domain

import (
	"check-price/src/common"
	"context"
)

type Service struct {
	Code             string
	Value            string
	OnBoardingStatus bool
}

type ServiceRepo interface {
	GetServicesPluckCodeByClientCode(ctx context.Context, clientCode string) ([]string, *common.Error)
	GetByServiceCode(ctx context.Context, clientCode string) ([]*Service, *common.Error)
}

func (Service) TableName() string {
	return "services"
}
