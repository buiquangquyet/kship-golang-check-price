package domain

import (
	"check-price/src/common"
	"context"
	"time"
)

type Setting struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SettingRepo interface {
	GetByName(ctx context.Context, name string) (*Setting, *common.Error)
}

func (Setting) TableName() string {
	return "settings"
}
