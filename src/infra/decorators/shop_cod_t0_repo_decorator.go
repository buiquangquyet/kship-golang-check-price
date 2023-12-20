package decorators

import (
	"check-price/src/common"
	"check-price/src/common/log"
	"check-price/src/core/domain"
	"check-price/src/infra/repo"
	"context"
	"encoding/json"
	"time"
)

const (
	expirationShopCodT0ByShopId = 24 * time.Hour
)

type ShopCodT0RepoDecorator struct {
	*baseDecorator
	shopCodT0Repo *repo.ShopCodT0Repo
}

func NewShopCodT0RepoDecorator(base *baseDecorator, shopCodT0Repo *repo.ShopCodT0Repo) domain.ShopCodT0Repo {
	return &ShopCodT0RepoDecorator{
		baseDecorator: base,
		shopCodT0Repo: shopCodT0Repo,
	}
}

func (c ShopCodT0RepoDecorator) GetByShopId(ctx context.Context, shopId int64) (*domain.ShopCodT0, *common.Error) {
	key := c.genKeyCacheShopCodT0ByShopId(shopId)
	var shopCodT0 domain.ShopCodT0
	val, err := c.get(ctx, key).Result()
	c.handleRedisError(ctx, err)
	if err == nil {
		err = json.Unmarshal([]byte(val), &shopCodT0)
		if err == nil {
			return &shopCodT0, nil
		}
		log.Warn(ctx, "unmarshall error")
	}
	shopCodT0DB, ierr := c.shopCodT0Repo.GetByShopId(ctx, shopId)
	if ierr != nil {
		return nil, ierr
	}
	go c.set(ctx, key, shopCodT0DB, expirationShopCodT0ByShopId)
	return shopCodT0DB, nil
}
