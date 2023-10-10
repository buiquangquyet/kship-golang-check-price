package decorators

import (
	"check-price/src/common"
	"check-price/src/core/domain"
	"check-price/src/core/enums"
	"check-price/src/infra/repo"
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

const (
	expirationSettingShopByRetailerId = 12 * time.Hour
)

type SettingShopRepoDecorator struct {
	*baseDecorator
	cache       redis.UniversalClient
	settingShop *repo.SettingShopRepo
}

func NewSettingShopRepoDecorator(base *baseDecorator, settingShop *repo.SettingShopRepo, cache redis.UniversalClient) domain.SettingShopRepo {
	return &SettingShopRepoDecorator{
		baseDecorator: base,
		cache:         cache,
		settingShop:   settingShop,
	}
}

func (s SettingShopRepoDecorator) GetByRetailerId(ctx context.Context, modelType enums.ModelTypeSettingShop, retailerId int64) ([]*domain.SettingShop, *common.Error) {
	key := s.genKeyCacheGetSettingShopByRetailerId(modelType, retailerId)
	var settingShops []*domain.SettingShop
	err := s.cache.Get(ctx, key).Scan(&settingShops)
	if err != nil {
		return settingShops, nil
	}
	s.handleRedisError(ctx, err)
	settingShopDB, ierr := s.settingShop.GetByRetailerId(ctx, modelType, retailerId)
	if ierr != nil {
		return nil, ierr
	}
	go s.cache.Set(ctx, key, settingShopDB, expirationSettingShopByRetailerId)
	return settingShopDB, nil
}
