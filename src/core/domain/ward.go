package domain

import (
	"check-price/src/common"
	"context"
	"time"
)

type Ward struct {
	Id            int64     `json:"id"`
	DistrictId    int64     `json:"district_id"`
	Name          string    `json:"name"`
	Code          string    `json:"code"`
	PostalCode    string    `json:"postal_code"`
	KvId          string    `json:"kv_id"`
	VtpId         string    `json:"vtp_id"`
	MappingStatus int       `json:"mapping_status"`
	GhnId         string    `json:"ghn_id"`
	SplId         string    `json:"spl_id"`
	JtId          string    `json:"jt_id"`
	JtMapping     string    `json:"jt_mapping"`
	VnpId         int       `json:"vnp_id"`
	VnpSortCode   string    `json:"vnp_sort_code"`
	KmsId         int64     `json:"kms_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type WardRepo interface {
	GetByKmsId(ctx context.Context, senderWardId int64) (*Ward, *common.Error)
	GetByKvId(ctx context.Context, senderWardId int64) (*Ward, *common.Error)
}

func (Ward) TableName() string {
	return "wards"
}
