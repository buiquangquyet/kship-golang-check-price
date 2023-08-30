package service

import (
	"check-price/src/common"
	"check-price/src/common/log"
	"check-price/src/core/constant"
	"check-price/src/core/domain"
	"check-price/src/helpers"
	"check-price/src/present/httpui/request"
	"context"
)

type PriceService struct {
	shopRepo  domain.ShopRepo
	priceRepo domain.PriceRepo
}

func NewPriceService(shopRepo domain.ShopRepo) *PriceService {
	return &PriceService{
		shopRepo: shopRepo,
	}
}

func (p *PriceService) GetPrice(ctx context.Context, req *request.GetPriceReRequest) ([]*domain.Price, *common.Error) {
	shop, err := p.shopRepo.GetByRetailerId(ctx, "")
	if helpers.IsInternalError(err) {
		log.Error(ctx, err.Error())
		return nil, err
	}
	//tai sao dk nhu nay moi vao cache check
	if shop == nil || req.ActiveKShip {
		// response từ cache
		prices, err := p.priceRepo.GetResponse(ctx)
		if err != nil {
			log.Error(ctx, err.Error())
			return nil, err
		}
		if prices != nil {
			return prices, nil
		}
		// get shop default, từ cache hay cả db phai hỏi lại
		shop, err := p.shopRepo.GetByCode(ctx, constant.ShopDefaultTrial)
		if err != nil {
			log.Error(ctx, err.Error())
			return nil, err
		}

	}

	return nil, nil
}

func (p *PriceService) GetShopDefault(ctx context.Context, req *request.GetPriceReRequest) ([]*domain.Price, *common.Error) {
	return nil, nil
}
