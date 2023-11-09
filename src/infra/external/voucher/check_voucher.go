package voucher

import (
	"check-price/src/core/dto"
	"time"
)

type checkVoucherOutput struct {
	Data struct {
		Id            int64     `json:"ID"`
		Code          string    `json:"Code"`
		Type          int       `json:"Type"`
		DiscountType  int       `json:"DiscountType"`
		DiscountValue int64     `json:"DiscountValue"`
		TimeStart     time.Time `json:"TimeStart"`
		ExpireTime    time.Time `json:"ExpireTime"`
		Description   string    `json:"Description"`
		Quantity      int       `json:"Quantity"`
		TotalUsed     int       `json:"TotalUsed"`
		IsUsedUp      int       `json:"IsUsedUp"`
		Active        int       `json:"Active"`
		RetailerId    string    `json:"RetailerId"`
		ClientIds     string    `json:"ClientIds"`
		CreatedBy     string    `json:"CreatedBy"`
		UpdatedBy     string    `json:"UpdatedBy"`
		CreatedAt     time.Time `json:"CreatedAt"`
		UpdatedAt     time.Time `json:"UpdatedAt"`
	} `json:"data"`
	Message    string `json:"message"`
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
}

func (c *checkVoucherOutput) ToDTO() *dto.Voucher {
	d := c.Data
	return &dto.Voucher{
		StatusCode:    c.StatusCode,
		Id:            d.Id,
		Code:          d.Code,
		Type:          d.Type,
		DiscountType:  d.DiscountType,
		DiscountValue: d.DiscountValue,
		TimeStart:     d.TimeStart,
		ExpireTime:    d.ExpireTime,
		Description:   d.Description,
		Quantity:      d.Quantity,
		TotalUsed:     d.IsUsedUp,
		IsUsedUp:      d.IsUsedUp,
		Active:        d.Active,
		RetailerId:    d.RetailerId,
		ClientIds:     d.ClientIds,
		CreatedBy:     d.CreatedBy,
		UpdatedBy:     d.UpdatedBy,
		CreatedAt:     d.CreatedAt,
		UpdatedAt:     d.UpdatedAt,
	}
}
