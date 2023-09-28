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
	priceRepo            domain.PriceRepo
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
	//validate + cache
	shop, err := p.shopRepo.GetByRetailerId(ctx, tokenInfo.RetailerId)
	if helpers.IsInternalError(err) {
		log.Error(ctx, err.Error())
		return nil, err
	}
	//tai sao dk nhu nay moi vao cache check
	if shop == nil || !req.ActiveKShip {
		// response từ cache
		prices, err := p.priceRepo.GetResponse(ctx, clientCode, req.SenderWardId, req.ReceiverWardId)
		if err != nil {
			log.Error(ctx, err.Error())
			return nil, err
		}
		if prices != nil {
			return prices, nil
		}
		//	// get shop default, từ cache hay cả db phai hỏi lại
		shop, err = p.shopRepo.GetByCode(ctx, constant.ShopDefaultTrial)
		if err != nil {
			log.Error(ctx, err.Error())
			return nil, err
		}
	}
	if ierr := p.validate(ctx, shop, req, tokenInfo); ierr != nil {
		return nil, ierr
	}
	// vao thang check price node js
	shipStrategy, exist := p.shipStrategyResolver.Resolve(clientCode)
	if !exist {
		log.Warn(ctx, "not support with partner:[%s]", clientCode)
		return nil, common.ErrBadRequest(ctx).SetDetail("partner not support").SetSource(common.SourceAPIService)
	}
	//example multi services
	services := []string{"service_1", "service_2", "service_3"}
	prices, err := shipStrategy.GetMultiplePriceV3(ctx, "shop_code", services)
	if err != nil {
		log.IErr(ctx, err)
		return nil, err
	}
	return prices, nil
}

func (p *PriceService) validate(ctx context.Context, shop *domain.Shop, req *request.GetPriceReRequest, tokenInfo *request.TokenInfo) *common.Error {
	ierr := p.validateShop(ctx, shop, req.ClientCode)
	if ierr != nil {
		return ierr
	}
	ierr = p.validateClient(ctx, req.ClientCode, tokenInfo.RetailerId)
	if ierr != nil {
		return ierr
	}
	return nil
}

// done
func (p *PriceService) validateShop(ctx context.Context, shop *domain.Shop, clientCode string) *common.Error {
	ierr := common.ErrBadRequest(ctx)
	switch clientCode {
	case constant.VTPFWDeliveryCode:
		if shop.VtpUsername == "" || shop.VtpPassword == "" {
			return ierr.SetCode(3005)
		}
	case constant.VNPDeliveryCode:
		if shop.VnpCmsCode == "" {
			return ierr.SetCode(3008)
		}
	case constant.GHTKDeliveryCode:
		if shop.GHTKUsername == "" || shop.GHTKPassword == "" {
			return ierr.SetCode(3012)
		}
	case constant.JTFWDeliveryCode:
		if shop.JtCustomerId == "" {
			return ierr.SetCode(3005)
		}
	case constant.GHNFWDeliveryCode:
		if shop.GHNGWShopId == "" || shop.GHNFWPhone == "" {
			return ierr.SetCode(3005)
		}
	case constant.BESTFWDeliveryCode:
		if shop.UsernameBestFw == "" || shop.PasswordBestFw == "" {
			return ierr.SetCode(3005)
		}
	}
	return nil
}

// done
func (p *PriceService) validateClient(ctx context.Context, clientCode string, retailerId int64) *common.Error {
	client, err := p.clientRepo.GetByCode(ctx, clientCode)
	if helpers.IsInternalError(err) {
		log.IErr(ctx, err)
		return err
	}
	ierr := common.ErrBadRequest(ctx)
	if client == nil {
		log.Warn(ctx, "client is null")
		return ierr.SetCode(3001)
	}
	if client.Status == constant.DisableStatus {
		return ierr.SetCode(3002)
	}
	clientUnAllowedShop, err := p.settingShopRepo.GetByRetailerId(ctx, retailerId)
	if err != nil {
		log.IErr(ctx, err)
		return err
	}
	if helpers.InArray(clientUnAllowedShop, client.Id) {
		return ierr.SetCode(3004)
	}

	if client.OnBoardingStatus == constant.OnboardingDisable {
		clientSettingShop, err := p.settingShopRepo.GetEnableShopByRetailerId(ctx, retailerId)
		if err != nil {
			log.IErr(ctx, err)
			return err
		}
		if !helpers.InArray(clientSettingShop, client.Id) {
			return ierr.SetCode(3004)
		}
	}
	return nil
}

