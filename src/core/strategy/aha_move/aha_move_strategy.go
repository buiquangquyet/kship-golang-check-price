package aha_move

import (
	"check-price/src/common"
	"check-price/src/common/log"
	"check-price/src/core/constant"
	"check-price/src/core/domain"
	"check-price/src/core/strategy"
	"check-price/src/infra/external/ahamove"
	"check-price/src/present/httpui/request"
	"context"
)

type AhaMoveStrategy struct {
	wardRepo          domain.WardRepo
	districtRepo      domain.DistrictRepo
	cityRepo          domain.CityRepo
	clientRepo        domain.ClientRepo
	serviceRepo       domain.ServiceRepo
	ahaMoveExtService *ahamove.AhaMoveExtService
}

func NewAhaMoveStrategy(
	wardRepo domain.WardRepo,
	districtRepo domain.DistrictRepo,
	cityRepo domain.CityRepo,
	clientRepo domain.ClientRepo,
	serviceRepo domain.ServiceRepo,
	ahaMoveExtService *ahamove.AhaMoveExtService,
) strategy.ShipStrategy {
	return &AhaMoveStrategy{
		wardRepo:          wardRepo,
		districtRepo:      districtRepo,
		cityRepo:          cityRepo,
		clientRepo:        clientRepo,
		serviceRepo:       serviceRepo,
		ahaMoveExtService: ahaMoveExtService,
	}
}

func (g *AhaMoveStrategy) Code() string {
	return constant.AHAMOVEDeliveryCode
}

func (g *AhaMoveStrategy) Validate(_ context.Context, _ *request.GetPriceRequest) *common.Error {
	return nil
}

func (g *AhaMoveStrategy) GetMultiplePriceV3(ctx context.Context, shop *domain.Shop, req *request.GetPriceRequest, coupon string) (map[string]*domain.Price, *common.Error) {
	mapPrices := make(map[string]*domain.Price)
	prices, err := g.ahaMoveExtService.CheckPrice(ctx, shop)
	if err != nil {
		log.Error(ctx, err.Error())
		return nil, err
	}
	for _, price := range prices {
		mapPrices[price.Code] = price
	}
	return mapPrices, nil
}
