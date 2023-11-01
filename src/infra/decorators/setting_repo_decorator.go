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
	expirationSettingByName = 30 * 24 * time.Hour
)

type SettingRepoDecorator struct {
	*baseDecorator
	settingRepo *repo.SettingRepo
}

func NewSettingRepoDecorator(base *baseDecorator, settingRepo *repo.SettingRepo) domain.SettingRepo {
	return &SettingRepoDecorator{
		baseDecorator: base,
		settingRepo:   settingRepo,
	}
}

func (c SettingRepoDecorator) GetByName(ctx context.Context, name string) (*domain.Setting, *common.Error) {
	key := c.keyCacheSettingByName(name)
	var setting domain.Setting
	val, err := c.get(ctx, key).Result()
	c.handleRedisError(ctx, err)
	if err == nil {
		err = json.Unmarshal([]byte(val), &setting)
		if err == nil {
			return &setting, nil
		}
		log.Warn(ctx, "unmarshall error")
	}
	settingDB, ierr := c.settingRepo.GetByName(ctx, name)
	if ierr != nil {
		return nil, ierr
	}
	go c.set(ctx, key, settingDB, expirationSettingByName)
	return settingDB, nil
}
