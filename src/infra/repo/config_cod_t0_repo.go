package repo

import (
	"check-price/src/common"
	"check-price/src/core/domain"
	"context"
	"gorm.io/gorm/clause"
)

type ConfigCodT0Repo struct {
	*baseRepo
}

func NewConfigCodT0Repo(base *baseRepo) *ConfigCodT0Repo {
	return &ConfigCodT0Repo{
		base,
	}
}

func (c ConfigCodT0Repo) GetByCodAndClientId(ctx context.Context, cod int64, clientId int64) ([]*domain.ConfigCodT0, *common.Error) {
	configCodT0 := make([]*domain.ConfigCodT0, 0)
	conds := []clause.Expression{
		clause.Gte{Column: "cod_to", Value: cod},
		clause.Lte{Column: "cod_from", Value: cod},
		clause.Eq{Column: "client_id", Value: clientId},
	}
	if err := c.db.WithContext(ctx).Clauses(conds...).Find(&configCodT0).Error; err != nil {
		return nil, c.returnError(ctx, err)
	}
	return configCodT0, nil
}
