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
	expirationConfigCodT0ByCodAndClientId = 1 * time.Hour
)

type ConfigCodT0RepoDecorator struct {
	*baseDecorator
	configCodT0Repo *repo.ConfigCodT0Repo
}

func NewConfigCodT0RepoDecorator(base *baseDecorator, configCodT0Repo *repo.ConfigCodT0Repo) domain.ConfigCodT0Repo {
	return &ConfigCodT0RepoDecorator{
		baseDecorator:   base,
		configCodT0Repo: configCodT0Repo,
	}
}

func (c ConfigCodT0RepoDecorator) GetByCodAndClientId(ctx context.Context, cod int64, clientId int64) ([]*domain.ConfigCodT0, *common.Error) {
	key := c.genKeyCacheGetConfigCodT0ByCodAndClientId(cod, clientId)
	var configCodT0s []*domain.ConfigCodT0
	val, err := c.get(ctx, key).Result()
	c.handleRedisError(ctx, err)
	if err == nil {
		err = json.Unmarshal([]byte(val), &configCodT0s)
		if err == nil {
			return configCodT0s, nil
		}
		log.Warn(ctx, "unmarshall error")
	}
	configCodT0sDB, ierr := c.configCodT0Repo.GetByCodAndClientId(ctx, cod, clientId)
	if ierr != nil {
		return nil, ierr
	}
	go c.set(ctx, key, configCodT0sDB, expirationConfigCodT0ByCodAndClientId)
	return configCodT0sDB, nil
}
