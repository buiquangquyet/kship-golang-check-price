package dto

import (
	"check-price/src/core/domain"
	"check-price/src/present/httpui/request"
)

type ValidatePrice struct {
	ClientCode       string
	Shop             *domain.Shop
	RetailerId       int64
	ValidateLocation *ValidateLocation
	Services         []*ValidateService
	ExtraServices    []*ValidateExtraService
}

type ValidateLocation struct {
	VersionLocation    int
	ReceiverLocationId int64
	SenderWardId       int64
	ReceiverWardId     int64
}

type ValidateService struct {
	Code          string
	OldTotalPrice string
}

type ValidateExtraService struct {
	Code     string
	Value    string
	ViewType string
	Name     string
}

func NewValidatePrice(shop *domain.Shop, req *request.GetPriceReRequest, token *request.TokenInfo) *ValidatePrice {
	services := make([]*ValidateService, len(req.Services))
	for i, s := range req.Services {
		services[i] = &ValidateService{
			Code:          s.Code,
			OldTotalPrice: s.OldTotalPrice,
		}
	}
	extraServices := make([]*ValidateExtraService, len(req.ExtraService))
	for i, s := range req.ExtraService {
		extraServices[i] = &ValidateExtraService{
			Code:     s.Code,
			Value:    s.Value,
			ViewType: s.ViewType,
			Name:     s.Name,
		}
	}
	return &ValidatePrice{
		ClientCode: req.ClientCode,
		Shop:       shop,
		RetailerId: token.RetailerId,
		ValidateLocation: &ValidateLocation{
			VersionLocation:    req.VersionLocation,
			ReceiverLocationId: req.ReceiverLocationId,
			SenderWardId:       req.SenderWardId,
			ReceiverWardId:     req.ReceiverWardId,
		},
		Services:      services,
		ExtraServices: extraServices,
	}
}
