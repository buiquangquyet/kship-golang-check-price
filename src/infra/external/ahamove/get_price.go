package ahamoveext

import (
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
	Address string
	Name    string
	Mobile  string
	Cod     int64
}

type ServiceInput struct {
	Id       string
	Requests []string
}

func newGetPriceInput(token string, p *param.GetPriceAhaMoveParam) *GetPriceInput {
	services := make([]*ServiceInput, len(p.Services))
	for i, s := range p.Services {
		services[i] = &ServiceInput{
			Id:       s.Id,
			Requests: s.Requests,
		}
	}
	return &GetPriceInput{
		Path: [2]*Path{
			{Address: p.Path[0].Address},
			{Address: p.Path[1].Address, Cod: p.Path[1].Cod},
		},
		PaymentMethod: p.PaymentMethod,
		PromoCode:     p.PromoCode,
		OrderTime:     p.OrderTime,
		Services:      services,
		Token:         token,
	}
}

type PriceOuput struct {
	VoucherDiscount int64 `json:"voucher_discount"`
	Discount        int64 `json:"discount"`
	DistancePrice   int64 `json:"distance_price"`
	DistanceFee     int64 `json:"distance_fee"`
	StoppointPrice  int64 `json:"stoppoint_price"`
	StopFee         int64 `json:"stop_fee"`
	Vat             int64 `json:"vat"`
	VatFee          int64 `json:"vat_fee"`
	SubtotalPrice   int64 `json:"subtotal_price"`
	TotalFee        int64 `json:"total_fee"`
	Surchage        int64 `json:"surchage"`
}

func (g *PriceOuput) ToDomain() *domain.Price {
	return &domain.Price{
		Id:            0,
		Code:          "",
		Name:          "",
		Image:         "",
		Description:   "",
		ClientCode:    "",
		GroupId:       "",
		InsuranceFee:  0,
		TransferFee:   0,
		CodFee:        0,
		Total:         0,
		Fee:           0,
		ConnFee:       0,
		CodstFee:      0,
		CodT0Fee:      0,
		TotalPrice:    0,
		OtherFee:      0,
		CouponSale:    0,
		OldTotalPrice: 0,
		Status:        false,
		Msg:           "",
		StatusCodT0:   false,
		MessageCodT0:  "",
	}
}
