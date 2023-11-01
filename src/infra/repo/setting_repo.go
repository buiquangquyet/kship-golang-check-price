package repo

import (
	"check-price/src/common"
	"check-price/src/core/domain"
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func NewSettingRepo(base *baseRepo) *SettingRepo {
	return &SettingRepo{
		base,
	}
}

type SettingRepo struct {
	*baseRepo
}

func (r SettingRepo) GetByName(ctx context.Context, name string) (*domain.Setting, *common.Error) {
	setting := &domain.Setting{}
	cond := clause.Eq{Column: "name", Value: name}
	if err := r.db.WithContext(ctx).Clauses(cond).Take(setting).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound(ctx, "setting", "not found").SetSource(common.SourceInfraService)
		}
		return nil, r.returnError(ctx, err)
	}
	return setting, nil
}
