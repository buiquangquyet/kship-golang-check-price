package domain

import (
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
	CodFee        int     `json:"codFee"`
	Total         int     `json:"total"`
	Fee           int64   `json:"fee"`
	ConnFee       int     `json:"connFee"`
	CodstFee      int64   `json:"codstFee"`
	CodT0Fee      float64 `json:"codT0Fee"`
	TotalPrice    int     `json:"totalPrice"`
	OtherPrice    int     `json:"otherPrice"`
	CouponSale    int64   `json:"couponSale"`
	OldTotalPrice int     `json:"oldTotalPrice"`
	Status        bool    `json:"status"`
	Msg           string  `json:"msg"`
	StatusCodT0   bool    `json:"status_codT0"`
	MessageCodT0  string  `json:"message_codT0"`
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

func (p *Price) SetCodT0Info(status bool, message string, codStFee float64) *Price {
	p.StatusCodT0 = status
	p.MessageCodT0 = message
	p.CodT0Fee = codStFee
	return p
}
