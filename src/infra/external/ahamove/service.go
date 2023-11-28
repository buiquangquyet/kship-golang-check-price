package ahamoveext

import (
	"check-price/src/common"
	"check-price/src/common/configs"
	"check-price/src/common/log"
	"check-price/src/core/domain"
	"check-price/src/core/param"
	"check-price/src/helpers"
	"check-price/src/infra/external"
	"context"
	"fmt"
	"github.com/imroc/req/v3"
	"github.com/redis/go-redis/v9"
	"time"
)

const (
	deliveryCode   = "AHAMOVE"
	timeoutAhaMove = 5 * time.Second
	expireToken    = 24 * time.Hour

	loginPath      = "/v1/partner/register_account"
	checkPricePath = "/v2/order/estimated_fee"
)

type Service struct {
	*external.BaseClient
	client *req.Client
	key    string
}

func NewService(base *external.BaseClient) *Service {
	cf := configs.Get().ExtService.AHAMOVE
	cli := req.C().SetBaseURL(cf.Host).SetTimeout(timeoutAhaMove)
	cli.SetCommonHeaders(map[string]string{
		"Content-Type": "application/json",
	})
	base.SetTracer(cli)
	return &Service{
		BaseClient: base,
		client:     cli,
		key:        cf.Key,
	}
}

func (g *Service) api(ctx context.Context) *req.Request {
	return g.client.R().SetContext(ctx)
}

func (g *Service) CheckPrice(ctx context.Context, shop *domain.Shop, p *param.GetPriceAhaMoveParam) ([]*domain.Price, *common.Error) {
	token, fromCache, err := g.getToken(ctx, shop, true)
	if err != nil {
		return nil, err
	}
	result, ierr := g.checkPrice(ctx, token, p)
	if ierr != nil {
		if fromCache && helpers.IsUnauthorizedError(ierr) {
			// retry once
			newToken, _, err := g.getToken(ctx, shop, false)
			if err != nil {
				return nil, err
			}
			return g.checkPrice(ctx, newToken, p)
		}
		return nil, ierr
	}
	return result, nil
}

func (g *Service) checkPrice(ctx context.Context, token string, p *param.GetPriceAhaMoveParam) ([]*domain.Price, *common.Error) {
	var output []*PriceOuput
	var outputErr OutputErr
	resp, err := g.api(ctx).
		SetBody(newGetPriceInput(token, p)).
		SetSuccessResult(&output).
		SetErrorResult(&outputErr).
		Get(checkPricePath)
	if err != nil {
		return nil, common.ErrSystemError(ctx, err.Error()).SetSource(common.SourceAHAMOVEService)
	}

	if resp.IsErrorState() {
		log.Debug(ctx, "Call AHAMOVE failed with body: %+v", output)
		detail := fmt.Sprintf("http: [%d], resp: [%s]", resp.StatusCode, resp.String())
		return nil, common.ErrSystemError(ctx, detail).SetSource(common.SourceAHAMOVEService)
	}
	prices := make([]*domain.Price, len(output))
	for i, p := range output {
		prices[i] = p.ToDomain()
	}
	return prices, nil
}

func (g *Service) getToken(ctx context.Context, shop *domain.Shop, allowFromCache bool) (string, bool, *common.Error) {
	if allowFromCache {
		token, err := g.GetTokenFromCache(ctx, deliveryCode, shop)
		if err == nil && token != "" {
			return token, true, nil
		}
		if err != nil && err != redis.Nil {
			log.Warn(ctx, "Get AHAMOVE Token of shop %s, error: %s", shop.Code, err.Error())
		}
	}
	newToken, err := g.newToken(ctx, shop)
	if err != nil {
		return "", false, err
	}
	go g.StoreToken(common.Detach(ctx), deliveryCode, shop, newToken, expireToken)
	return newToken, false, nil
}

func (g *Service) newToken(ctx context.Context, shop *domain.Shop) (string, *common.Error) {
	var output LoginOutput
	var outputErr OutputErr
	resp, err := g.api(ctx).
		SetQueryParams(map[string]string{
			"name":    "Kiotviet.vn",
			"api_key": g.key,
			"mobile":  shop.Phone,
		}).
		SetSuccessResult(&output).
		SetErrorResult(&outputErr).
		Get(loginPath)
	if err != nil {
		return "", common.ErrSystemError(ctx, err.Error())
	}

	if resp.IsErrorState() {
		log.Debug(ctx, "Call AHAMOVE failed with body: %+v", output)
		detail := fmt.Sprintf("http: [%d], resp: [%s]", resp.StatusCode, resp.String())
		return "", common.ErrSystemError(ctx, detail).SetSource(common.SourceAHAMOVEService)
	}
	return output.Token, nil
}

func (g *Service) handleError(_ context.Context) *common.Error {
	//Todo code

	return nil
}
