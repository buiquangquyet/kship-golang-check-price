package decorators

import (
	"check-price/src/common"
	"check-price/src/core/domain"
	"check-price/src/infra/repo"
	"context"
	"github.com/go-redis/redis/v8"
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

func (s SettingShopRepoDecorator) GetByRetailerId(ctx context.Context, retailerId int64) ([]int64, *common.Error) {
	//TODO implement me
	panic("implement me")
}

func (s SettingShopRepoDecorator) GetEnableShopByRetailerId(ctx context.Context, retailerId int64) ([]*domain.SettingShop, *common.Error) {
	return nil, nil
}

func (s SettingShopRepoDecorator) GetServiceExtraEnableShop(ctx context.Context, retailerId int64) (bool, *common.Error) {
	//TODO implement me
	panic("implement me")
}
