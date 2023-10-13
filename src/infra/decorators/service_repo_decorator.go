package decorators

import (
	"check-price/src/common"
	"check-price/src/core/domain"
	"check-price/src/core/enums"
	"check-price/src/infra/repo"
	"context"
	"time"
)

const (
	expirationServiceByClientId   = 12 * time.Hour
	expirationServiceByClientCode = 12 * time.Hour
)

type ServiceRepoDecorator struct {
	*baseDecorator
	serviceRepo *repo.ServiceRepo
}

func NewServiceRepoDecorator(base *baseDecorator, serviceRepo *repo.ServiceRepo) domain.ServiceRepo {
	return &ServiceRepoDecorator{
		baseDecorator: base,
		serviceRepo:   serviceRepo,
	}
}

func (s ServiceRepoDecorator) GetByClientId(ctx context.Context, typeService enums.TypeService, status int, clientId int64) ([]*domain.Service, *common.Error) {
	key := s.genKeyCacheGetServiceByClientId(typeService, status, clientId)
	var services []*domain.Service
	err := s.get(ctx, key).Scan(&services)
	if err == nil {
		return services, nil
	}
	s.handleRedisError(ctx, err)
	servicesDB, ierr := s.serviceRepo.GetByClientId(ctx, typeService, status, clientId)
	if ierr != nil {
		return nil, ierr
	}
	go s.set(ctx, key, servicesDB, expirationServiceByClientId)
	return servicesDB, nil
}

func (s ServiceRepoDecorator) GetByClientCode(ctx context.Context, typeService enums.TypeService, status int, clientCode string) ([]*domain.Service, *common.Error) {
	key := s.genKeyCacheGetServiceByClientCode(typeService, status, clientCode)
	var services []*domain.Service
	err := s.get(ctx, key).Scan(&services)
	if err == nil {
		return services, nil
	}
	s.handleRedisError(ctx, err)
	servicesDB, ierr := s.serviceRepo.GetByClientCode(ctx, typeService, status, clientCode)
	if ierr != nil {
		return nil, ierr
	}
	go s.set(ctx, key, servicesDB, expirationServiceByClientCode)
	return servicesDB, nil
}
