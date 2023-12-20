package ghtk

import (
	"check-price/src/common"
	"check-price/src/common/log"
	"check-price/src/core/constant"
	"check-price/src/core/domain"
	"check-price/src/core/param"
	"check-price/src/core/strategy"
	"check-price/src/infra/external/ghtk"
	"check-price/src/present/httpui/request"
	"context"
	"strconv"
	"sync"
)

type Strategy struct {
	ghtkExtService *ghtkext.Service
	baseStrategy   *strategy.BaseStrategy
}

func NewStrategy(
	ghtkExtService *ghtkext.Service,
	baseStrategy *strategy.BaseStrategy,
) strategy.ShipStrategy {
	return &Strategy{
		ghtkExtService: ghtkExtService,
		baseStrategy:   baseStrategy,
	}
}

func (s *Strategy) Code() string {
	return constant.GHTKDeliveryCode
}

func (s *Strategy) Validate(ctx context.Context, req *request.GetPriceRequest) *common.Error {
	ierr := common.ErrBadRequest(ctx)
	if req.ReceiverWardId == 0 {
		return ierr.SetMessage("Vui lòng nhập xã phường người nhận")
	}
	if req.ReceiverLocationId == 0 {
		return ierr.SetMessage("Vui lòng nhập quận huyện người nhận")
	}
	if req.MoneyCollection != 0 && req.MoneyCollection > 20000000 {
		return ierr.SetMessage("Chỉ nhận thu hộ COD tối đa 20,000,000").SetCode(4405)
	}
	for _, extraService := range req.ExtraService {
		if extraService.Code == constant.ServiceExtraGbh {
			value, err := strconv.ParseInt(extraService.Value, 10, 64)
			if err != nil {
				return common.ErrBadRequest(ctx).SetDetail("value extra service invalid")
			}
			if value > 20000000 {
				return common.ErrBadRequest(ctx).SetCode(4406)
			}
		}
	}
	return nil
}

func (s *Strategy) GetMultiplePriceV3(ctx context.Context, shop *domain.Shop, _ *domain.Client, req *request.GetPriceRequest, _ string) (map[string]*domain.Price, *common.Error) {
	address, err := s.baseStrategy.GetAddress(ctx, req)
	if err != nil {
		return nil, err
	}
	value, inspectFee, tags, err := s.getExtraService(ctx, req.ExtraService)
	if err != nil {
		return nil, err
	}
	isBBS, weight := s.getWeight(req)
	products := s.getProduct(req, isBBS)
	p := &param.GetPriceGHTKParam{
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
		NotDeliveredFee:  inspectFee,
	}

	var wg sync.WaitGroup
	mapPrices := make(map[string]*domain.Price)
	for _, service := range req.Services {
		wg.Add(1)
		go func(service *request.Service) {
			defer wg.Done()
			var price *domain.Price
			var err *common.Error
			if isBBS {
				price, err = s.ghtkExtService.GetPriceOver20(ctx, shop, service.Code, p)
			} else {
				price, err = s.ghtkExtService.GetPriceUnder20(ctx, shop, service.Code, p)
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

func (s *Strategy) getWeight(req *request.GetPriceRequest) (bool, int64) {
	isBBS := false
	product := req.Product
	weight := int64(product.ProductLength * product.ProductWidth * product.ProductHeight / 6)
	if weight < product.ProductWeight {
		weight = product.ProductWeight
	}
	if weight > 20000 {
		isBBS = true
	}
	return isBBS, weight
}

func (s *Strategy) getProduct(req *request.GetPriceRequest, isBBS bool) []*param.ProductGHTK {
	products := make([]*param.ProductGHTK, 0)
	if isBBS {
		products = append(products, &param.ProductGHTK{
			Name:     "kiện hàng",
			Quantity: 1,
			Weight:   req.ProductWeight / 1000,
			Width:    req.ProductWidth,
			Length:   req.ProductLength,
			Height:   req.ProductHeight,
		})
	}
	return products
}

func (s *Strategy) getExtraService(ctx context.Context, extraServices []*request.ExtraService) (int64, int64, []int, *common.Error) {
	var value, inspectFee int64
	var err error
	tags := make([]int, 0)
	for _, extraService := range extraServices {
		if extraService.Code == constant.ServiceExtraGbh {
			valueString, err := strconv.ParseInt(extraService.Value, 10, 64)
			if err != nil {
				return 0, 0, nil, common.ErrBadRequest(ctx).SetDetail("value extra service invalid")
			}
			value = valueString
		}
		if tag, exist := constant.MapGHTKTagByServiceExtra[extraService.Code]; exist {
			tags = append(tags, tag)
		}

		if tag, exist := constant.MapGHTKTagByShipperNote[extraService.Value]; exist {
			tags = append(tags, tag)
		}
		if extraService.Code == constant.ServiceExtraInspectFee {
			inspectFee, err = strconv.ParseInt(extraService.Value, 10, 64)
			if err != nil {
				return 0, 0, nil, common.ErrBadRequest(ctx).SetDetail("value extra service invalid")
			}
			if inspectFee < constant.MinFailDeliveryTaking {
				inspectFee = constant.MinFailDeliveryTaking
			}
			if inspectFee > constant.MaxFailDeliveryTaking {
				inspectFee = constant.MaxFailDeliveryTaking
			}
		}
	}
	return value, inspectFee, tags, nil
}
