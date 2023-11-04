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
	expirationShopLevelById = 24 * time.Hour
)

type ShopLevelRepoDecorator struct {
	*baseDecorator
	shopLevelRepo *repo.ShopLevelRepo
}

func NewShopLevelRepoDecorator(base *baseDecorator, shopLevelRepo *repo.ShopLevelRepo) domain.ShopLevelRepo {
	return &ShopLevelRepoDecorator{
		baseDecorator: base,
		shopLevelRepo: shopLevelRepo,
	}
}

func (s *ShopLevelRepoDecorator) GetById(ctx context.Context, id int64) (*domain.ShopLevel, *common.Error) {
	key := s.genKeyCacheShopLevelById(id)
	var shopLevel domain.ShopLevel
	val, err := s.get(ctx, key).Result()
	s.handleRedisError(ctx, err)
	if err == nil {
		err = json.Unmarshal([]byte(val), &shopLevel)
		if err == nil {
			return &shopLevel, nil
		}
		log.Warn(ctx, "unmarshall error")
	}
	shopLevelDB, ierr := s.shopLevelRepo.GetById(ctx, id)
	if ierr != nil {
		return nil, ierr
	}
	go s.set(ctx, key, shopLevelDB, expirationShopLevelById)
	return shopLevelDB, nil
}
