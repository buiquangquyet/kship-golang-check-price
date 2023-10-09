package decorators

import (
	"check-price/src/common"
	"check-price/src/core/domain"
	"check-price/src/infra/repo"
	"context"
	"github.com/go-redis/redis/v8"
)

type CityRepoDecorator struct {
	*baseDecorator
	cache    redis.UniversalClient
	cityRepo *repo.CityRepo
}

func NewCityRepoDecorator(base *baseDecorator, cityRepo *repo.CityRepo, cache redis.UniversalClient) domain.CityRepo {
	return &CityRepoDecorator{
		baseDecorator: base,
		cache:         cache,
		cityRepo:      cityRepo,
	}
}
func (c CityRepoDecorator) GetById(ctx context.Context, id int64) (*domain.City, *common.Error) {
	//TODO implement me
	panic("implement me")
}
