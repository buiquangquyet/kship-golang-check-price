package decorators

import (
	"check-price/src/common"
	"check-price/src/core/domain"
	"check-price/src/infra/repo"
	"context"
	"time"
)

const (
	expirationCityById = 24 * time.Hour
)

type CityRepoDecorator struct {
	*baseDecorator
	cityRepo *repo.CityRepo
}

func NewCityRepoDecorator(base *baseDecorator, cityRepo *repo.CityRepo) domain.CityRepo {
	return &CityRepoDecorator{
		baseDecorator: base,
		cityRepo:      cityRepo,
	}
}

func (c CityRepoDecorator) GetById(ctx context.Context, id int64) (*domain.City, *common.Error) {
	key := c.genKeyCacheGetCityById(id)
	var city domain.City
	err := c.get(ctx, key).Scan(&city)
	if err == nil {
		return &city, nil
	}
	c.handleRedisError(ctx, err)
	cityDB, ierr := c.cityRepo.GetById(ctx, id)
	if ierr != nil {
		return nil, ierr
	}
	go c.set(ctx, key, cityDB, expirationCityById)
	return cityDB, nil
}
