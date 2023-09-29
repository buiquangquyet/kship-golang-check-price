package ghtk

import (
	"check-price/src/common"
	"check-price/src/common/configs"
	"check-price/src/core/domain"
	"check-price/src/infra/external"
	"context"
	"github.com/imroc/req/v3"
	"time"
)

const (
	timeoutGHTK = 5 * time.Second
)

type GHTKExtService struct {
	*external.BaseClient
	client *req.Client
	cf     *configs.GHTK
}

func NewGHTKExtService(base *external.BaseClient) *GHTKExtService {
	cf := configs.Get().ExtService.GHTK
	cli := req.C().SetBaseURL(cf.Host).SetTimeout(timeoutGHTK)
	cli.SetCommonHeaders(map[string]string{
		"Content-Type": "application/json",
	})
	base.SetTracer(cli)
	return &GHTKExtService{
		BaseClient: base,
		client:     cli,
		cf:         cf,
	}
}

func (g *GHTKExtService) Connect(ctx context.Context, shopCode string) string {
	// get token from cache
	return g.cf.Token
}

func (g *GHTKExtService) GetPriceFromDelivery(ctx context.Context, service string) (*domain.Price, *common.Error) {
	return &domain.Price{}, nil
}
