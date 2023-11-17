package ahamove

import (
	"check-price/src/common"
	"check-price/src/common/log"
	"check-price/src/core/constant"
	"check-price/src/core/domain"
	"check-price/src/core/strategy"
	ahamoveext "check-price/src/infra/external/ahamove"
	"check-price/src/infra/external/aieliminating"
	"check-price/src/present/httpui/request"
	"context"
	"fmt"
)

type AhaMoveStrategy struct {
	baseStrategy            *strategy.BaseStrategy
	clientRepo              domain.ClientRepo
	serviceRepo             domain.ServiceRepo
	ahaMoveExtService       *ahamoveext.AhaMoveExtService
	aiEliminatingExtService *aieliminating.Service
}

func NewAhaMoveStrategy(
	baseStrategy *strategy.BaseStrategy,
	clientRepo domain.ClientRepo,
	serviceRepo domain.ServiceRepo,
	ahaMoveExtService *ahamoveext.AhaMoveExtService,
	aiEliminatingExtService *aieliminating.Service,
) strategy.ShipStrategy {
	return &AhaMoveStrategy{
		baseStrategy:            baseStrategy,
		clientRepo:              clientRepo,
		serviceRepo:             serviceRepo,
		ahaMoveExtService:       ahaMoveExtService,
		aiEliminatingExtService: aiEliminatingExtService,
	}
}

func (s *AhaMoveStrategy) Code() string {
	return constant.AHAMOVEDeliveryCode
}

func (s *AhaMoveStrategy) Validate(_ context.Context, _ *request.GetPriceRequest) *common.Error {
	return nil
}

func (s *AhaMoveStrategy) GetMultiplePriceV3(ctx context.Context, shop *domain.Shop, req *request.GetPriceRequest, coupon string) (map[string]*domain.Price, *common.Error) {
	senderAddress, receiverAddress, err := s.getAddressValue(ctx, req)
	if err != nil {
		return nil, err
	}
	fmt.Println(senderAddress, receiverAddress)
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

func (s *AhaMoveStrategy) getAddressValue(ctx context.Context, req *request.GetPriceRequest) (string, string, *common.Error) {
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

	return senderAddress, receiverAddress, nil
}
