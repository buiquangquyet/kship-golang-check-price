package decorators

import (
	"check-price/src/common"
	"check-price/src/common/log"
	"check-price/src/core/domain"
	"check-price/src/core/enums"
	"check-price/src/infra/repo"
	"context"
	"encoding/json"
	"time"
)

const (
	expirationSettingShopByRetailerId = 12 * time.Hour
)

type SettingShopRepoDecorator struct {
	*baseDecorator
	settingShop *repo.SettingShopRepo
}

func (s SettingShopRepoDecorator) GetByRetailerIdAndModelId(ctx context.Context, modelType enums.ModelTypeSettingShop, retailerId int64, modelId int64) ([]*domain.SettingShop, *common.Error) {
	//TODO implement me
	panic("implement me")
}

func NewSettingShopRepoDecorator(base *baseDecorator, settingShop *repo.SettingShopRepo) domain.SettingShopRepo {
	return &SettingShopRepoDecorator{
		baseDecorator: base,
		settingShop:   settingShop,
	}
}

func (s SettingShopRepoDecorator) GetByRetailerId(ctx context.Context, modelType enums.ModelTypeSettingShop, retailerId int64) ([]*domain.SettingShop, *common.Error) {
	key := s.genKeyCacheGetSettingShopByRetailerId(modelType, retailerId)
	var settingShops []*domain.SettingShop
	val, err := s.get(ctx, key).Result()
	s.handleRedisError(ctx, err)
	if err == nil {
		err = json.Unmarshal([]byte(val), &settingShops)
		if err == nil {
			return settingShops, nil
		}
		log.Warn(ctx, "unmarshall error")
	}
	settingShopDB, ierr := s.settingShop.GetByRetailerId(ctx, modelType, retailerId)
	if ierr != nil {
		return nil, ierr
	}
	go s.set(ctx, key, settingShopDB, expirationSettingShopByRetailerId)
	return settingShopDB, nil
}
