package dto

import "time"

type Voucher struct {
	StatusCode    int
	Id            int64
	Code          string
	Type          int
	DiscountType  int
	DiscountValue int64
	TimeStart     time.Time
	ExpireTime    time.Time
	Description   string
	Quantity      int
	TotalUsed     int
	IsUsedUp      int
	Active        int
	RetailerId    string
	ClientIds     string
	CreatedBy     string
	UpdatedBy     string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
