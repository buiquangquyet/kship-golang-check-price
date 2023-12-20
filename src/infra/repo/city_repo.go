package repo

import (
	"check-price/src/common"
	"check-price/src/core/domain"
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	city := &domain.City{}
	cond := clause.Eq{Column: "id", Value: id}
	if err := c.db.WithContext(ctx).Clauses(cond).Take(city).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound(ctx, "city", "not found").SetSource(common.SourceInfraService)
		}
		return nil, c.returnError(ctx, err)
	}
	return city, nil
}
