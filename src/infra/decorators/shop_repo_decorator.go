package decorators

import (
	"check-price/src/common"
	"check-price/src/common/log"
	"check-price/src/core/domain"
	"check-price/src/infra/repo"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

const (
	keyCacheShopByCode   = "cache_shop_by_code"
	expirationShopByCode = 24 * time.Hour

	keyCacheShopByRetailerId   = "cache_shop_by_retailer_id"
	expirationShopByRetailerId = 24 * time.Hour
)

type ShopRepoDecorator struct {
	cache    redis.UniversalClient
	shopRepo *repo.ShopRepo
}

func NewShopRepoDecorator(shopRepo *repo.ShopRepo, cache redis.UniversalClient) domain.ShopRepo {
	return &ShopRepoDecorator{
		cache:    cache,
		shopRepo: shopRepo,
	}
}

func (s *ShopRepoDecorator) GetByRetailerId(ctx context.Context, retailerId int64) (*domain.Shop, *common.Error) {
	key := s.genKeyCacheGetByRetailerId(retailerId)
	var shop domain.Shop
	err := s.cache.Get(ctx, key).Scan(&shop)
	if err != nil {
		return &shop, nil
	}
	if err != redis.Nil {
		log.Error(ctx, "get redis error")
	}
	shopDB, ierr := s.shopRepo.GetByRetailerId(ctx, retailerId)
	if ierr != nil {
		return nil, ierr
	}
	go func() {
		s.cache.Set(ctx, key, shopDB, expirationShopByRetailerId)
	}()
	return shopDB, nil
}

func (s *ShopRepoDecorator) genKeyCacheGetByRetailerId(retailerId int64) string {
	return fmt.Sprintf("%s_%v", keyCacheShopByRetailerId, retailerId)
}

func (s *ShopRepoDecorator) GetByCode(ctx context.Context, code string) (*domain.Shop, *common.Error) {
	key := s.genKeyCacheGetByCode(code)
	var shop domain.Shop
	err := s.cache.Get(ctx, key).Scan(&shop)
	if err != nil {
		return &shop, nil
	}
	if err != redis.Nil {
		log.Error(ctx, "get redis error")
	}
	shopDB, ierr := s.shopRepo.GetByCode(ctx, code)
	if ierr != nil {
		return nil, ierr
	}
	go func() {
		s.cache.Set(ctx, key, shopDB, expirationShopByCode)
	}()
	return shopDB, nil
}

func (s *ShopRepoDecorator) genKeyCacheGetByCode(code string) string {
	return fmt.Sprintf("%s_%s", keyCacheShopByCode, code)
}
