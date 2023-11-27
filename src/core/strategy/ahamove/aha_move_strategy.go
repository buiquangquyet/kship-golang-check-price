package ahamove

import (
	"check-price/src/common"
	"check-price/src/common/log"
	"check-price/src/core/constant"
	"check-price/src/core/domain"
	"check-price/src/core/param"
	"check-price/src/core/strategy"
	ahamoveext "check-price/src/infra/external/ahamove"
	"check-price/src/infra/external/aieliminating"
	"check-price/src/present/httpui/request"
	"context"
	"fmt"
	"time"
)

type Strategy struct {
	baseStrategy            *strategy.BaseStrategy
	ahaMoveExtService       *ahamoveext.Service
	aiEliminatingExtService *aieliminating.Service
}

func NewStrategy(
	baseStrategy *strategy.BaseStrategy,
	ahaMoveExtService *ahamoveext.Service,
	aiEliminatingExtService *aieliminating.Service,
) strategy.ShipStrategy {
	return &Strategy{
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

func (s *Strategy) GetMultiplePriceV3(ctx context.Context, shop *domain.Shop, req *request.GetPriceRequest, coupon string) (map[string]*domain.Price, *common.Error) {
	if shop.Code == constant.ShopDefaultTrial {
		return nil, common.ErrBadRequest(ctx).SetCode(2002)
	}
	senderAddress, receiverAddress, err := s.getAddressValue(ctx, req)
	if err != nil {
		return nil, err
	}

	paymentMethod := "CASH"
	var orderTime int64 = 0
	for _, extraService := range req.ExtraService {
		if extraService.Code == constant.ServiceExtraPrepaid {
			paymentMethod = "BALANCE"
		}
		if extraService.Code == constant.ServiceExtraScheduled {
			dateString := "02-01-2006 15:04"
			orderTimeParse, err := time.Parse(dateString, extraService.Value)
			if err != nil {
				//log
				return nil, common.ErrBadRequest(ctx).SetCode(1111)
			}
			orderTime = orderTimeParse.Unix()
		}
	}

	_ = &param.GetPriceAhaMoveParam{
		Path: [2]*param.Path{
			{Address: senderAddress},
			{Address: receiverAddress, Cod: req.MoneyCollection},
		},
		PaymentMethod: paymentMethod,
		PromoCode:     coupon,
		OrderTime:     orderTime,
		Services:      nil,
	}

	mapPrices := make(map[string]*domain.Price)
	prices, err := s.ahaMoveExtService.CheckPrice(ctx, shop)
	if err != nil {
		log.Error(ctx, err.Error())
		return nil, err
	}
	for _, price := range prices {
		mapPrices[price.Code] = price
	}
	return mapPrices, nil
}

func (s *Strategy) getAddressValue(ctx context.Context, req *request.GetPriceRequest) (string, string, *common.Error) {
	var senderAddress string
	var receiverAddress string
	address, err := s.baseStrategy.GetAddress(ctx, req)
	if err != nil {
		return "", "", err
	}
	if req.SenderAddress != "" {
		senderAddress, err = s.aiEliminatingExtService.Redundancy(ctx, req.SenderAddress, address.PickWard.Name, address.PickDistrict.Name, address.PickProvince.Name)
		if err != nil {
			log.Error(ctx, err.Error())
			return "", "", err
		}
	}
	if req.ReceiverAddress != "" {
		receiverAddress, err = s.aiEliminatingExtService.Redundancy(ctx, req.ReceiverAddress, address.ReceiverWard.Name, address.ReceiverDistrict.Name, address.ReceiverProvince.Name)
		if err != nil {
			log.Error(ctx, err.Error())
			return "", "", err
		}
	}
	senderAddress = fmt.Sprintf("%s, %s, %s, %s", senderAddress, address.PickWard.Name, address.PickDistrict.Name, address.PickProvince.Name)
	receiverAddress = fmt.Sprintf("%s, %s, %s, %s", receiverAddress, address.ReceiverWard.Name, address.ReceiverDistrict.Name, address.ReceiverProvince.Name)
	if address.PickProvince.Name != address.ReceiverProvince.Name {
		log.Error(ctx, "")
		//return
	}
	return senderAddress, receiverAddress, nil
}

func (s *Strategy) getServices(ctx context.Context, services []*request.Service) ([]param.ServiceAhaMove, *common.Error) {
	return nil, nil
}
