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
)

type PriceService struct {
	shopRepo             domain.ShopRepo
	serviceRepo          domain.ServiceRepo
	clientRepo           domain.ClientRepo
	validateService      *ValidateService
	voucherService       *VoucherService
	settingShopRepo      domain.SettingShopRepo
	shipStrategyResolver *strategy.ShipStrategyFilterResolver
	extraService         *ExtraService
}

func NewPriceService(
	shopRepo domain.ShopRepo,
	serviceRepo domain.ServiceRepo,
	clientRepo domain.ClientRepo,
	validateService *ValidateService,
	voucherService *VoucherService,
	settingShopRepo domain.SettingShopRepo,
	shipStrategyResolver *strategy.ShipStrategyFilterResolver,
	extraService *ExtraService,
) *PriceService {
	return &PriceService{
		shopRepo:             shopRepo,
		serviceRepo:          serviceRepo,
		clientRepo:           clientRepo,
		validateService:      validateService,
		voucherService:       voucherService,
		settingShopRepo:      settingShopRepo,
		shipStrategyResolver: shipStrategyResolver,
		extraService:         extraService,
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
		req.RetailerId = shop.RetailerId
	}
	client, ierr := p.validateService.validatePrice(ctx, shop, req)
	if ierr != nil {
		return nil, ierr
	}
	addInfoDTO := dto.NewAddInfoDTO(shop, client, req)
	callTo, voucher, err := p.voucherService.checkVoucher(ctx, addInfoDTO)
	if err != nil {
		return nil, err
	}
	coupon := ""
	switch callTo {
	case enums.TypeVoucherUseKv:
		//handle trong SetCouponInfo
		//voucher !=0
	case enums.TypeVoucherUseDelivery:
		//ban sang cac hang
		coupon = addInfoDTO.Coupon
		voucher = 0
	case enums.TypeVoucherNotUse:
		//khong lam gi
	}
	addInfoDTO.CallTo = callTo
	addInfoDTO.Voucher = voucher

	shipStrategy, exist := p.shipStrategyResolver.Resolve(req.ClientCode)
	if !exist {
		log.Warn(ctx, "not support with partner:[%s]", req.ClientCode)
		return nil, common.ErrBadRequest(ctx).SetDetail("partner not support").SetSource(common.SourceAPIService)
	}
	ierr = shipStrategy.Validate(ctx, req)
	if ierr != nil {
		return nil, ierr
	}
	mapPrices, err := shipStrategy.GetMultiplePriceV3(ctx, shop, client, req, coupon)
	if err != nil {
		return nil, err
	}

	prices, ierr := p.addInfo(ctx, addInfoDTO, mapPrices)
	if ierr != nil {
		return nil, ierr
	}
	return prices, nil
}

func (p *PriceService) addInfo(ctx context.Context, addInfoDto *dto.AddInfoDto, mapPrices map[string]*domain.Price) ([]*domain.Price, *common.Error) {
	servicesCode := make([]string, len(addInfoDto.Services))
	for i, service := range addInfoDto.Services {
		servicesCode[i] = service.Code
	}
	services, ierr := p.serviceRepo.GetByClientIdAndCodes(ctx, enums.TypeServiceDV, servicesCode, addInfoDto.Client.Id)
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
		price.SetFee()
		if addInfoDto.Voucher == 0 && addInfoDto.CallTo == enums.TypeVoucherNotUse {
			price.SetOldTotalPrice(0)
		}
		price.SetCode(serviceCode)
		price.SetClientInfo(addInfoDto.Client)
		price.SetServiceInfo(mapServices[serviceCode])
		err := p.extraService.handlePriceSpecialService(ctx, price, addInfoDto)
		if err != nil {
			return nil, err
		}
		price.SetTotalPrice(addInfoDto.Payer)
		if addInfoDto.CallTo == enums.TypeVoucherUseKv {
			price.SetCouponInfo(addInfoDto.Voucher)
		} else {
			price.SetOtherFee()
		}
		if addInfoDto.Client.Code == constant.GHTKDeliveryCode && addInfoDto.Payer == constant.PaymentByTo {
			price.HandleFeeGHTK()
		}
		prices = append(prices, price)
	}

	return prices, nil
}
