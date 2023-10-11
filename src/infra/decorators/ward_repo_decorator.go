package decorators

import (
	"check-price/src/common"
	"check-price/src/core/domain"
	"check-price/src/infra/repo"
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

const (
	expirationWardByKmsId = 30 * 24 * time.Hour
	expirationWardByKvId  = 30 * 24 * time.Hour
)

type WardRepoDecorator struct {
	*baseDecorator
	cache    redis.UniversalClient
	wardRepo *repo.WardRepo
}

func NewWardRepoDecorator(base *baseDecorator, wardRepo *repo.WardRepo, cache redis.UniversalClient) domain.WardRepo {
	return &WardRepoDecorator{
		baseDecorator: base,
		cache:         cache,
		wardRepo:      wardRepo,
	}
}

func (w WardRepoDecorator) GetByKmsId(ctx context.Context, kmsId int64) (*domain.Ward, *common.Error) {
	key := w.genKeyCacheGetWardByKmsId(kmsId)
	var ward domain.Ward
	err := w.cache.Get(ctx, key).Scan(&ward)
	if err != nil {
		return &ward, nil
	}
	w.handleRedisError(ctx, err)
	wardDB, ierr := w.wardRepo.GetByKmsId(ctx, kmsId)
	if ierr != nil {
		return nil, ierr
	}
	go w.cache.Set(ctx, key, wardDB, expirationWardByKmsId)
	return wardDB, nil
}

func (w WardRepoDecorator) GetByKvId(ctx context.Context, kvId int64) (*domain.Ward, *common.Error) {
	key := w.genKeyCacheGetWardByKvId(kvId)
	var ward domain.Ward
	err := w.cache.Get(ctx, key).Scan(&ward)
	if err != nil {
		return &ward, nil
	}
	w.handleRedisError(ctx, err)
	wardDB, ierr := w.wardRepo.GetByKvId(ctx, kvId)
	if ierr != nil {
		return nil, ierr
	}
	go w.cache.Set(ctx, key, wardDB, expirationWardByKvId)
	return wardDB, nil
}
