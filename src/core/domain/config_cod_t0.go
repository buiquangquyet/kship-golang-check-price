package domain

import (
	"check-price/src/common"
	"context"
	"time"
)

type ConfigCodT0 struct {
	Id               int64     `json:"id"`
	CodFrom          int       `json:"cod_from"`
	CodTo            int       `json:"cod_to"`
	Type             int       `json:"type"`
	Value            float64   `json:"value"`
	ClientId         int64     `json:"client_id"`
	ServicesPossible string    `json:"services_possible"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type ConfigCodT0Repo interface {
	GetByCodAndClientId(ctx context.Context, cod int64, clientId int64) ([]*ConfigCodT0, *common.Error)
}
