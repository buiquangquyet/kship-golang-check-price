package domain

import "strconv"

type Price struct {
	Id            int64  `json:"id,omitempty"`
	Code          string `json:"code,omitempty"`
	Name          string `json:"name,omitempty"`
	Image         string `json:"image,omitempty"`
	Description   string `json:"description,omitempty"`
	ClientCode    string `json:"clientCode,omitempty"`
	GroupId       string `json:"groupId,omitempty"`
	InsuranceFee  int64  `json:"insuranceFee,omitempty"`
	TransferFee   int64  `json:"transferFee,omitempty"`
	CodFee        int    `json:"codFee,omitempty"`
	Total         int    `json:"total,omitempty"`
	Fee           int64  `json:"fee,omitempty"`
	ConnFee       int    `json:"connFee,omitempty"`
	CodstFee      int    `json:"codstFee,omitempty"`
	TotalPrice    int    `json:"totalPrice,omitempty"`
	OtherPrice    int    `json:"otherPrice,omitempty"`
	CouponSale    string `json:"couponSale"`
	OldTotalPrice int    `json:"oldTotalPrice,omitempty"`
	Status        bool   `json:"status"`
	Msg           string `json:"msg,omitempty"`
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
