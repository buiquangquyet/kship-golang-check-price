package domain

import (
	"check-price/src/common"
	"context"
)

type Ward struct {
}

type WardRepo interface {
	GetByKmsId(ctx context.Context, senderWardId int64) (*Ward, *common.Error)
	GetByKvId(ctx context.Context, senderWardId int64) (*Ward, *common.Error)
}

func (Ward) TableName() string {
	return "wards"
}
