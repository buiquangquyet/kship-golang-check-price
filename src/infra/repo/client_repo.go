package repo

import (
	"check-price/src/common"
	"check-price/src/core/domain"
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func NewClientRepo(base *baseRepo) *ClientRepo {
	return &ClientRepo{
		base,
	}
}

type ClientRepo struct {
	*baseRepo
}

func (c ClientRepo) GetByCode(ctx context.Context, clientCode string) (*domain.Client, *common.Error) {
	client := &domain.Client{}
	cond := clause.Eq{Column: "code", Value: clientCode}
	if err := c.db.WithContext(ctx).Clauses(cond).Take(client).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound(ctx, "client", "not found").SetSource(common.SourceInfraService)
		}
		return nil, c.returnError(ctx, err)
	}
	return client, nil
}
