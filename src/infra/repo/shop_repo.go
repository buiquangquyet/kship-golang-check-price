package repo

import (
	"check-price/src/common"
	"check-price/src/core/domain"
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func NewShopRepo(base *baseRepo) *ShopRepo {
	return &ShopRepo{
		base,
	}
}

type ShopRepo struct {
	*baseRepo
}

func (s *ShopRepo) GetByRetailerId(ctx context.Context, retailerId int64) (*domain.Shop, *common.Error) {
	shop := &domain.Shop{}
	cond := clause.Eq{Column: "retailer_id", Value: retailerId}
	if err := s.db.WithContext(ctx).Clauses(cond).Take(shop).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound(ctx, "shop", "not found").SetSource(common.SourceInfraService)
		}
		return nil, s.returnError(ctx, err)
	}
	return shop.DecryptPassword(), nil
}

func (s *ShopRepo) GetByCode(ctx context.Context, code string) (*domain.Shop, *common.Error) {
	shop := &domain.Shop{}
	cond := clause.Eq{Column: "code", Value: code}
	if err := s.db.WithContext(ctx).Clauses(cond).Take(shop).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound(ctx, "shop", "not found").SetSource(common.SourceInfraService)
		}
		return nil, s.returnError(ctx, err)
	}
	return shop.DecryptPassword(), nil
}
