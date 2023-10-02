package domain

import (
	"check-price/src/common"
	"context"
	"time"
)

type District struct {
	Id            int64     `json:"id"`
	CityId        int64     `json:"city_id"`
	Name          string    `json:"name"`
	Code          string    `json:"code"`
	PostalCode    string    `json:"postal_code"`
	KvId          string    `json:"kv_id"`
	VtpId         string    `json:"vtp_id"`
	MappingStatus int       `json:"mapping_status"`
	GhnId         int       `json:"ghn_id"`
	SplId         string    `json:"spl_id"`
	JtName        string    `json:"jt_name"`
	VnpId         int       `json:"vnp_id"`
	KmsId         int       `json:"kms_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type DistrictRepo interface {
	GetByKmsId(ctx context.Context, senderLocationId int64) (*District, *common.Error)
	GetByKvId(ctx context.Context, senderLocationId int64) (*District, *common.Error)
}

func (District) TableName() string {
	return "districts"
}
