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
	expirationClientByCode = 12 * time.Hour
)

type ClientRepoDecorator struct {
	*baseDecorator
	clientRepo *repo.ClientRepo
}

func NewClientRepoDecorator(base *baseDecorator, clientRepo *repo.ClientRepo) domain.ClientRepo {
	return &ClientRepoDecorator{
		baseDecorator: base,
		clientRepo:    clientRepo,
	}
}

func (c ClientRepoDecorator) GetByCode(ctx context.Context, clientCode string) (*domain.Client, *common.Error) {
	key := c.genKeyCacheGetClientByCode(clientCode)
	var client domain.Client
	val, err := c.get(ctx, key).Result()
	c.handleRedisError(ctx, err)
	if err == nil {
		err = json.Unmarshal([]byte(val), &client)
		if err == nil {
			return &client, nil
		}
		log.Warn(ctx, "unmarshall error")
	}
	clientDB, ierr := c.clientRepo.GetByCode(ctx, clientCode)
	if ierr != nil {
		return nil, ierr
	}
	go c.set(ctx, key, clientDB, expirationClientByCode)
	return clientDB, nil
}
