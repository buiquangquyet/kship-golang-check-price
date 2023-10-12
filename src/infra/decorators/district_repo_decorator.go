package decorators

import (
	"check-price/src/common"
	"check-price/src/core/domain"
	"check-price/src/infra/repo"
	"context"
	"time"
)

const (
	expirationDistrictByKmsId = 30 * 24 * time.Hour
	expirationDistrictByKvId  = 30 * 24 * time.Hour
)

type DistrictRepoDecorator struct {
	*baseDecorator
	districtRepo *repo.DistrictRepo
}

func NewDistrictRepoDecorator(base *baseDecorator, districtRepo *repo.DistrictRepo) domain.DistrictRepo {
	return &DistrictRepoDecorator{
		baseDecorator: base,
		districtRepo:  districtRepo,
	}
}

func (d DistrictRepoDecorator) GetByKmsId(ctx context.Context, kmsId int64) (*domain.District, *common.Error) {
	key := d.genKeyCacheGetDistrictByKmsId(kmsId)
	var district domain.District
	err := d.get(ctx, key).Scan(&district)
	if err != nil {
		return &district, nil
	}
	d.handleRedisError(ctx, err)
	districtDB, ierr := d.districtRepo.GetByKmsId(ctx, kmsId)
	if ierr != nil {
		return nil, ierr
	}
	go d.set(ctx, key, districtDB, expirationDistrictByKmsId)
	return districtDB, nil
}

func (d DistrictRepoDecorator) GetByKvId(ctx context.Context, kvId int64) (*domain.District, *common.Error) {
	key := d.genKeyCacheGetDistrictByKvId(kvId)
	var district domain.District
	err := d.get(ctx, key).Scan(&district)
	if err != nil {
		return &district, nil
	}
	d.handleRedisError(ctx, err)
	shopDB, ierr := d.districtRepo.GetByKvId(ctx, kvId)
	if ierr != nil {
		return nil, ierr
	}
	go d.set(ctx, key, shopDB, expirationDistrictByKvId)
	return shopDB, nil
}

func (d DistrictRepoDecorator) GetById(ctx context.Context, id int64) (*domain.District, *common.Error) {
	//TODO implement me
	panic("implement me")
}
