package ahamove

import (
	"check-price/src/common"
	"check-price/src/common/configs"
	"check-price/src/common/log"
	"check-price/src/core/domain"
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
	codeSuccess    = "SUCCESS"
	timeoutAhaMove = 5 * time.Second
	expireToken    = 24 * time.Hour

	loginPath      = "/v1/partner/register_account"
	checkPricePath = "/v2/order/estimated_fee"
)

type AhaMoveExtService struct {
	*external.BaseClient
	client *req.Client
	key    string
}

func NewAhaMoveExtService(base *external.BaseClient) *AhaMoveExtService {
	cf := configs.Get().ExtService.AHAMOVE
	cli := req.C().SetBaseURL(cf.Host).SetTimeout(timeoutAhaMove)
	cli.SetCommonHeaders(map[string]string{
		"Content-Type": "application/json",
	})
	base.SetTracer(cli)
	return &AhaMoveExtService{
		BaseClient: base,
		client:     cli,
		key:        cf.Key,
	}
}

func (g *AhaMoveExtService) api(ctx context.Context) *req.Request {
	return g.client.R().SetContext(ctx)
}

func (g *AhaMoveExtService) CheckPrice(ctx context.Context, shop *domain.Shop) ([]*domain.Price, *common.Error) {
	token, fromCache, err := g.getToken(ctx, shop, true)
	if err != nil {
		return nil, err
	}

	result, ierr := g.checkPrice(ctx, token)
	if ierr != nil {
		if fromCache && helpers.IsUnauthorizedError(err) {
			// retry once
			newToken, _, err := g.getToken(ctx, shop, false)
			if err != nil {
				return nil, err
			}
			return g.checkPrice(ctx, newToken)
		} else {
			return nil, ierr
		}
	}
	return result, nil
}

func (g *AhaMoveExtService) checkPrice(ctx context.Context, token string) ([]*domain.Price, *common.Error) {
	var output []*PriceOuput
	var outputErr OutputErr
	resp, err := g.api(ctx).
		SetBody(newGetPriceInput(token)).
		SetSuccessResult(&output).
		SetErrorResult(&outputErr).
		Get(checkPricePath)
	if err != nil {
		return nil, common.ErrSystemError(ctx, err.Error()).SetSource(common.SourceGHTKService)
	}

	if resp.IsErrorState() {
		log.Debug(ctx, "Call GetCompany MISA failed with body: %+v", output)
		//Todo consider error code
		//detail := fmt.Sprintf("http: [%d], resp: [%s]", resp.StatusCode, resp.String())
		//return nil, g.handleError(ctx, resp.StatusCode, &output).SetSource(common.SourceGHTKService).SetDetail(detail)
	}
	prices := make([]*domain.Price, len(output))
	for i, p := range output {
		prices[i] = p.ToDomain()
	}
	return prices, nil
}

func (g *AhaMoveExtService) getToken(ctx context.Context, shop *domain.Shop, allowFromCache bool) (string, bool, *common.Error) {
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

func (g *AhaMoveExtService) newToken(ctx context.Context, shop *domain.Shop) (string, *common.Error) {
	var output LoginOutput
	var outputErr OutputErr
	resp, err := g.api(ctx).
		SetPathParams(map[string]string{
			"name":    "Kiotviet.vn",
			"api_key": g.key,
			"mobile":  shop.Phone,
		}).
		SetSuccessResult(&output).
		SetErrorResult(&outputErr).
		Post(loginPath)
	if err != nil {
		return "", common.ErrSystemError(ctx, err.Error())
	}

	if resp.IsErrorState() {
		log.Debug(ctx, "Call ghtk failed with body: %+v", output)
		detail := fmt.Sprintf("http: [%d], resp: [%s]", resp.StatusCode, resp.String())
		return "", common.ErrSystemError(ctx, detail).SetSource(common.SourceGHTKService)
	}
	return output.Token, nil
}

func (g *AhaMoveExtService) success(code string) bool {
	return code == codeSuccess
}

func (g *AhaMoveExtService) handleError(ctx context.Context) *common.Error {
	//Todo code
	return nil
}
