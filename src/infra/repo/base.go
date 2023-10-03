package repo

import (
	"check-price/src/common"
	"context"
	"gorm.io/gorm"
)

type baseRepo struct {
	db *gorm.DB
}

func NewBaseRepo(db *gorm.DB) *baseRepo {
	return &baseRepo{
		db: db,
	}
}

func (b *baseRepo) returnError(ctx context.Context, err error) *common.Error {
	return common.ErrSystemError(ctx, err.Error()).SetSource(common.SourceInfraService)
}
