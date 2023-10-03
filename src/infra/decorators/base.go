package decorators

import (
	"check-price/src/common/log"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

const (
	keyCacheShopByCode       = "cache_shop_by_code"
	keyCacheShopByRetailerId = "cache_shop_by_retailer_id"
)

type baseDecorator struct {
}

func NewBaseDecorator() *baseDecorator {
	return &baseDecorator{}
}

func (b *baseDecorator) handleRedisError(ctx context.Context, err error) {
	if err != redis.Nil {
		log.Error(ctx, "get redis error")
	}
}

// shop domain
func (b *baseDecorator) genKeyCacheGetShopByRetailerId(retailerId int64) string {
	return fmt.Sprintf("%s_%v", keyCacheShopByRetailerId, retailerId)
}

func (b *baseDecorator) genKeyCacheGetShopByCode(code string) string {
	return fmt.Sprintf("%s_%s", keyCacheShopByCode, code)
}

//
