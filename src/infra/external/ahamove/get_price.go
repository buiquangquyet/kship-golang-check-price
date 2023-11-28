package ahamoveext

import (
	"check-price/src/core/constant"
	"check-price/src/core/domain"
	"check-price/src/core/param"
)

type GetPriceInput struct {
	Path          [2]*Path        `json:"path"`
	PaymentMethod string          `json:"payment_method"`
	PromoCode     string          `json:"promo_code"`
	OrderTime     int64           `json:"order_time"`
	Services      []*ServiceInput `json:"services"`
	Token         string          `json:"token"`
}

type Path struct {
	Address string `json:"address,omitempty"`
	Name    string `json:"name,omitempty"`
	Mobile  string `json:"mobile,omitempty"`
	Cod     int64  `json:"cod,omitempty"`
}

type ServiceInput struct {
	Id       string     `json:"_id,omitempty"`
	Requests []*Request `json:"requests"`
}

type Request struct {
	Id       string `json:"_id,omitempty"`
	Num      int    `json:"num,omitempty"`
	TierCode string `json:"tier_code,omitempty"`
}

func newGetPriceInput(token string, p *param.GetPriceAhaMoveParam) *GetPriceInput {
	services := make([]*ServiceInput, len(p.Services))
	for i, s := range p.Services {
		requests := make([]*Request, len(s.Requests))
		for j, request := range s.Requests {
			requests[j] = &Request{
				Id:       request.Id,
				Num:      request.Num,
				TierCode: request.TierCode,
			}
		}
		services[i] = &ServiceInput{
			Id:       s.Id,
			Requests: requests,
		}
	}
	cod := p.Path[1].Cod
	if constant.IsDevEnv() {
		cod = 0
	}
	return &GetPriceInput{
		Path: [2]*Path{
			{Address: p.Path[0].Address},
			{Address: p.Path[1].Address, Cod: cod},
		},
		PaymentMethod: p.PaymentMethod,
		PromoCode:     p.PromoCode,
		OrderTime:     p.OrderTime,
		Services:      services,
		Token:         token,
	}
}

type PriceOuput struct {
	VoucherDiscount int64  `json:"voucher_discount"`
	Discount        int64  `json:"discount"`
	DistancePrice   int64  `json:"distance_price"`
	DistanceFee     int64  `json:"distance_fee"`
	StoppointPrice  int64  `json:"stoppoint_price"`
	StopFee         int64  `json:"stop_fee"`
	Vat             int64  `json:"vat"`
	VatFee          int64  `json:"vat_fee"`
	SubtotalPrice   int64  `json:"subtotal_price"`
	TotalFee        int64  `json:"total_fee"`
	TotalPrice      int64  `json:"total_price"`
	Id              string `json:"_id"`
}

func (g *PriceOuput) ToDomain() *domain.Price {
	return &domain.Price{
		Code:          g.Id,
		TransferFee:   g.DistancePrice + g.StoppointPrice + g.VatFee,
		TotalPrice:    g.TotalPrice,
		CouponSale:    g.VoucherDiscount,
		OldTotalPrice: g.SubtotalPrice,
		Status:        true,
		Msg:           "Check price success",
	}
}
