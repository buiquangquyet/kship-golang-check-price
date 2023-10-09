package decorators

import (
	"check-price/src/common"
	"check-price/src/core/domain"
	"check-price/src/infra/repo"
	"context"
	"github.com/go-redis/redis/v8"
)

type DistrictRepoDecorator struct {
	*baseDecorator
	cache        redis.UniversalClient
	districtRepo *repo.DistrictRepo
}

func NewDistrictRepoDecorator(base *baseDecorator, districtRepo *repo.DistrictRepo, cache redis.UniversalClient) domain.DistrictRepo {
	return &DistrictRepoDecorator{
		baseDecorator: base,
		cache:         cache,
		districtRepo:  districtRepo,
	}
}

func (d DistrictRepoDecorator) GetByKmsId(ctx context.Context, senderLocationId int64) (*domain.District, *common.Error) {
	//TODO implement me
	panic("implement me")
}

func (d DistrictRepoDecorator) GetByKvId(ctx context.Context, senderLocationId int64) (*domain.District, *common.Error) {
	//TODO implement me
	panic("implement me")
}

func (d DistrictRepoDecorator) GetById(ctx context.Context, id int64) (*domain.District, *common.Error) {
	//TODO implement me
	panic("implement me")
}
