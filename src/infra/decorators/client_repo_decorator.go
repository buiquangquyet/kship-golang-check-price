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
	expirationClientByCode = 12 * time.Hour
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
	key := c.genKeyCacheGetClientByCode(clientCode)
	var client domain.Client
	err := c.cache.Get(ctx, key).Scan(&client)
	if err != nil {
		return &client, nil
	}
	c.handleRedisError(ctx, err)
	clientDB, ierr := c.clientRepo.GetByCode(ctx, clientCode)
	if ierr != nil {
		return nil, ierr
	}
	go c.cache.Set(ctx, key, clientDB, expirationClientByCode)
	return clientDB, nil
}
