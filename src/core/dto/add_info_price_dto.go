package dto

import (
	"check-price/src/core/constant"
	"check-price/src/core/domain"
	"check-price/src/present/httpui/request"
)

type AddInfoDto struct {
	Shop         *domain.Shop
	Client       *domain.Client
	Cod          int64
	RetailerId   int64
	Payer        string
	Coupon       string
	Services     []*Service
	ExtraService []*ExtraService
}

func NewAddInfoDTO(shop *domain.Shop, client *domain.Client, req *request.GetPriceRequest) *AddInfoDto {
	services := make([]*Service, len(req.Services))
	for i, service := range req.Services {
		services[i] = &Service{
			Code:          service.Code,
			OldTotalPrice: service.OldTotalPrice,
		}
	}
	payer := ""
	coupon := ""
	extraServices := make([]*ExtraService, len(req.ExtraService))
	for i, extraService := range req.ExtraService {
		if extraService.Code == constant.ServiceExtraCodePayment {
			payer = extraService.Value
		}
		if extraService.Code == "Coupon" {
			coupon = extraService.Value
		}
		extraServices[i] = &ExtraService{
			Code:     extraService.Code,
			Value:    extraService.Value,
			ViewType: extraService.ViewType,
			Name:     extraService.Name,
		}
	}
	return &AddInfoDto{
		Shop:         shop,
		Client:       client,
		Cod:          req.MoneyCollection,
		RetailerId:   req.RetailerId,
		Payer:        payer,
		Coupon:       coupon,
		Services:     services,
		ExtraService: extraServices,
	}
}

type Service struct {
	Code          string
	OldTotalPrice string
}

type ExtraService struct {
	Code     string
	Value    string
	ViewType string
	Name     string
}
