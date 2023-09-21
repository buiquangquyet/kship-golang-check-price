package strategy

import (
	"check-price/src/common"
	"check-price/src/common/log"
	"check-price/src/core/constant"
	"check-price/src/core/domain"
	"check-price/src/infra/external"
	"context"
	"sync"
)

type GHTKStrategy struct {
	ghtkExtService *external.GHTKExtService
}

func NewGHTKStrategy(ghtkExtService *external.GHTKExtService) ShipStrategy {
	return &GHTKStrategy{
		ghtkExtService: ghtkExtService,
	}
}

func (g *GHTKStrategy) Code() string {
	return constant.GHTKDeliveryCode
}

// tam thoi de services []string
func (g *GHTKStrategy) GetMultiplePriceV3(ctx context.Context, shopCode string, services []string) ([]*domain.Price, *common.Error) {
	_ = g.ghtkExtService.Connect(ctx, shopCode)
	//call ghtk
	var wg sync.WaitGroup
	mapPrices := make(map[string]*domain.Price)
	for _, service := range services {
		wg.Add(1)
		go func(service string) {
			defer wg.Done()
			price, err := g.ghtkExtService.GetPriceFromDelivery(ctx, service)
			if err != nil {
				log.Error(ctx, err.Error())
				return
			}
			mapPrices[service] = price
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
