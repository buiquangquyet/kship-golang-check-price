package strategy

import (
	"check-price/src/core/domain"
	"check-price/src/infra/external"
	"context"
	"github.com/opentracing/opentracing-go/log"
	"sync"
)

type GhtkStrategy struct {
	ghtkExtService *external.GHTKExtService
}

func (g *GhtkStrategy) Code() string {
	return "ghtk"
}

// tam thoi de services []string
func (g *GhtkStrategy) GetMultiplePriceV3(ctx context.Context, shopCode string, services []string) {
	_ = g.ghtkExtService.Connect(ctx, shopCode)
	//call ghtk
	var wg sync.WaitGroup
	prices := make(map[string]*domain.Price)
	for _, service := range services {
		wg.Add(1)
		go func(service string) {
			price, err := g.ghtkExtService.GetPriceFromDelivery(ctx, service)
			if err != nil {
				log.Error(err)
			}
			prices[service] = price
			defer wg.Done()
		}(service)
	}
	wg.Wait()

}
