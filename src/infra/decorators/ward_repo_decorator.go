package decorators

import (
	"check-price/src/common"
	"check-price/src/common/log"
	"check-price/src/core/domain"
	"check-price/src/infra/repo"
	"context"
	"encoding/json"
	"time"
)

const (
	expirationWardByKmsId = 30 * 24 * time.Hour
	expirationWardByKvId  = 30 * 24 * time.Hour
)

type WardRepoDecorator struct {
	*baseDecorator
	wardRepo *repo.WardRepo
}

func NewWardRepoDecorator(base *baseDecorator, wardRepo *repo.WardRepo) domain.WardRepo {
	return &WardRepoDecorator{
		baseDecorator: base,
		wardRepo:      wardRepo,
	}
}

func (w WardRepoDecorator) GetByKmsId(ctx context.Context, kmsId int64) (*domain.Ward, *common.Error) {
	key := w.genKeyCacheGetWardByKmsId(kmsId)
	var ward domain.Ward
	val, err := w.get(ctx, key).Result()
	w.handleRedisError(ctx, err)
	if err == nil {
		err = json.Unmarshal([]byte(val), &ward)
		if err == nil {
			return &ward, nil
		}
		log.Warn(ctx, "unmarshall error")
	}
	wardDB, ierr := w.wardRepo.GetByKmsId(ctx, kmsId)
	if ierr != nil {
		return nil, ierr
	}
	go w.set(ctx, key, wardDB, expirationWardByKmsId)
	return wardDB, nil
}

func (w WardRepoDecorator) GetByKvId(ctx context.Context, kvId int64) (*domain.Ward, *common.Error) {
	key := w.genKeyCacheGetWardByKvId(kvId)
	var ward domain.Ward
	val, err := w.get(ctx, key).Result()
	w.handleRedisError(ctx, err)
	if err == nil {
		err = json.Unmarshal([]byte(val), &ward)
		if err == nil {
			return &ward, nil
		}
		log.Warn(ctx, "unmarshall error")
	}
	w.handleRedisError(ctx, err)
	wardDB, ierr := w.wardRepo.GetByKvId(ctx, kvId)
	if ierr != nil {
		return nil, ierr
	}
	go w.set(ctx, key, wardDB, expirationWardByKvId)
	return wardDB, nil
}
