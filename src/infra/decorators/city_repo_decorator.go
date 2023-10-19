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
	val, err := c.get(ctx, key).Result()
	c.handleRedisError(ctx, err)
	if err == nil {
		err = json.Unmarshal([]byte(val), &city)
		if err == nil {
			return &city, nil
		}
		log.Warn(ctx, "unmarshall error")
	}
	cityDB, ierr := c.cityRepo.GetById(ctx, id)
	if ierr != nil {
		return nil, ierr
	}
	go c.set(ctx, key, cityDB, expirationCityById)
	return cityDB, nil
}
