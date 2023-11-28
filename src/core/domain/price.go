package domain

import (
	"check-price/src/core/constant"
	"strconv"
)

type Price struct {
	Id            int64   `json:"id"`
	Code          string  `json:"code"`
	Name          string  `json:"name"`
	Image         string  `json:"image"`
	Description   string  `json:"description"`
	ClientCode    string  `json:"clientCode"`
	GroupId       string  `json:"groupId"`
	InsuranceFee  int64   `json:"insuranceFee"`
	TransferFee   int64   `json:"transferFee"`
	CodFee        int64   `json:"codFee"`
	Total         int64   `json:"total"`
	Fee           int64   `json:"fee"`
	ConnFee       float64 `json:"connFee"`
	CodstFee      int64   `json:"codstFee"`
	CodT0Fee      float64 `json:"codT0Fee"`
	TotalPrice    int64   `json:"totalPrice"`
	OtherFee      int64   `json:"otherFee"`
	CouponSale    int64   `json:"couponSale"`
	OldTotalPrice int64   `json:"oldTotalPrice"`
	Status        bool    `json:"status"`
	Msg           string  `json:"msg"`
	StatusCodT0   bool    `json:"status_codT0,omitempty"`
	MessageCodT0  string  `json:"message_codT0,omitempty"`
}

func (p *Price) SetCode(code string) *Price {
	p.Code = code
	return p
}

func (p *Price) SetClientInfo(client *Client) *Price {
	p.ClientCode = client.Code
	p.Image = client.Image
	return p
}

func (p *Price) SetServiceInfo(service *Service) *Price {
	p.Id = service.Id
	p.GroupId = strconv.Itoa(service.GroupId)
	p.Name = service.Name
	p.Description = service.Description
	return p
}

func (p *Price) SetConnFee(connFee float64) *Price {
	p.ConnFee = connFee
	return p
}

func (p *Price) SetCodStFee(shop *Shop, cod int64) *Price {
	var codStFee int64
	isShopType := shop.Type == constant.ShopVip
	for i := 0; i < constant.MaxLevel; i++ {
		if constant.CodLevelMin[i] <= cod && cod <= constant.CodLevelMax[i] {
			if isShopType {
				codStFee = constant.PriceVip[i]
			} else {
				codStFee = constant.PriceNormal[i]
			}
		}
	}
	p.CodstFee = codStFee
	return p
}

func (p *Price) SetCodT0Info(status bool, message string, codStFee float64) *Price {
	p.StatusCodT0 = status
	p.MessageCodT0 = message
	p.CodT0Fee = codStFee
	return p
}

func (p *Price) SetTotalPrice(payer string) *Price {
	total := p.CodstFee + int64(p.ConnFee) + int64(p.CodT0Fee)
	totalFeeExtraService := total
	if payer == constant.PaymentByFrom {
		total = p.Fee + total
	}
	p.Total = total
	p.TotalPrice += totalFeeExtraService
	return p
}

func (p *Price) SetCouponInfo(discountVoucher int64) *Price {
	totalPrice := p.TotalPrice
	if discountVoucher > totalPrice {
		discountVoucher = totalPrice
	}
	p.CouponSale = discountVoucher
	p.OldTotalPrice = totalPrice
	p.TotalPrice = totalPrice - discountVoucher
	p.Total = p.TotalPrice
	p.OtherFee = totalPrice - (p.TransferFee + p.InsuranceFee + p.CodFee)
	return p
}

func (p *Price) SetOtherFee() *Price {
	p.OtherFee = p.TotalPrice - (p.TransferFee + p.InsuranceFee + p.CodFee) + p.CouponSale +
		p.CodstFee + int64(p.ConnFee) + int64(p.CodT0Fee)
	return p
}

func (p *Price) HandleFeeGHTK() *Price {
	p.Total = p.InsuranceFee - (p.CodstFee + int64(p.ConnFee) + int64(p.CodT0Fee))
	return p
}
