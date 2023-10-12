package repo

import (
	"check-price/src/common"
	"check-price/src/core/domain"
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func NewDistrictRepo(base *baseRepo) *DistrictRepo {
	return &DistrictRepo{
		base,
	}
}

type DistrictRepo struct {
	*baseRepo
}

func (d DistrictRepo) GetByKmsId(ctx context.Context, kmsId int64) (*domain.District, *common.Error) {
	district := &domain.District{}
	conds := []clause.Expression{
		clause.Eq{Column: "kms_id", Value: kmsId},
		clause.Eq{Column: "mapping_status", Value: 1},
	}
	if err := d.db.WithContext(ctx).Clauses(conds...).Take(district).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound(ctx, "district", "not found").SetSource(common.SourceInfraService)
		}
		return nil, d.returnError(ctx, err)
	}
	return district, nil
}

func (d DistrictRepo) GetByKvId(ctx context.Context, kvId int64) (*domain.District, *common.Error) {
	district := &domain.District{}
	conds := []clause.Expression{
		clause.Eq{Column: "kv_id", Value: kvId},
		clause.Eq{Column: "mapping_status", Value: 1},
	}
	if err := d.db.WithContext(ctx).Clauses(conds...).Take(district).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound(ctx, "district", "not found").SetSource(common.SourceInfraService)
		}
		return nil, d.returnError(ctx, err)
	}
	return district, nil
}

func (d DistrictRepo) GetById(ctx context.Context, id int64) (*domain.District, *common.Error) {
	district := &domain.District{}
	conds := []clause.Expression{
		clause.Eq{Column: "id", Value: id},
		clause.Eq{Column: "mapping_status", Value: 1},
	}
	if err := d.db.WithContext(ctx).Clauses(conds...).Take(district).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound(ctx, "district", "not found").SetSource(common.SourceInfraService)
		}
		return nil, d.returnError(ctx, err)
	}
	return district, nil
}
