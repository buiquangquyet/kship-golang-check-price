package strategy

import (
	"check-price/src/common"
	"check-price/src/common/log"
	"check-price/src/core/constant"
	"check-price/src/core/domain"
	"check-price/src/core/dto"
	"check-price/src/helpers"
	"check-price/src/present/httpui/request"
	"context"
)

type BaseStrategy struct {
	wardRepo     domain.WardRepo
	districtRepo domain.DistrictRepo
	cityRepo     domain.CityRepo
}

func NewBaseStrategy(
	wardRepo domain.WardRepo,
	districtRepo domain.DistrictRepo,
	cityRepo domain.CityRepo,
) *BaseStrategy {
	return &BaseStrategy{
		wardRepo:     wardRepo,
		districtRepo: districtRepo,
		cityRepo:     cityRepo,
	}
}

func (g *BaseStrategy) GetAddress(ctx context.Context, req *request.GetPriceRequest) (*dto.Address, *common.Error) {
	isVer2 := req.VersionLocation == constant.VersionLocation2
	ierr := common.ErrBadRequest(ctx)
	var pickWard *domain.Ward
	var receiverWard *domain.Ward
	if isVer2 {
		pickWard, ierr = g.wardRepo.GetByKmsId(ctx, req.SenderWardId)
		if helpers.IsInternalError(ierr) {
			log.Error(ctx, ierr.Error())
			return nil, ierr
		}
	} else {
		pickWard, ierr = g.wardRepo.GetByKvId(ctx, req.SenderWardId)
		if helpers.IsInternalError(ierr) {
			log.Error(ctx, ierr.Error())
			return nil, ierr
		}
	}
	if ierr != nil {
		return nil, ierr.SetCode(4004)
	}

	if isVer2 {
		receiverWard, ierr = g.wardRepo.GetByKmsId(ctx, req.ReceiverWardId)
		if helpers.IsInternalError(ierr) {
			log.Error(ctx, ierr.Error())
			return nil, ierr
		}
	} else {
		receiverWard, ierr = g.wardRepo.GetByKvId(ctx, req.ReceiverWardId)
		if helpers.IsInternalError(ierr) {
			log.Error(ctx, ierr.Error())
			return nil, ierr
		}
	}
	if ierr != nil {
		return nil, ierr.SetCode(4006)
	}

	pickDistrict, err := g.districtRepo.GetById(ctx, pickWard.DistrictId)
	if err != nil {
		log.Error(ctx, err.Error())
		return nil, err
	}
	pickProvince, err := g.cityRepo.GetById(ctx, pickDistrict.CityId)
	if err != nil {
		log.Error(ctx, err.Error())
		return nil, err
	}
	receiverDistrict, err := g.districtRepo.GetById(ctx, receiverWard.DistrictId)
	if err != nil {
		log.Error(ctx, err.Error())
		return nil, err
	}
	receiverProvince, err := g.cityRepo.GetById(ctx, receiverDistrict.CityId)
	if err != nil {
		log.Error(ctx, err.Error())
		return nil, err
	}
	return &dto.Address{
		PickProvince:     pickProvince,
		PickDistrict:     pickDistrict,
		PickWard:         pickWard,
		ReceiverProvince: receiverProvince,
		ReceiverDistrict: receiverDistrict,
		ReceiverWard:     receiverWard,
	}, nil
}
