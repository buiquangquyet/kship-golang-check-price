package domain

import (
	"check-price/src/common"
	"context"
)

type Service struct {
}

type ServiceRepo interface {
	GetServicesPluckCodeByClientCode(ctx context.Context, clientCode string) ([]string, *common.Error)
}
