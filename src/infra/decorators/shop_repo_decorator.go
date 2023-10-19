package decorators

import (
	"check-price/src/common"
	"check-price/src/common/log"
	"check-price/src/core/domain"
	"check-price/src/infra/repo"
	"context"
	"encoding/json"
	"time"
)

const (
	expirationShopByCode       = 24 * time.Hour
	expirationShopByRetailerId = 24 * time.Hour
)

type ShopRepoDecorator struct {
	*baseDecorator
	shopRepo *repo.ShopRepo
}

func NewShopRepoDecorator(base *baseDecorator, shopRepo *repo.ShopRepo) domain.ShopRepo {
	return &ShopRepoDecorator{
		baseDecorator: base,
		shopRepo:      shopRepo,
	}
}

func (s *ShopRepoDecorator) GetByRetailerId(ctx context.Context, retailerId int64) (*domain.Shop, *common.Error) {
	key := s.genKeyCacheGetShopByRetailerId(retailerId)
	var shop domain.Shop
	val, err := s.get(ctx, key).Result()
	s.handleRedisError(ctx, err)
	if err == nil {
		err = json.Unmarshal([]byte(val), &shop)
		if err == nil {
			return &shop, nil
		}
		log.Warn(ctx, "unmarshall error")
	}
	shopDB, ierr := s.shopRepo.GetByRetailerId(ctx, retailerId)
	if ierr != nil {
		return nil, ierr
	}
	go s.set(ctx, key, shopDB, expirationShopByRetailerId)
	return shopDB, nil
}

func (s *ShopRepoDecorator) GetByCode(ctx context.Context, code string) (*domain.Shop, *common.Error) {
	key := s.genKeyCacheGetShopByCode(code)
	var shop domain.Shop
	val, err := s.get(ctx, key).Result()
	s.handleRedisError(ctx, err)
	if err == nil {
		err = json.Unmarshal([]byte(val), &shop)
		if err == nil {
			return &shop, nil
		}
		log.Warn(ctx, "unmarshall error")
	}
	shopDB, ierr := s.shopRepo.GetByCode(ctx, code)
	if ierr != nil {
		return nil, ierr
	}
	go s.set(ctx, key, shopDB, expirationShopByCode)
	return shopDB, nil
}
