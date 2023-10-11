package repo

import (
	"check-price/src/common"
	"check-price/src/core/domain"
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func NewWardRepo(base *baseRepo) *WardRepo {
	return &WardRepo{
		base,
	}
}

type WardRepo struct {
	*baseRepo
}

func (w WardRepo) GetByKmsId(ctx context.Context, kmsId int64) (*domain.Ward, *common.Error) {
	ward := &domain.Ward{}
	conds := []clause.Expression{
		clause.Eq{Column: "kms_id", Value: kmsId},
		clause.Eq{Column: "retailer_id", Value: 1},
	}
	if err := w.db.WithContext(ctx).Clauses(conds...).Take(ward).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound(ctx, "ward", "not found").SetSource(common.SourceInfraService)
		}
		return nil, w.returnError(ctx, err)
	}
	return ward, nil
}

func (w WardRepo) GetByKvId(ctx context.Context, kvId int64) (*domain.Ward, *common.Error) {
	ward := &domain.Ward{}
	conds := []clause.Expression{
		clause.Eq{Column: "kv_id", Value: kvId},
		clause.Eq{Column: "retailer_id", Value: 1},
	}
	if err := w.db.WithContext(ctx).Clauses(conds...).Take(ward).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound(ctx, "ward", "not found").SetSource(common.SourceInfraService)
		}
		return nil, w.returnError(ctx, err)
	}
	return ward, nil
}
