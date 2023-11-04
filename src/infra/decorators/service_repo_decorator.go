package decorators

import (
	"check-price/src/common"
	"check-price/src/common/log"
	"check-price/src/core/domain"
	"check-price/src/core/enums"
	"check-price/src/infra/repo"
	"context"
	"encoding/json"
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

func (s ServiceRepoDecorator) GetByCode(ctx context.Context, code string) (*domain.Service, *common.Error) {
	//TODO implement me
	panic("implement me")
}

func NewServiceRepoDecorator(base *baseDecorator, serviceRepo *repo.ServiceRepo) domain.ServiceRepo {
	return &ServiceRepoDecorator{
		baseDecorator: base,
		serviceRepo:   serviceRepo,
	}
}

func (s ServiceRepoDecorator) GetByClientIdAndStatus(ctx context.Context, typeService enums.TypeService, status int, clientId int64) ([]*domain.Service, *common.Error) {
	key := s.genKeyCacheGetServiceByClientId(typeService, status, clientId)
	var services []*domain.Service
	val, err := s.get(ctx, key).Result()
	s.handleRedisError(ctx, err)
	if err == nil {
		err = json.Unmarshal([]byte(val), &services)
		if err == nil {
			return services, nil
		}
		log.Warn(ctx, "unmarshall error")
	}
	servicesDB, ierr := s.serviceRepo.GetByClientIdAndStatus(ctx, typeService, status, clientId)
	if ierr != nil {
		return nil, ierr
	}
	go s.set(ctx, key, servicesDB, expirationServiceByClientId)
	return servicesDB, nil
}

func (s ServiceRepoDecorator) GetByClientCodeAndStatus(ctx context.Context, typeService enums.TypeService, status int, clientCode string) ([]*domain.Service, *common.Error) {
	key := s.genKeyCacheGetServiceByClientCode(typeService, status, clientCode)
	var services []*domain.Service
	val, err := s.get(ctx, key).Result()
	s.handleRedisError(ctx, err)
	if err == nil {
		err = json.Unmarshal([]byte(val), &services)
		if err == nil {
			return services, nil
		}
		log.Warn(ctx, "unmarshall error")
	}
	servicesDB, ierr := s.serviceRepo.GetByClientCodeAndStatus(ctx, typeService, status, clientCode)
	if ierr != nil {
		return nil, ierr
	}
	go s.set(ctx, key, servicesDB, expirationServiceByClientCode)
	return servicesDB, nil
}

func (s ServiceRepoDecorator) GetByClientIdAndCodes(ctx context.Context, typeService enums.TypeService, codes []string, clientId int64) ([]*domain.Service, *common.Error) {
	key := s.genKeyCacheGetServiceByClientIdAndCodes(typeService, codes, clientId)
	var services []*domain.Service
	val, err := s.get(ctx, key).Result()
	s.handleRedisError(ctx, err)
	if err == nil {
		err = json.Unmarshal([]byte(val), &services)
		if err == nil {
			return services, nil
		}
		log.Warn(ctx, "unmarshall error")
	}
	servicesDB, ierr := s.serviceRepo.GetByClientIdAndCodes(ctx, typeService, codes, clientId)
	if ierr != nil {
		return nil, ierr
	}
	go s.set(ctx, key, servicesDB, expirationServiceByClientCode)
	return servicesDB, nil
}
