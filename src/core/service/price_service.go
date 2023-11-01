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
	serviceRepo          domain.ServiceRepo
	clientRepo           domain.ClientRepo
	validateService      *ValidateService
	shipStrategyResolver *strategy.ShipStrategyFilterResolver
	codT0Service         *CodT0Service
}

func NewPriceService(
	shopRepo domain.ShopRepo,
	serviceRepo domain.ServiceRepo,
	clientRepo domain.ClientRepo,
	validateService *ValidateService,
	shipStrategyResolver *strategy.ShipStrategyFilterResolver,
	codT0Service *CodT0Service,
) *PriceService {
	return &PriceService{
		shopRepo:             shopRepo,
		serviceRepo:          serviceRepo,
		clientRepo:           clientRepo,
		validateService:      validateService,
		shipStrategyResolver: shipStrategyResolver,
		codT0Service:         codT0Service,
	}
}

func (p *PriceService) GetPrice(ctx context.Context, req *request.GetPriceRequest) ([]*domain.Price, *common.Error) {
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
	ierr := p.validateService.validatePrice(ctx, shop, req.RetailerId, req)
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
	prices, ierr := p.addInfo(ctx, shop, req.ClientCode, req, mapPrices)
	if ierr != nil {
		return nil, ierr
	}
	return prices, nil
}

func (p *PriceService) addInfo(ctx context.Context, shop *domain.Shop, clientCode string, req *request.GetPriceRequest, mapPrices map[string]*domain.Price) ([]*domain.Price, *common.Error) {
	client, ierr := p.clientRepo.GetByCode(ctx, clientCode)
	if ierr != nil {
		log.Error(ctx, ierr.Error())
		return nil, ierr
	}
	servicesCode := make([]string, len(req.Services))
	for i, service := range req.Services {
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
		err := p.handlePriceSpecialService(ctx, price, shop, req.ExtraService, req.MoneyCollection)
		if err != nil {
			return nil, err
		}
		err = p.codT0Service.addCodT0Price(ctx, price, req, client, shop)
		if err != nil {
			return nil, err
		}
		prices = append(prices, price)
	}
	return prices, nil
}

func (p *PriceService) handlePriceSpecialService(ctx context.Context, price *domain.Price, shop *domain.Shop, extraService []*request.ExtraService, cod int64) *common.Error {
	extraServiceCode := make([]string, len(extraService))
	//payer := ""
	for i, service := range extraService {
		if service.Code == "PAYMENT_BY" {
			//payer = service.Code
		}
		extraServiceCode[i] = service.Code
	}
	if helpers.InArray(extraServiceCode, constant.ServiceExtraCODST) && p.checkServiceExtraIsPossible(ctx) {
		price.CalculatorCODST(shop, cod)
	}
	if helpers.InArray(extraServiceCode, constant.ServiceExtraCODT0) {

	}
	return nil
}

func (p *PriceService) checkServiceExtraIsPossible(_ context.Context) bool {

	//Todo code
	return true
}
