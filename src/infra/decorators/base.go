package decorators

import (
	"check-price/src/common/log"
	"check-price/src/core/enums"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

const (
	keyCacheShopByCode       = "cache_shop_by_code"
	keyCacheShopByRetailerId = "cache_shop_by_retailer_id"

	keyCacheSettingShopByRetailerId = "cache_setting_shop"

	keyCacheClientByCode = "cache_client_by_code"
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

func (b *baseDecorator) genKeyCacheGetSettingShopByRetailerId(modelType enums.ModelTypeSettingShop, retailerId int64) string {
	return fmt.Sprintf("%s_%s_%v", keyCacheSettingShopByRetailerId, modelType.ToString(), retailerId)
}

func (b *baseDecorator) genKeyCacheGetClientByCode(code string) string {
	return fmt.Sprintf("%s_%s", keyCacheClientByCode, code)
}
