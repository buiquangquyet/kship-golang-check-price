package repo

import (
	"check-price/src/common"
	"check-price/src/core/domain"
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func NewShopLevelRepo(base *baseRepo) *ShopLevelRepo {
	return &ShopLevelRepo{
		base,
	}
}

type ShopLevelRepo struct {
	*baseRepo
}

func (c ShopLevelRepo) GetById(ctx context.Context, id int64) (*domain.ShopLevel, *common.Error) {
	shopLevel := &domain.ShopLevel{}
	cond := clause.Eq{Column: "id", Value: id}
	if err := c.db.WithContext(ctx).Clauses(cond).Take(shopLevel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound(ctx, "shop level", "not found").SetSource(common.SourceInfraService)
		}
		return nil, c.returnError(ctx, err)
	}
	return shopLevel, nil
}
