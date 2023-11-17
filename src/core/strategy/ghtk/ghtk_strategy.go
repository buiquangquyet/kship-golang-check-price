package ghtk

import (
	"check-price/src/common"
	"check-price/src/common/log"
	"check-price/src/core/constant"
	"check-price/src/core/domain"
	"check-price/src/core/dto"
	"check-price/src/core/strategy"
	"check-price/src/infra/external/ghtk"
	"check-price/src/present/httpui/request"
	"context"
	"strconv"
	"sync"
)

type Strategy struct {
	wardRepo       domain.WardRepo
	districtRepo   domain.DistrictRepo
	cityRepo       domain.CityRepo
	clientRepo     domain.ClientRepo
	serviceRepo    domain.ServiceRepo
	ghtkExtService *ghtkext.Service
	baseStrategy   *strategy.BaseStrategy
}

func NewStrategy(
	wardRepo domain.WardRepo,
	districtRepo domain.DistrictRepo,
	cityRepo domain.CityRepo,
	clientRepo domain.ClientRepo,
	serviceRepo domain.ServiceRepo,
	ghtkExtService *ghtkext.Service,
) strategy.ShipStrategy {
	return &Strategy{
		wardRepo:       wardRepo,
		districtRepo:   districtRepo,
		cityRepo:       cityRepo,
		clientRepo:     clientRepo,
		serviceRepo:    serviceRepo,
		ghtkExtService: ghtkExtService,
	}
}

func (g *Strategy) Code() string {
	return constant.GHTKDeliveryCode
}

func (g *Strategy) Validate(ctx context.Context, req *request.GetPriceRequest) *common.Error {
	ierr := common.ErrBadRequest(ctx)
	if req.ReceiverWardId == 0 {
		return ierr.SetMessage("Vui lòng nhập xã phường người nhận")
	}
	if req.ReceiverLocationId == 0 {
		return ierr.SetMessage("Vui lòng nhập quận huyện người nhận")
	}
	if req.MoneyCollection != 0 && req.MoneyCollection > 20000000 {
		return ierr.SetMessage("Chỉ nhận thu hộ COD tối đa 20,000,000")
	}
	if req.ProductPrice != 0 && req.ProductPrice > 20000000 {
		return ierr.SetMessage("Chỉ nhận Khai giá tối đa 20,000,000")
	}
	return nil
}

func (g *Strategy) GetMultiplePriceV3(ctx context.Context, shop *domain.Shop, req *request.GetPriceRequest, _ string) (map[string]*domain.Price, *common.Error) {
	var wg sync.WaitGroup
	mapPrices := make(map[string]*domain.Price)
	isBBS := false
	product := req.Product
	weight := int64(product.ProductLength * product.ProductWidth * product.ProductHeight / 6)
	if weight < product.ProductWeight {
		weight = product.ProductWeight
	}
	if weight > 20000 {
		isBBS = true
	}

	getPriceParam, err := g.getPriceInput(ctx, isBBS, weight, req)
	if err != nil {
		return nil, err
	}
	for _, service := range req.Services {
		wg.Add(1)
		go func(service *request.Service) {
			defer wg.Done()
			var price *domain.Price
			var err *common.Error
			if isBBS {
				price, err = g.ghtkExtService.GetPriceOver20(ctx, shop, service.Code, getPriceParam)
			} else {
				price, err = g.ghtkExtService.GetPriceUnder20(ctx, shop, service.Code, getPriceParam)
			}
			if err != nil {
				log.Error(ctx, err.Error())
				return
			}
			mapPrices[service.Code] = price
		}(service)
	}
	wg.Wait()
	return mapPrices, nil
}

func (g *Strategy) getPriceInput(ctx context.Context, isBBS bool, weight int64, req *request.GetPriceRequest) (*dto.GetPriceInputDto, *common.Error) {
	address, err := g.baseStrategy.GetAddress(ctx, req)
	if err != nil {
		return nil, err
	}
	products := make([]*dto.Product, 0)
	if isBBS {
		products = append(products, &dto.Product{
			Name:     "kiện hàng",
			Quantity: 1,
			Weight:   req.ProductWeight / 1000,
			Width:    req.ProductWidth,
			Length:   req.ProductLength,
			Height:   req.ProductHeight,
		})
	}
	var value int64 = 0
	for _, extraService := range req.ExtraService {
		if extraService.Code == "GBH" {
			valueString, err := strconv.ParseInt(extraService.Value, 10, 64)
			if err != nil {
				return nil, common.ErrBadRequest(ctx).SetDetail("value extra service invalid")
			}
			value = valueString
		}
	}
	tags := make([]int, 0)
	for _, service := range req.ExtraService {
		if tag, exist := constant.MapGHTKExtraService[service.Code]; exist {
			tags = append(tags, tag)
		}
	}
	return &dto.GetPriceInputDto{
		PickProvince:     address.PickProvince.Name,
		PickDistrict:     address.PickDistrict.Name,
		PickWard:         address.PickWard.Name,
		ReceiverProvince: address.ReceiverProvince.Name,
		ReceiverDistrict: address.ReceiverDistrict.Name,
		ReceiverWard:     address.ReceiverWard.Name,
		Address:          req.ReceiverAddress,
		Products:         products,
		Weight:           weight,
		Value:            value,
		Tags:             tags,
	}, nil
}
