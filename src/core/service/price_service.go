package service

import (
	"check-price/src/common"
	"check-price/src/common/log"
	"check-price/src/core/constant"
	"check-price/src/core/domain"
	"check-price/src/core/enums"
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
	ierr := p.validate(ctx, shop, req.RetailerId, req)
	if ierr != nil {
		return nil, ierr
	}
	shipStrategy, exist := p.shipStrategyResolver.Resolve(clientCode)
	if !exist {
		log.Warn(ctx, "not support with partner:[%s]", clientCode)
		return nil, common.ErrBadRequest(ctx).SetDetail("partner not support").SetSource(common.SourceAPIService)
	}
	ierr = shipStrategy.Validate(ctx, req)
	if ierr != nil {
		return nil, ierr
	}
	mapPrices, err := shipStrategy.GetMultiplePriceV3(ctx, shop, req)
	if err != nil {
		log.IErr(ctx, err)
		return nil, err
	}
	prices, ierr := p.addInfo(ctx, req.ClientCode, req.Services, mapPrices)
	if ierr != nil {
		return nil, ierr
	}
	return prices, nil
}

func (p *PriceService) addInfo(ctx context.Context, clientCode string, servicesReq []*request.Service, mapPrices map[string]*domain.Price) ([]*domain.Price, *common.Error) {
	client, ierr := p.clientRepo.GetByCode(ctx, clientCode)
	if ierr != nil {
		log.Error(ctx, ierr.Error())
		return nil, ierr
	}
	servicesCode := make([]string, len(servicesReq))
	for i, service := range servicesReq {
		servicesCode[i] = service.Code
	}
	services, ierr := p.serviceRepo.GetByClientIdAndCodes(ctx, enums.TypeServiceDV, servicesCode, client.Id)
	if ierr != nil {
		log.Error(ctx, ierr.Error())
		return nil, ierr
	}
	mapServices := make(map[string]*domain.Service)
	for _, service := range services {
		mapServices[service.Code] = service
	}
	prices := make([]*domain.Price, 0)
	for serviceCode, price := range mapPrices {
		price.Code = serviceCode
		price.SetClientInfo(client)
		price.SetServiceInfo(mapServices[serviceCode])

		prices = append(prices, price)
	}
	return prices, nil
}

func (p *PriceService) handlePriceSpecialService(ctx context.Context, price *domain.Price, shop *domain.Shop, extraService []*request.ExtraService, cod int64) *common.Error {
	extraServiceCode := make([]string, len(extraService))
	payer := ""
	for i, service := range extraService {
		if service.Code == "PAYMENT_BY" {
			payer = service.Code
		}
		extraServiceCode[i] = service.Code
	}
	if helpers.InArray(extraServiceCode, constant.ServiceExtraCODST) && p.checkServiceExtraIsPossible(ctx) {
		price.CalculatorCODST(shop, cod)
	}
	return nil
}

func (p *PriceService) checkServiceExtraIsPossible(ctx context.Context) bool {

	//Todo code
	return true
}
