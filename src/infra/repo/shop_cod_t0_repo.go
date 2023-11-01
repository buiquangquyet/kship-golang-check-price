package repo

import (
	"check-price/src/common"
	"check-price/src/core/domain"
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func NewShopCodT0Repo(base *baseRepo) *ShopCodT0Repo {
	return &ShopCodT0Repo{
		base,
	}
}

type ShopCodT0Repo struct {
	*baseRepo
}

func (r ShopCodT0Repo) GetByShopId(ctx context.Context, shopId int64) (*domain.ShopCodT0, *common.Error) {
	shopCodT0 := &domain.ShopCodT0{}
	cond := clause.Eq{Column: "shop_id", Value: shopId}
	if err := r.db.WithContext(ctx).Clauses(cond).Take(shopCodT0).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound(ctx, "shop cod t0", "not found").SetSource(common.SourceInfraService)
		}
		return nil, r.returnError(ctx, err)
	}
	return shopCodT0, nil
}
