package decorators

import (
	"check-price/src/common"
	"check-price/src/core/domain"
	"check-price/src/infra/repo"
	"context"
	"github.com/go-redis/redis/v8"
)

type ServiceRepoDecorator struct {
	*baseDecorator
	cache       redis.UniversalClient
	serviceRepo *repo.ServiceRepo
}

func NewServiceRepoDecorator(base *baseDecorator, serviceRepo *repo.ServiceRepo, cache redis.UniversalClient) domain.ServiceRepo {
	return &ServiceRepoDecorator{
		baseDecorator: base,
		cache:         cache,
		serviceRepo:   serviceRepo,
	}
}

func (s ServiceRepoDecorator) GetServicesPluckCodeByClientCode(ctx context.Context, clientCode string) ([]string, *common.Error) {
	//TODO implement me
	panic("implement me")
}

func (s ServiceRepoDecorator) GetByServiceCode(ctx context.Context, clientCode string) ([]*domain.Service, *common.Error) {
	//TODO implement me
	panic("implement me")
}
