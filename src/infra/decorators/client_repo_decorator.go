package decorators

import (
	"check-price/src/common"
	"check-price/src/core/domain"
	"check-price/src/infra/repo"
	"context"
	"github.com/go-redis/redis/v8"
)

type ClientRepoDecorator struct {
	*baseDecorator
	cache      redis.UniversalClient
	clientRepo *repo.ClientRepo
}

func NewClientRepoDecorator(base *baseDecorator, clientRepo *repo.ClientRepo, cache redis.UniversalClient) domain.ClientRepo {
	return &ClientRepoDecorator{
		baseDecorator: base,
		cache:         cache,
		clientRepo:    clientRepo,
	}
}

func (c ClientRepoDecorator) GetById(ctx context.Context, id int64) (*domain.Client, *common.Error) {
	//TODO implement me
	panic("implement me")
}

func (c ClientRepoDecorator) GetByCode(ctx context.Context, clientCode string) (*domain.Client, *common.Error) {
	//TODO implement me
	panic("implement me")
}
