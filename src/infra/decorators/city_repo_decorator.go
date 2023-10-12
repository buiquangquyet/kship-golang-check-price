package decorators

import (
	"check-price/src/common"
	"check-price/src/core/domain"
	"check-price/src/infra/repo"
	"context"
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
	//TODO implement me
	panic("implement me")
}
