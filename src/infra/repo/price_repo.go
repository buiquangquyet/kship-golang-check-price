package repo

import (
	"check-price/src/common"
	"check-price/src/core/domain"
	"context"
)

func NewPriceRepo(base *baseRepo) domain.PriceRepo {
	return &priceRepo{
		base,
	}
}

type priceRepo struct {
	*baseRepo
}

func (p *priceRepo) GetResponse(ctx context.Context) ([]*domain.Price, *common.Error) {
	var prices []*domain.Price
	key := ""
	if err := p.cache.Get(ctx, key).Scan(&prices); err != nil {
		return nil, p.returnError(ctx, err)
	}
	return prices, nil
}

func (p *priceRepo) GetById(ctx context.Context, id int64) (*domain.Price, *common.Error) {
	//TODO implement me
	panic("implement me")
}
