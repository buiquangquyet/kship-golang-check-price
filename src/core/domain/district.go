package domain

import (
	"check-price/src/common"
	"context"
)

type District struct {
}

type DistrictRepo interface {
	GetByKmsId(ctx context.Context, receiverLocationId int64) (*District, *common.Error)
	GetByKvId(ctx context.Context, receiverLocationId int64) (*District, *common.Error)
}
