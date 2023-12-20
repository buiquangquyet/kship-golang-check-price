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
	expirationDistrictByKmsId = 30 * 24 * time.Hour
	expirationDistrictByKvId  = 30 * 24 * time.Hour
	expirationDistrictById    = 30 * 24 * time.Hour
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
	val, err := d.get(ctx, key).Result()
	d.handleRedisError(ctx, err)
	if err == nil {
		err = json.Unmarshal([]byte(val), &district)
		if err == nil {
			return &district, nil
		}
		log.Warn(ctx, "unmarshall error")
	}
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
	val, err := d.get(ctx, key).Result()
	d.handleRedisError(ctx, err)
	if err == nil {
		err = json.Unmarshal([]byte(val), &district)
		if err == nil {
			return &district, nil
		}
		log.Warn(ctx, "unmarshall error")
	}
	districtDB, ierr := d.districtRepo.GetByKvId(ctx, kvId)
	if ierr != nil {
		return nil, ierr
	}
	go d.set(ctx, key, districtDB, expirationDistrictByKvId)
	return districtDB, nil
}

func (d DistrictRepoDecorator) GetById(ctx context.Context, id int64) (*domain.District, *common.Error) {
	key := d.genKeyCacheGetDistrictById(id)
	var district domain.District
	val, err := d.get(ctx, key).Result()
	d.handleRedisError(ctx, err)
	if err == nil {
		err = json.Unmarshal([]byte(val), &district)
		if err == nil {
			return &district, nil
		}
		log.Warn(ctx, "unmarshall error")
	}
	d.handleRedisError(ctx, err)
	districtDB, ierr := d.districtRepo.GetById(ctx, id)
	if ierr != nil {
		return nil, ierr
	}
	go d.set(ctx, key, districtDB, expirationDistrictById)
	return districtDB, nil
}
