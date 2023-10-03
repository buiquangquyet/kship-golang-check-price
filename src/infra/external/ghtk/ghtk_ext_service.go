package ghtk

import (
	"check-price/src/common"
	"check-price/src/common/configs"
	"check-price/src/common/log"
	"check-price/src/core/domain"
	"check-price/src/helpers"
	"check-price/src/infra/external"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/imroc/req/v3"
	"time"
)

const (
	deliveryCode = "GHTK"
	codeSuccess  = "SUCCESS"
	timeoutGHTK  = 5 * time.Second
	expireToken  = 24 * 30 * time.Hour

	loginPath = "/services/shops/token"
)

type GHTKExtService struct {
	*external.BaseClient
	client *req.Client
	token  string
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
		token:      cf.Token,
	}
}

func (g *GHTKExtService) GetPriceFromDelivery(ctx context.Context, shop *domain.Shop, service string) (*domain.Price, *common.Error) {
	token, fromCache, err := g.getToken(ctx, shop, true)
	if err != nil {
		return nil, err
	}

	result, ierr := g.getPriceFromDelivery(ctx, shop, service, token)
	if ierr != nil {
		if fromCache && helpers.IsUnauthorizedError(err) {
			// retry once
			newToken, _, err := g.getToken(ctx, shop, false)
			if err != nil {
				return nil, err
			}
			return g.getPriceFromDelivery(ctx, shop, service, newToken)
		} else {
			return nil, ierr
		}
	}
	return result, nil
}
func (g *GHTKExtService) getPriceFromDelivery(ctx context.Context, shop *domain.Shop, service string, token string) (*domain.Price, *common.Error) {

	return nil, nil
}
func (g *GHTKExtService) getToken(ctx context.Context, shop *domain.Shop, allowFromCache bool) (string, bool, *common.Error) {
	if allowFromCache {
		token, err := g.GetTokenFromCache(ctx, deliveryCode, shop)
		if err == nil && token != "" {
			return token, true, nil
		}
		if err != nil && err != redis.Nil {
			log.Warn(ctx, "Get GHTK Token of shop %s, error: %s", shop.Code, err.Error())
		}
	}
	newToken, err := g.newToken(ctx, shop)
	if err != nil {
		return "", false, err
	}
	go func() {
		g.StoreToken(common.Detach(ctx), deliveryCode, shop, newToken, expireToken)
	}()
	return newToken, false, nil
}

func (g *GHTKExtService) newToken(ctx context.Context, shop *domain.Shop) (string, *common.Error) {
	var output LoginOutput
	resp, err := g.client.R().SetContext(ctx).
		SetHeader("Token", g.token).
		SetBody(newLoginInput(shop)).
		SetSuccessResult(&output).
		SetErrorResult(&output).
		Post(loginPath)
	if err != nil {
		return "", common.ErrSystemError(ctx, err.Error())
	}

	if resp.IsErrorState() || !g.success(output.Code) {
		log.Debug(ctx, "Call ghtk failed with body: %+v", output)
		detail := fmt.Sprintf("http: [%d], resp: [%s]", resp.StatusCode, resp.String())
		return "", common.ErrSystemError(ctx, err.Error()).SetDetail(detail).SetSource(common.SourceGHTKService)
	}
	return output.Data.Code, nil
}

func (g *GHTKExtService) success(code string) bool {
	return code == codeSuccess
}

func (g *GHTKExtService) handleError(ctx context.Context, code string) *common.Error {
	switch code {
	case "ERROR_INVALID_ACCOUNT":
		return common.ErrBadRequest(ctx)
	default:
		detail := fmt.Sprintf("call ghtk error:[%s]", code)
		return common.ErrSystemError(ctx, detail)
	}
}
