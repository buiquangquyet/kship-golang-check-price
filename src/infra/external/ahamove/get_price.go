package ahamoveext

import (
	"check-price/src/core/domain"
)

type GetPriceInput struct {
	Path          []*Path    `json:"path"`
	PaymentMethod string     `json:"payment_method"`
	PromoCode     string     `json:"promo_code"`
	OrderTime     int        `json:"order_time"`
	Services      []*Service `json:"services"`
	Token         string     `json:"token"`
}

type Path struct {
	Address string
	Name    string
	Mobile  string
	Cod     int
}

type Service struct {
	Id       int
	Requests []string
}

func newGetPriceInput(token string) *GetPriceInput {
	return &GetPriceInput{
		Path:          nil,
		PaymentMethod: "",
		PromoCode:     "",
		OrderTime:     0,
		Services:      nil,
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
