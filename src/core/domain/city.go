package domain

import (
	"check-price/src/common"
	"context"
	"time"
)

type City struct {
	Id            int64     `json:"id"`
	Zone          string    `json:"zone"`
	Name          string    `json:"name"`
	Code          string    `json:"code"`
	PostalCode    string    `json:"postal_code"`
	KvId          string    `json:"kv_id"`
	VtpId         string    `json:"vtp_id"`
	MappingStatus int       `json:"mapping_status"`
	GhnId         int64     `json:"ghn_id"`
	SplId         string    `json:"spl_id"`
	JtName        string    `json:"jt_name"`
	VnpId         int64     `json:"vnp_id"`
	KmsId         int64     `json:"kms_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type CityRepo interface {
	GetById(ctx context.Context, id int64) (*City, *common.Error)
}

func (City) TableName() string {
	return "cities"
}
