package dto

import "check-price/src/core/domain"

type Address struct {
	PickProvince *domain.City
	PickDistrict *domain.District
	PickWard     *domain.Ward

	ReceiverProvince *domain.City
	ReceiverDistrict *domain.District
	ReceiverWard     *domain.Ward
}
