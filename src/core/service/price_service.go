package service

import (
	"check-price/src/common"
	"check-price/src/common/log"
	"check-price/src/core/constant"
	"check-price/src/core/domain"
	"check-price/src/core/dto"
	"check-price/src/core/enums"
	"check-price/src/core/strategy"
	"check-price/src/helpers"
	"check-price/src/present/httpui/request"
	"context"
	"strings"
)

type PriceService struct {
	shopRepo             domain.ShopRepo
	serviceRepo          domain.ServiceRepo
	clientRepo           domain.ClientRepo
	validateService      *ValidateService
	settingShopRepo      domain.SettingShopRepo
	shipStrategyResolver *strategy.ShipStrategyFilterResolver
	codT0Service         *CodT0Service
}

func NewPriceService(
	shopRepo domain.ShopRepo,
	serviceRepo domain.ServiceRepo,
	clientRepo domain.ClientRepo,
	validateService *ValidateService,
	settingShopRepo domain.SettingShopRepo,
	shipStrategyResolver *strategy.ShipStrategyFilterResolver,
	codT0Service *CodT0Service,
) *PriceService {
	return &PriceService{
		shopRepo:             shopRepo,
		serviceRepo:          serviceRepo,
		clientRepo:           clientRepo,
		validateService:      validateService,
		settingShopRepo:      settingShopRepo,
		shipStrategyResolver: shipStrategyResolver,
		codT0Service:         codT0Service,
	}
}

func (p *PriceService) GetPrice(ctx context.Context, req *request.GetPriceRequest) ([]*domain.Price, *common.Error) {
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
	client, ierr := p.validateService.validatePrice(ctx, shop, req)
	if ierr != nil {
		return nil, ierr
	}
	shipStrategy, exist := p.shipStrategyResolver.Resolve(req.ClientCode)
	if !exist {
		log.Warn(ctx, "not support with partner:[%s]", req.ClientCode)
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

	prices, ierr := p.addInfo(ctx, dto.NewAddInfoDTO(shop, client, req), mapPrices)
	if ierr != nil {
		return nil, ierr
	}
	return prices, nil
}

func (p *PriceService) addInfo(ctx context.Context, addInfoDto *dto.AddInfoDto, mapPrices map[string]*domain.Price) ([]*domain.Price, *common.Error) {
	client, ierr := p.clientRepo.GetByCode(ctx, addInfoDto.Client.Code)
	if ierr != nil {
		log.Error(ctx, ierr.Error())
		return nil, ierr
	}
	servicesCode := make([]string, len(addInfoDto.Services))
	for i, service := range addInfoDto.Services {
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
		err := p.handlePriceSpecialService(ctx, addInfoDto, price)
		if err != nil {
			return nil, err
		}

		prices = append(prices, price)
	}
	return prices, nil
}

func (p *PriceService) handlePriceSpecialService(ctx context.Context, addInfoDto *dto.AddInfoDto, price *domain.Price) *common.Error {
	extraServiceCode := make([]string, len(addInfoDto.ExtraService))
	//payer := ""
	for i, service := range addInfoDto.ExtraService {
		if service.Code == "PAYMENT_BY" {
			//payer = service.Code
		}
		extraServiceCode[i] = service.Code
	}
	if helpers.InArray(extraServiceCode, constant.ServiceExtraCODST) && p.checkServiceExtraIsPossible(ctx, addInfoDto, constant.ServiceExtraCODST) {
		price.CalculatorCODST(addInfoDto.Shop, addInfoDto.Cod)
	}
	if helpers.InArray(extraServiceCode, constant.ServiceExtraCODT0) {
		err := p.codT0Service.addCodT0Price(ctx, price, addInfoDto)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *PriceService) checkServiceExtraIsPossible(ctx context.Context, addInfoDto *dto.AddInfoDto, extraServiceCode string) bool {
	extraService, err := p.serviceRepo.GetByCode(ctx, extraServiceCode)
	if helpers.IsInternalError(err) {
		log.Error(ctx, err.Error())
		return false
	}
	if err != nil {
		return false
	}
	if extraService.ClientsPossible == "" && strings.Contains(extraService.ClientsPossible, addInfoDto.Client.Code) {
		if extraService.OnBoardingStatus == constant.StatusEnableServiceExtra {
			return true
		}
		serviceExtraEnableShop, err := p.settingShopRepo.GetByRetailerId(ctx, enums.ModelTypeServiceExtraSettingShop, addInfoDto.RetailerId)
		if err != nil {
			log.Error(ctx, err.Error())
			return false
		}
		for _, service := range serviceExtraEnableShop {
			if extraService.Id == service.ModelId {
				return true
			}
		}
	}
	return true
}
