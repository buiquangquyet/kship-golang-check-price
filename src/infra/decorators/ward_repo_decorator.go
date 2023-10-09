package decorators

import (
	"check-price/src/common"
	"check-price/src/core/domain"
	"check-price/src/infra/repo"
	"context"
	"github.com/go-redis/redis/v8"
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

func (w WardRepoDecorator) GetByKmsId(ctx context.Context, senderWardId int64) (*domain.Ward, *common.Error) {
	//TODO implement me
	panic("implement me")
}

func (w WardRepoDecorator) GetByKvId(ctx context.Context, senderWardId int64) (*domain.Ward, *common.Error) {
	//TODO implement me
	panic("implement me")
}
