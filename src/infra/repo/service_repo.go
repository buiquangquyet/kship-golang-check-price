package repo

import (
	"check-price/src/common"
	"check-price/src/core/domain"
	"check-price/src/core/enums"
	"check-price/src/helpers"
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

func NewServiceRepo(base *baseRepo) *ServiceRepo {
	return &ServiceRepo{
		base,
	}
}

type ServiceRepo struct {
	*baseRepo
}

func (s ServiceRepo) GetByClientIdAndStatus(ctx context.Context, typeService enums.TypeService, status int, clientId int64) ([]*domain.Service, *common.Error) {
	services := make([]*domain.Service, 0)
	conds := []clause.Expression{
		clause.Eq{Column: "type", Value: typeService.ToInt()},
		clause.Eq{Column: "client_id", Value: clientId},
		clause.Eq{Column: "status", Value: status},
	}
	if err := s.db.WithContext(ctx).Clauses(conds...).Find(&services).Error; err != nil {
		return nil, s.returnError(ctx, err)
	}
	return services, nil
}

func (s ServiceRepo) GetByClientCodeAndStatus(ctx context.Context, typeService enums.TypeService, status int, clientCode string) ([]*domain.Service, *common.Error) {
	services := make([]*domain.Service, 0)
	conds := []clause.Expression{
		clause.Eq{Column: "type", Value: typeService.ToInt()},
		clause.Like{Column: "clients_possible", Value: "%" + strings.TrimSpace(clientCode) + "%"},
		clause.Eq{Column: "status", Value: status},
	}
	if err := s.db.WithContext(ctx).Clauses(conds...).Find(&services).Error; err != nil {
		return nil, s.returnError(ctx, err)
	}
	return services, nil
}

func (s ServiceRepo) GetByClientIdAndCodes(ctx context.Context, typeService enums.TypeService, code []string, clientId int64) ([]*domain.Service, *common.Error) {
	services := make([]*domain.Service, 0)
	codeInterface := helpers.ConvertTypesToInterfaces(code)
	conds := []clause.Expression{
		clause.Eq{Column: "client_id", Value: clientId},
		clause.IN{Column: "code", Values: codeInterface},
		clause.Eq{Column: "type", Value: typeService.ToInt()},
	}
	if err := s.db.WithContext(ctx).Clauses(conds...).Find(&services).Error; err != nil {
		return nil, s.returnError(ctx, err)
	}
	return services, nil
}

func (s ServiceRepo) GetByCode(ctx context.Context, code string) (*domain.Service, *common.Error) {
	service := &domain.Service{}
	cond := clause.Eq{Column: "code", Value: code}
	if err := s.db.WithContext(ctx).Clauses(cond).Take(service).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound(ctx, "service", "not found").SetSource(common.SourceInfraService)
		}
		return nil, s.returnError(ctx, err)
	}
	return service, nil
}
