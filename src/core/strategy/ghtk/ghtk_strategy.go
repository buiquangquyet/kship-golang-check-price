package ghtk

import (
	"check-price/src/common"
	"check-price/src/common/log"
	"check-price/src/core/constant"
	"check-price/src/core/domain"
	"check-price/src/core/strategy"
	"check-price/src/infra/external/ghtk"
	"check-price/src/present/httpui/request"
	"context"
	"sync"
)

type GHTKStrategy struct {
	ghtkExtService *ghtk.GHTKExtService
}

func NewGHTKStrategy(ghtkExtService *ghtk.GHTKExtService) strategy.ShipStrategy {
	return &GHTKStrategy{
		ghtkExtService: ghtkExtService,
	}
}

func (g *GHTKStrategy) Code() string {
	return constant.GHTKDeliveryCode
}

func (g *GHTKStrategy) Validate(ctx context.Context, req *request.GetPriceReRequest) *common.Error {
	ierr := common.ErrBadRequest(ctx)
	if req.SenderWardId == 0 {
		return ierr.SetMessage("Vui lòng nhập xã phường người gửi")
	}
	if req.ReceiverWardId == 0 {
		return ierr.SetMessage("Vui lòng nhập xã phường người nhận")
	}
	if req.SenderLocationId == 0 {
		return ierr.SetMessage("Vui lòng nhập quận huyện người gửi")
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

// tam thoi de services []string
func (g *GHTKStrategy) GetMultiplePriceV3(ctx context.Context, shop *domain.Shop, req *request.GetPriceReRequest) ([]*domain.Price, *common.Error) {
	var wg sync.WaitGroup
	mapPrices := make(map[string]*domain.Price)
	isBBS := false
	product := req.Product
	weight := product.ProductLength * product.ProductWidth * product.ProductHeight / 6
	if weight < product.ProductWeight {
		weight = product.ProductWeight
	}
	if weight > 20000 {
		isBBS = true
	}
	for _, service := range req.Services {
		wg.Add(1)
		go func(service *request.Service) {
			defer wg.Done()
			var price *domain.Price
			var err *common.Error
			if isBBS {
				price, err = g.ghtkExtService.GetPriceOver20(ctx, shop, service.Code)
			} else {
				price, err = g.ghtkExtService.GetPriceUnder20(ctx, shop, service.Code)
			}
			if err != nil {
				log.Error(ctx, err.Error())
				return
			}
			mapPrices[service.Code] = price
		}(service)
	}
	wg.Wait()
	prices := make([]*domain.Price, 0)
	for service, price := range mapPrices {
		price.Code = service
		prices = append(prices, price)
	}
	return prices, nil
}
