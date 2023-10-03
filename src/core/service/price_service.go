package service

import (
	"check-price/src/common"
	"check-price/src/common/log"
	"check-price/src/core/constant"
	"check-price/src/core/domain"
	"check-price/src/core/dto"
	"check-price/src/core/strategy"
	"check-price/src/helpers"
	"check-price/src/present/httpui/request"
	"context"
)

type PriceService struct {
	shopRepo             domain.ShopRepo
	clientRepo           domain.ClientRepo
	settingShopRepo      domain.SettingShopRepo
	districtRepo         domain.DistrictRepo
	wardRepo             domain.WardRepo
	serviceRepo          domain.ServiceRepo
	shipStrategyResolver *strategy.ShipStrategyFilterResolver
}

func NewPriceService(shipStrategyResolver *strategy.ShipStrategyFilterResolver) *PriceService {
	return &PriceService{
		shipStrategyResolver: shipStrategyResolver,
	}
}

func (p *PriceService) GetPrice(ctx context.Context, req *request.GetPriceReRequest, tokenInfo *request.TokenInfo) ([]*domain.Price, *common.Error) {
	clientCode := req.ClientCode
	shop, err := p.shopRepo.GetByRetailerId(ctx, tokenInfo.RetailerId)
	if helpers.IsInternalError(err) {
		log.Error(ctx, err.Error())
		return nil, err
	}
	if shop == nil || !req.ActiveKShip {
		shop, err = p.shopRepo.GetByCode(ctx, constant.ShopDefaultTrial)
		if err != nil {
			log.Error(ctx, err.Error())
			return nil, err
		}
	}
	if ierr := p.validate(ctx, dto.NewValidatePrice(shop, req, tokenInfo)); ierr != nil {
		return nil, ierr
	}
	shipStrategy, exist := p.shipStrategyResolver.Resolve(clientCode)
	if !exist {
		log.Warn(ctx, "not support with partner:[%s]", clientCode)
		return nil, common.ErrBadRequest(ctx).SetDetail("partner not support").SetSource(common.SourceAPIService)
	}
	services := []string{"service_1", "service_2", "service_3"}
	prices, err := shipStrategy.GetMultiplePriceV3(ctx, "shop_code", services)
	if err != nil {
		log.IErr(ctx, err)
		return nil, err
	}
	return prices, nil
}
