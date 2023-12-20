package ahamove

import (
	"check-price/src/common"
	"check-price/src/common/log"
	"check-price/src/core/constant"
	"check-price/src/core/domain"
	"check-price/src/core/enums"
	"check-price/src/core/param"
	"check-price/src/core/strategy"
	ahamoveext "check-price/src/infra/external/ahamove"
	"check-price/src/infra/external/aieliminating"
	"check-price/src/present/httpui/request"
	"context"
	"fmt"
	"strconv"
	"time"
)

type Strategy struct {
	settingShop             domain.SettingShopRepo
	serviceRepo             domain.ServiceRepo
	baseStrategy            *strategy.BaseStrategy
	ahaMoveExtService       *ahamoveext.Service
	aiEliminatingExtService *aieliminating.Service
}

func NewStrategy(
	settingShop domain.SettingShopRepo,
	serviceRepo domain.ServiceRepo,
	baseStrategy *strategy.BaseStrategy,
	ahaMoveExtService *ahamoveext.Service,
	aiEliminatingExtService *aieliminating.Service,
) strategy.ShipStrategy {
	return &Strategy{
		settingShop:             settingShop,
		serviceRepo:             serviceRepo,
		baseStrategy:            baseStrategy,
		ahaMoveExtService:       ahaMoveExtService,
		aiEliminatingExtService: aiEliminatingExtService,
	}
}

func (s *Strategy) Code() string {
	return constant.AHAMOVEDeliveryCode
}

func (s *Strategy) Validate(_ context.Context, _ *request.GetPriceRequest) *common.Error {
	return nil
}

func (s *Strategy) GetMultiplePriceV3(ctx context.Context, shop *domain.Shop, client *domain.Client, req *request.GetPriceRequest, coupon string) (map[string]*domain.Price, *common.Error) {
	if shop.Code == constant.ShopDefaultTrial {
		return nil, common.ErrBadRequest(ctx).SetCode(2002)
	}
	provinceId, senderAddress, receiverAddress, err := s.getAddressValue(ctx, req)
	if err != nil {
		return nil, err
	}
	paymentMethod, orderTime, err := s.getExtraService(ctx, req.ExtraService)
	if err != nil {
		return nil, err
	}
	services, err := s.getServices(ctx, req.Services, req.ExtraService, client, provinceId)
	if err != nil {
		return nil, err
	}
	getPriceParam := &param.GetPriceAhaMoveParam{
		Path: [2]*param.Path{
			{Address: senderAddress},
			{Address: receiverAddress, Cod: req.MoneyCollection},
		},
		PaymentMethod: paymentMethod,
		PromoCode:     coupon,
		OrderTime:     orderTime,
		Services:      services,
	}

	mapPrices := make(map[string]*domain.Price)
	prices, err := s.ahaMoveExtService.CheckPrice(ctx, shop, getPriceParam)
	if err != nil {
		log.Error(ctx, err.Error())
		return nil, err
	}
	for _, price := range prices {
		mapPrices[price.Code] = price
	}
	return mapPrices, nil
}

func (s *Strategy) getAddressValue(ctx context.Context, req *request.GetPriceRequest) (int64, string, string, *common.Error) {
	var senderAddress string
	var receiverAddress string
	address, err := s.baseStrategy.GetAddress(ctx, req)
	if err != nil {
		return 0, "", "", err
	}
	if req.SenderAddress != "" {
		senderAddress, err = s.aiEliminatingExtService.Redundancy(ctx, req.SenderAddress, address.PickWard.Name, address.PickDistrict.Name, address.PickProvince.Name)
		if err != nil {
			log.Error(ctx, err.Error())
			return 0, "", "", err
		}
	}
	if req.ReceiverAddress != "" {
		receiverAddress, err = s.aiEliminatingExtService.Redundancy(ctx, req.ReceiverAddress, address.ReceiverWard.Name, address.ReceiverDistrict.Name, address.ReceiverProvince.Name)
		if err != nil {
			log.Error(ctx, err.Error())
			return 0, "", "", err
		}
	}
	senderAddress = fmt.Sprintf("%s, %s, %s, %s", senderAddress, address.PickWard.Name, address.PickDistrict.Name, address.PickProvince.Name)
	receiverAddress = fmt.Sprintf("%s, %s, %s, %s", receiverAddress, address.ReceiverWard.Name, address.ReceiverDistrict.Name, address.ReceiverProvince.Name)
	if address.PickProvince.Name != address.ReceiverProvince.Name {
		log.Error(ctx, "province invalid")
		return 0, "", "", common.ErrBadRequest(ctx).SetCode(1002)
	}
	return address.PickProvince.Id, senderAddress, receiverAddress, nil
}

