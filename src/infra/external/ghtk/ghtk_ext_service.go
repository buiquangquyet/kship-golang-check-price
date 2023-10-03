package ghtk

import (
	"check-price/src/common"
	"check-price/src/common/configs"
	"check-price/src/common/log"
	"check-price/src/core/domain"
	"check-price/src/infra/external"
	"context"
	"fmt"
	"github.com/imroc/req/v3"
	"time"
)

const (
	deliveryCode = "GHTK"
	codeSuccess  = "SUCCESS"
	timeoutGHTK  = 5 * time.Second

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
	return &domain.Price{}, nil
}

//func (z *GHTKExtService) getToken(ctx context.Context, shop *domain.Shop, allowFromCache bool) (string, bool, *common.Error) {
//	if allowFromCache {
//		token, err := z.getTokenFromCache(ctx, merchant)
//		if err == nil && token != "" {
//			return token, true, nil
//		}
//		if err != nil && err != redis.Nil {
//			log.Warn(ctx, "Get ZALO Token of merchant %d %s, error: %s", merchant.MerchantId, merchant.MerchantCode.ToString(), err.Error())
//		}
//	}
//	refreshToken, err := z.refreshTokenRepo.Get(ctx, merchant)
//	if err != nil {
//		if helpers.IsNotFoundError(err) {
//			return "", false, common.ErrUnauthorized(ctx).SetCode(common.ErrorZALOUnauthorized).SetDetail(err.Detail)
//		}
//		log.Error(ctx, "get token")
//		return "", false, err
//	}
//	if time.Now().After(refreshToken.ExpireIn) {
//		return "", false, common.ErrBadRequest(ctx)
//	}
//	newToken, err := z.newToken(ctx, refreshToken.Token)
//	if err != nil {
//		return "", false, err
//	}
//
//	//Todo go routine
//	z.storeToken(common.Detach(ctx), merchant, refreshToken.Id, newToken)
//	return newToken.AccessToken, false, nil
//}

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
