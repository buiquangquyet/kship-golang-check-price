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

func (p *priceRepo) GetResponse(ctx context.Context, clientCode string, senderWardId string, receiverWardId string) ([]*domain.Price, *common.Error) {
	//TODO implement me
	panic("implement me")
}

func (p *priceRepo) GetById(ctx context.Context, id int64) (*domain.Price, *common.Error) {
	//TODO implement me
	panic("implement me")
}
