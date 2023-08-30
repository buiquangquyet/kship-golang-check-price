package repo

import (
	"check-price/src/common"
	"context"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type baseRepo struct {
	db    *gorm.DB
	cache redis.UniversalClient
}

func NewBaseRepo(db *gorm.DB, cache redis.UniversalClient) *baseRepo {
	return &baseRepo{
		db:    db,
		cache: cache,
	}
}

func (b *baseRepo) returnError(ctx context.Context, err error) *common.Error {
	return common.ErrSystemError(ctx, err.Error()).SetSource(common.SourceInfraService)
}
