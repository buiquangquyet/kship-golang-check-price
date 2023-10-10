package repo

import (
	"check-price/src/common"
	"check-price/src/core/domain"
	"context"
)

func NewCityRepo(base *baseRepo) *CityRepo {
	return &CityRepo{
		base,
	}
}

type CityRepo struct {
	*baseRepo
}

func (c CityRepo) GetById(ctx context.Context, id int64) (*domain.City, *common.Error) {
	return nil, nil
}