func (p *PriceService) validateLocation(ctx context.Context, req *request.GetPriceReRequest) *common.Error {
	ierr := common.ErrBadRequest(ctx)
	if req.VersionLocation == constant.VersionLocation2 {
		_, ierr = p.districtRepo.GetByKmsId(ctx, req.ReceiverLocationId)
	} else {
		_, ierr = p.districtRepo.GetByKvId(ctx, req.ReceiverLocationId)
	}
	if ierr != nil {
		return ierr.SetCode(4003)
	}
	if req.ClientCode == constant.GHTKDeliveryCode && req.SenderWardId == 0 {
		return ierr.SetCode(4004)
	}
	if req.SenderWardId != 0 {
		//Todo duplicate code
		if req.VersionLocation == constant.VersionLocation2 {
			_, ierr = p.wardRepo.GetByKmsId(ctx, req.SenderWardId)
		} else {
			_, ierr = p.wardRepo.GetByKvId(ctx, req.SenderWardId)
		}
		if ierr != nil {
			return ierr.SetCode(4004)
		}
	}
	ierr = common.ErrBadRequest(ctx)
	if !helpers.InArray(constant.ReceiverWardIdClientCode, req.ClientCode) && req.ReceiverWardId != 0 {
		return ierr.SetCode(4006)
	}
	//validate width, height, ...
	return nil
}

func (p *PriceService) validateService(ctx context.Context, req *request.GetPriceReRequest) *common.Error {
	ierr := common.ErrBadRequest(ctx)
	if req.Services == nil {
		return ierr.SetCode(4001)
	}
	serviceEnable, ierr := p.serviceRepo.GetServicesPluckCodeByClientCode(ctx, req.ClientCode)
	if ierr != nil {
		log.IErr(ctx, ierr)
		return ierr
	}
	for _, service := range req.Services {
		serviceCode := service.Code
		if helpers.InArray(serviceEnable, serviceCode) {
			return common.ErrBadRequest(ctx).SetCode(4009)
		}
	}
	return nil
}

func (p *PriceService) validateExtraService(ctx context.Context, req *request.GetPriceReRequest, tokenInfo *request.TokenInfo) *common.Error {
	ierr := common.ErrBadRequest(ctx)
	extraServiceRequestCodes := make([]string, len(req.ExtraService))
	for i, service := range req.ExtraService {
		extraServiceRequestCodes[i] = service.Code
	}
	if !helpers.InArray(extraServiceRequestCodes, constant.ServiceExtraCodePayment) {
		return ierr.SetCode(4013)
	}
	servicesExtraByClientCodes, ierr := p.serviceRepo.GetByServiceCode(ctx, req.ClientCode)
	if ierr != nil {
		log.IErr(ctx, ierr)
		return ierr
	}
	clientExtraServicesAllow := make([]string, 0)
	clientExtraServicesPaymentByAllow := make([]string, 0)
	for _, servicesExtraByClientCode := range servicesExtraByClientCodes {
		if servicesExtraByClientCode.Code == constant.ServiceExtraCodePayment {
			serviceExtraValue := constant.PaymentByFrom
			if servicesExtraByClientCode.Value != "" {
				serviceExtraValue = servicesExtraByClientCode.Value
			}
			clientExtraServicesPaymentByAllow = append(clientExtraServicesPaymentByAllow, serviceExtraValue)
		}
		if servicesExtraByClientCode.OnBoardingStatus {
			clientExtraServicesAllow = append(clientExtraServicesAllow, servicesExtraByClientCode.Code)
			continue
		}
		serviceExtraEnableShop, ierr := p.settingShopRepo.GetServiceExtraEnableShop(ctx, tokenInfo.RetailerId)
		if ierr != nil {
			log.IErr(ctx, ierr)
			return ierr
		}
		if serviceExtraEnableShop {
			clientExtraServicesAllow = append(clientExtraServicesAllow, servicesExtraByClientCode.Code)
		}
	}
	return nil
}
