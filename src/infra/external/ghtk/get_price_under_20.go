package ghtk

import (
	"check-price/src/core/domain"
	"check-price/src/core/dto"
)

type GetPriceUnder20Input struct {
	Address      string `json:"address"`
	PickProvince string `json:"pick_province"`
	PickDistrict string `json:"pick_district"`
	PickWard     string `json:"pick_ward"`
	Province     string `json:"province"`
	District     string `json:"district"`
	Ward         string `json:"ward"`
	Weight       int64  `json:"weight"`
	Value        int64  `json:"value"`
	Transport    string `json:"transport"`
	Tags         []int  `json:"tags"`
}

func newGetPriceUnder20Input(serviceCode string, p *dto.GetPriceInputDto) *GetPriceUnder20Input {
	return &GetPriceUnder20Input{
		Address:      p.Address,
		PickProvince: p.PickProvince,
		PickDistrict: p.PickDistrict,
		PickWard:     p.PickWard,
		Province:     p.ReceiverProvince,
		District:     p.ReceiverDistrict,
		Ward:         p.ReceiverWard,
		Weight:       p.Weight,
		Value:        p.Value,
		Transport:    serviceCode,
		Tags:         p.Tags,
	}
}

type GetPriceUnder20Output struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Fee     struct {
		Name         string `json:"name"`
		Fee          int64  `json:"fee"`
		InsuranceFee int64  `json:"insurance_fee"`
		IncludeVat   string `json:"include_vat"`
		CostId       string `json:"cost_id"`
		DeliveryType string `json:"delivery_type"`
		A            string `json:"a"`
		Dt           string `json:"dt"`
		ExtFees      []*Fee `json:"ext_fees"`
		ShipFeeOnly  int64  `json:"ship_fee_only"`
		PromotionKey string `json:"promotion_key"`
		Delivery     bool   `json:"delivery"`
	} `json:"fee"`
}

func (g *GetPriceUnder20Output) ToDomainPrice() *domain.Price {
	return &domain.Price{
		Id:            0,
		Code:          "",
		Name:          "",
		Image:         "",
		Description:   "",
		ClientCode:    "",
		GroupId:       "",
		InsuranceFee:  g.Fee.InsuranceFee,
		TransferFee:   g.Fee.ShipFeeOnly,
		CodeFee:       0,
		Total:         0,
		Fee:           g.Fee.Fee,
		ConnFee:       0,
		CodstFee:      0,
		TotalPrice:    0,
		OtherPrice:    0,
		OldTotalPrice: 0,
		Status:        g.Fee.Delivery,
		Msg:           "",
	}
}