func (s *Strategy) getExtraService(ctx context.Context, extraServices []*request.ExtraService) (string, int64, *common.Error) {
	paymentMethod := "CASH"
	var orderTime int64 = 0
	for _, extraService := range extraServices {
		if extraService.Code == constant.ServiceExtraPrepaid {
			paymentMethod = "BALANCE"
		}
		if extraService.Code == constant.ServiceExtraScheduled {
			dateString := "02-01-2006 15:04"
			orderTimeParse, err := time.Parse(dateString, extraService.Value)
			if err != nil {
				log.Error(ctx, "scheduled invalid, scheduled:[%s]", extraService.Value)
				return "", 0, common.ErrBadRequest(ctx).SetCode(1111)
			}
			orderTime = orderTimeParse.Unix()
		}
	}
	return paymentMethod, orderTime, nil
}

func (s *Strategy) getServices(ctx context.Context, services []*request.Service, extraServices []*request.ExtraService, client *domain.Client, provinceId int64) ([]*param.ServiceAhaMove, *common.Error) {
	codes := make([]string, len(services))
	for i, service := range services {
		codes[i] = service.Code
	}
	servicesDB, err := s.serviceRepo.GetByClientIdAndCodes(ctx, enums.TypeServiceDV, codes, client.Id)
	if err != nil {
		log.Error(ctx, err.Error())
		return nil, err
	}
	servicesMap := make(map[int64]*domain.Service)
	serviceIds := make([]int64, len(servicesDB))
	for i, service := range servicesDB {
		servicesMap[service.Id] = service
		serviceIds[i] = service.Id
	}
	settingShops, err := s.settingShop.GetByModelIdsAndValue(ctx, enums.ModelTypeCitiesPossible, serviceIds, strconv.FormatInt(provinceId, 10))
	if err != nil {
		log.Error(ctx, err.Error())
		return nil, err
	}
	servicesParam := make([]string, 0)
	if len(settingShops) == 0 {
		servicesParam = codes
	} else {
		for _, settingShop := range settingShops {
			servicesParam = append(servicesParam, servicesMap[settingShop.ModelId].Code)
		}
	}
	servicesAhaMoveParam := make([]*param.ServiceAhaMove, len(servicesParam))
	requestParam := make([]*param.Request, 0)
	for i, service := range servicesParam {
		for _, extraService := range extraServices {
			switch extraService.Code {
			case constant.ServiceExtraInspectFee, constant.ServiceExtraRoundTrip, "BOCXEP-2", "BOCXEP-3", constant.ServiceExtraThermalBag:
				num, err := strconv.Atoi(extraService.Value)
				if err != nil {
					log.Error(ctx, "value extra service invalid, scheduled:[%s]", extraService.Value)
					return nil, common.ErrBadRequest(ctx).SetCode(1111)
				}
				requestParam = append(requestParam, &param.Request{
					Id:  fmt.Sprintf("%s-%s", service, extraService.Code),
					Num: num,
				})
			case "BULKY":
				requestParam = append(requestParam, &param.Request{
					Id:       fmt.Sprintf("%s-%s", service, extraService.Code),
					TierCode: extraService.Value,
				})
			default:
				continue
			}
		}
		servicesAhaMoveParam[i] = &param.ServiceAhaMove{
			Id:       service,
			Requests: requestParam,
		}
	}
	return servicesAhaMoveParam, nil
}
