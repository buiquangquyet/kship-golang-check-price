package repo

import (
	"check-price/src/common"
	"check-price/src/core/domain"
	"context"
)

func NewShopRepo(base *baseRepo) domain.ShopRepo {
	return &ShopRepo{
		base,
	}
}

type ShopRepo struct {
	*baseRepo
}

func (s *ShopRepo) GetByRetailerId(ctx context.Context, retailerId int64) (*domain.Shop, *common.Error) {
	//TODO implement me
	panic("implement me")
}

func (s *ShopRepo) GetByCode(ctx context.Context, shopCode string) (*domain.Shop, *common.Error) {
	//TODO implement me
	panic("implement me")
}
