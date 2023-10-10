package decorators

import (
	"check-price/src/common"
	"check-price/src/core/domain"
	"check-price/src/core/enums"
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

func (s SettingShopRepoDecorator) GetByRetailerId(ctx context.Context, modelType enums.ModelTypeSettingShop, retailerId int64) ([]*domain.SettingShop, *common.Error) {
	//TODO implement me
	panic("implement me")
}
