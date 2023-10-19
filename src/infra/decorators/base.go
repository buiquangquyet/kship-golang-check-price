package decorators

import (
	"check-price/src/common"
	"check-price/src/common/log"
	"check-price/src/core/enums"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

const (
	keyCacheShopByCode       = "cache_shop_by_code"
	keyCacheShopByRetailerId = "cache_shop_by_retailer_id"

	keyCacheSettingShopByRetailerId = "cache_setting_shop"

	keyCacheClientByCode = "cache_client_by_code"

	keyCacheDistrictByKmsId = "cache_district_by_kms_id"
	keyCacheDistrictByKvId  = "cache_district_by_kv_id"
	keyCacheDistrictById    = "cache_district_by_id"

	keyCacheWardByKmsId = "cache_ward_by_kms_id"
	keyCacheWardByKvId  = "cache_ward_by_kv_id"

	keyCacheServiceByClientId   = "cache_service_by_client_id"
	keyCacheServiceByClientCode = "cache_service_by_client_code"

	keyCacheCityById = "cache_city_by_id"
)

type baseDecorator struct {
	cache redis.UniversalClient
}

func NewBaseDecorator(cache redis.UniversalClient) *baseDecorator {
	return &baseDecorator{
		cache: cache,
	}
}

func (b *baseDecorator) handleRedisError(ctx context.Context, err error) {
	if err == nil {
		return
	}
	if err != redis.Nil {
		log.Warn(ctx, "get redis error, err:[%s]", err.Error())
	}
}

func (b *baseDecorator) get(ctx context.Context, key string) *redis.StringCmd {
	return b.cache.Get(ctx, key)
}

func (b *baseDecorator) set(ctx context.Context, key string, value interface{}, exp time.Duration) {
	valueByte, err := json.Marshal(value)
	if err != nil {
		log.Warn(ctx, "marshal error")
	}
	err = b.cache.Set(common.Detach(ctx), key, valueByte, exp).Err()
	if err != nil {
		log.Warn(ctx, "set redis error, err:[%s]", err.Error())
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

func (b *baseDecorator) genKeyCacheGetDistrictByKmsId(kmsId int64) string {
	return fmt.Sprintf("%s_%v", keyCacheDistrictByKmsId, kmsId)
}

func (b *baseDecorator) genKeyCacheGetDistrictByKvId(kvId int64) string {
	return fmt.Sprintf("%s_%v", keyCacheDistrictByKvId, kvId)
}

func (b *baseDecorator) genKeyCacheGetDistrictById(id int64) string {
	return fmt.Sprintf("%s_%v", keyCacheDistrictById, id)
}

func (b *baseDecorator) genKeyCacheGetWardByKmsId(kmsId int64) string {
	return fmt.Sprintf("%s_%v", keyCacheWardByKmsId, kmsId)
}

func (b *baseDecorator) genKeyCacheGetWardByKvId(kvId int64) string {
	return fmt.Sprintf("%s_%v", keyCacheWardByKvId, kvId)
}

func (b *baseDecorator) genKeyCacheGetServiceByClientId(typeService enums.TypeService, status int, clientId int64) string {
	return fmt.Sprintf("%s_%v_%v_%v", keyCacheServiceByClientId, typeService.ToInt(), status, clientId)
}

func (b *baseDecorator) genKeyCacheGetServiceByClientCode(typeService enums.TypeService, status int, clientCode string) string {
	return fmt.Sprintf("%s_%v_%v_%s", keyCacheServiceByClientCode, typeService, status, clientCode)
}

func (b *baseDecorator) genKeyCacheGetCityById(id int64) string {
	return fmt.Sprintf("%s_%v", keyCacheCityById, id)
}
