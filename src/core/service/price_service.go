package service

import (
	"check-price/src/common"
	"check-price/src/common/log"
	"check-price/src/core/constant"
	"check-price/src/core/domain"
	"check-price/src/core/strategy"
	"check-price/src/helpers"
	"check-price/src/present/httpui/request"
	"context"
)

type PriceService struct {
	shopRepo             domain.ShopRepo
	districtRepo         domain.DistrictRepo
	wardRepo             domain.WardRepo
	serviceRepo          domain.ServiceRepo
	settingShopRepo      domain.SettingShopRepo
	clientRepo           domain.ClientRepo
	shipStrategyResolver *strategy.ShipStrategyFilterResolver
}

func NewPriceService(
	shopRepo domain.ShopRepo,
	districtRepo domain.DistrictRepo,
	wardRepo domain.WardRepo,
	serviceRepo domain.ServiceRepo,
	settingShopRepo domain.SettingShopRepo,
	clientRepo domain.ClientRepo,
	shipStrategyResolver *strategy.ShipStrategyFilterResolver) *PriceService {
	return &PriceService{
		shopRepo:             shopRepo,
		districtRepo:         districtRepo,
		wardRepo:             wardRepo,
		serviceRepo:          serviceRepo,
		settingShopRepo:      settingShopRepo,
		clientRepo:           clientRepo,
		shipStrategyResolver: shipStrategyResolver,
	}
}

func (p *PriceService) GetPrice(ctx context.Context, req *request.GetPriceReRequest) ([]*domain.Price, *common.Error) {
	clientCode := req.ClientCode
	shop, err := p.shopRepo.GetByRetailerId(ctx, req.RetailerId)
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
	if ierr := p.validate(ctx, shop, req.RetailerId, req); ierr != nil {
		return nil, ierr
	}
	shipStrategy, exist := p.shipStrategyResolver.Resolve(clientCode)
	if !exist {
		log.Warn(ctx, "not support with partner:[%s]", clientCode)
		return nil, common.ErrBadRequest(ctx).SetDetail("partner not support").SetSource(common.SourceAPIService)
	}
	ierr := shipStrategy.Validate(ctx, req)
	if ierr != nil {
		return nil, ierr
	}
	prices, err := shipStrategy.GetMultiplePriceV3(ctx, shop, req)
	if err != nil {
		log.IErr(ctx, err)
		return nil, err
	}
	return prices, nil
}
