package repo

import (
	"check-price/src/common"
	"check-price/src/core/domain"
	"check-price/src/core/enums"
	"check-price/src/helpers"
	"context"
	"gorm.io/gorm/clause"
)

func NewSettingShopRepo(base *baseRepo) *SettingShopRepo {
	return &SettingShopRepo{
		base,
	}
}

type SettingShopRepo struct {
	*baseRepo
}

func (s SettingShopRepo) GetByRetailerIdAndModelId(ctx context.Context, modelType enums.ModelTypeSettingShop, retailerId int64, modelId int64) ([]*domain.SettingShop, *common.Error) {
	settingShops := make([]*domain.SettingShop, 0)
	conds := []clause.Expression{
		clause.Eq{Column: "retailer_id", Value: retailerId},
		clause.Eq{Column: "model_type", Value: modelType.ToString()},
		clause.Eq{Column: "model_id", Value: modelId},
	}
	if err := s.db.WithContext(ctx).Clauses(conds...).Find(&settingShops).Error; err != nil {
		return nil, s.returnError(ctx, err)
	}
	return settingShops, nil
}

func (s SettingShopRepo) GetByRetailerId(ctx context.Context, modelType enums.ModelTypeSettingShop, retailerId int64) ([]*domain.SettingShop, *common.Error) {
	settingShops := make([]*domain.SettingShop, 0)
	conds := []clause.Expression{
		clause.Eq{Column: "retailer_id", Value: retailerId},
		clause.Eq{Column: "model_type", Value: modelType.ToString()},
	}
	if err := s.db.WithContext(ctx).Clauses(conds...).Find(&settingShops).Error; err != nil {
		return nil, s.returnError(ctx, err)
	}
	return settingShops, nil
}

func (s SettingShopRepo) GetByModelIdsAndValue(ctx context.Context, modelType enums.ModelTypeSettingShop, modelIds []int64, value string) ([]*domain.SettingShop, *common.Error) {
	settingShops := make([]*domain.SettingShop, 0)
	conds := []clause.Expression{
		clause.Eq{Column: "model_type", Value: modelType.ToString()},
		clause.IN{Column: "model_id", Values: helpers.ConvertTypesToInterfaces(modelIds)},
		clause.Eq{Column: "value", Value: value},
	}
	if err := s.db.WithContext(ctx).Clauses(conds...).Find(&settingShops).Error; err != nil {
		return nil, s.returnError(ctx, err)
	}
	return settingShops, nil
}
