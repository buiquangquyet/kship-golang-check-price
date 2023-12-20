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

	keyCacheSettingShopByRetailerId           = "cache_setting_shop_by_retailer_id"
	keyCacheSettingShopByRetailerIdAndModelId = "cache_setting_shop_by_retailer_id_and_model_id"
	keyCacheSettingShopByModelsAndValue       = "cache_setting_shop_by_models_and_value"

	keyCacheClientByCode = "cache_client_by_code"

	keyCacheDistrictByKmsId = "cache_district_by_kms_id"
	keyCacheDistrictByKvId  = "cache_district_by_kv_id"
	keyCacheDistrictById    = "cache_district_by_id"

	keyCacheWardByKmsId = "cache_ward_by_kms_id"
	keyCacheWardByKvId  = "cache_ward_by_kv_id"

	keyCacheServiceByClientId        = "cache_service_by_client_id"
	keyCacheServiceByClientCode      = "cache_service_by_client_code"
	keyCacheServiceByClientIdAndCode = "cache_service_by_code_and_client_code"
	keyCacheServiceByCode            = "cache_service_by_code"

	keyCacheCityById = "cache_city_by_id"

	keyCacheConfigCodT0ByCodAndClientId = "cache_config_cod_t0_by_cod_and_client_id"

	keyCacheSettingByName = "cache_setting_by_name"

	keyCacheShopCodT0ByShopId = "cache_shop_cod_t0_by_shop_id"

	keyCacheShopLevelById = "cache_shop_level_by_id"
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

func (b *baseDecorator) genKeyCacheGetShopByRetailerId(retailerId int64) string {
	return fmt.Sprintf("%s_%v", keyCacheShopByRetailerId, retailerId)
}

func (b *baseDecorator) genKeyCacheGetShopByCode(code string) string {
	return fmt.Sprintf("%s_%s", keyCacheShopByCode, code)
}

func (b *baseDecorator) genKeyCacheGetSettingShopByRetailerId(modelType enums.ModelTypeSettingShop, retailerId int64) string {
	return fmt.Sprintf("%s_%s_%v", keyCacheSettingShopByRetailerId, modelType.ToString(), retailerId)
}

func (b *baseDecorator) genKeyCacheGetSettingShopByRetailerIdAndModelId(modelType enums.ModelTypeSettingShop, retailerId int64, modelId int64) string {
	return fmt.Sprintf("%s_%s_%v_%v", keyCacheSettingShopByRetailerIdAndModelId, modelType.ToString(), retailerId, modelId)
}

func (b *baseDecorator) genKeyCacheGetSettingShopByModelsAndValue(modelType enums.ModelTypeSettingShop, value string) string {
	return fmt.Sprintf("%s_%s_%s", keyCacheSettingShopByModelsAndValue, modelType.ToString(), value)
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

func (b *baseDecorator) genKeyCacheGetServiceByClientIdAndCodes(typeService enums.TypeService, codes []string, clientId int64) string {
	return fmt.Sprintf("%s_%v_%v_%v", keyCacheServiceByClientIdAndCode, typeService, codes, clientId)
}

func (b *baseDecorator) genKeyCacheGetServiceByCode(code string) string {
	return fmt.Sprintf("%s_%s", keyCacheServiceByCode, code)
}

func (b *baseDecorator) genKeyCacheGetCityById(id int64) string {
	return fmt.Sprintf("%s_%v", keyCacheCityById, id)
}

func (b *baseDecorator) genKeyCacheGetConfigCodT0ByCodAndClientId(cod int64, clientId int64) string {
	return fmt.Sprintf("%s_%v_%v", keyCacheConfigCodT0ByCodAndClientId, cod, clientId)
}

func (b *baseDecorator) genKeyCacheSettingByName(name string) string {
	return fmt.Sprintf("%s_%s", keyCacheSettingByName, name)
}

func (b *baseDecorator) genKeyCacheShopCodT0ByShopId(shopId int64) string {
	return fmt.Sprintf("%s_%v", keyCacheShopCodT0ByShopId, shopId)
}

func (b *baseDecorator) genKeyCacheShopLevelById(id int64) string {
	return fmt.Sprintf("%s_%v", keyCacheShopLevelById, id)
}
