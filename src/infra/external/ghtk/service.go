package ghtkext

import (
	"check-price/src/common"
	"check-price/src/common/configs"
	"check-price/src/common/log"
	"check-price/src/core/domain"
	"check-price/src/core/dto"
	"check-price/src/helpers"
	"check-price/src/infra/external"
	"context"
	"fmt"
	"github.com/imroc/req/v3"
	"github.com/redis/go-redis/v9"
	"time"
)

const (
	deliveryCode = "GHTK"
	codeSuccess  = "SUCCESS"
	timeoutGHTK  = 5 * time.Second
	expireToken  = 24 * 30 * time.Hour

	loginPath           = "/services/shops/token"
	getPriceUnder20Path = "/services/shipment/fee"
	getPriceOver20Path  = "/services/shipment/3pl/fee"
)

type Service struct {
	*external.BaseClient
	client *req.Client
	token  string
}

func NewService(base *external.BaseClient) *Service {
	cf := configs.Get().ExtService.GHTK
	cli := req.C().SetBaseURL(cf.Host).SetTimeout(timeoutGHTK)
	cli.SetCommonHeaders(map[string]string{
		"Content-Type": "application/json",
	})
	base.SetTracer(cli)
	return &Service{
		BaseClient: base,
		client:     cli,
		token:      cf.Token,
	}
}

func (g *Service) GetPriceUnder20(ctx context.Context, shop *domain.Shop, serviceCode string, getPriceParam *dto.GetPriceInputDto) (*domain.Price, *common.Error) {
	token, fromCache, err := g.getToken(ctx, shop, true)
	if err != nil {
		return nil, err
	}

	result, ierr := g.getPriceUnder20(ctx, serviceCode, getPriceParam, token)
	if ierr != nil {
		if fromCache && helpers.IsUnauthorizedError(ierr) {
			// retry once
			newToken, _, err := g.getToken(ctx, shop, false)
			if err != nil {
				return nil, err
			}
			return g.getPriceUnder20(ctx, serviceCode, getPriceParam, newToken)
		} else {
			return nil, ierr
		}
	}
	return result, nil
}

func (g *Service) getPriceUnder20(ctx context.Context, serviceCode string, getPriceParam *dto.GetPriceInputDto, token string) (*domain.Price, *common.Error) {
	var output GetPriceUnder20Output
	resp, err := g.api(ctx, token).
		SetBody(newGetPriceUnder20Input(serviceCode, getPriceParam)).
		SetSuccessResult(&output).
		SetErrorResult(&output).
		Get(getPriceUnder20Path)
	if err != nil {
		return nil, common.ErrSystemError(ctx, err.Error()).SetSource(common.SourceGHTKService)
	}

	if !output.Success || resp.IsErrorState() {
		log.Debug(ctx, "Call GetCompany MISA failed with body: %+v", output)
		//Todo consider error code
		//detail := fmt.Sprintf("http: [%d], resp: [%s]", resp.StatusCode, resp.String())
		//return nil, g.handleError(ctx, resp.StatusCode, &output).SetSource(common.SourceGHTKService).SetDetail(detail)
	}
	return output.ToDomainPrice(), nil
}

func (g *Service) GetPriceOver20(ctx context.Context, shop *domain.Shop, serviceCode string, getPriceParam *dto.GetPriceInputDto) (*domain.Price, *common.Error) {
	token, fromCache, err := g.getToken(ctx, shop, true)
	if err != nil {
		return nil, err
	}

	result, ierr := g.getPriceOver20(ctx, serviceCode, getPriceParam, token)
	if ierr != nil {
		if fromCache && helpers.IsUnauthorizedError(ierr) {
			// retry once
			newToken, _, err := g.getToken(ctx, shop, false)
			if err != nil {
				return nil, err
			}
			return g.getPriceOver20(ctx, serviceCode, getPriceParam, newToken)
		} else {
			return nil, ierr
		}
	}
	return result, nil
}

func (g *Service) getPriceOver20(ctx context.Context, serviceCode string, getPriceParam *dto.GetPriceInputDto, token string) (*domain.Price, *common.Error) {
	var output GetPriceOver20Output
	resp, err := g.api(ctx, token).
		//Todo fix this
		SetBody(newGetPriceOver20Input(serviceCode, getPriceParam)).
		SetSuccessResult(&output).
		SetErrorResult(&output).
		Post(getPriceOver20Path)
	if err != nil {
		return nil, common.ErrSystemError(ctx, err.Error()).SetSource(common.SourceGHTKService)
	}

	if !output.Success || resp.IsErrorState() {
		log.Debug(ctx, "Call GetCompany MISA failed with body: %+v", output)
		//Todo consider error code
		//detail := fmt.Sprintf("http: [%d], resp: [%s]", resp.StatusCode, resp.String())
		//return nil, g.handleError(ctx, resp.StatusCode, &output).SetSource(common.SourceGHTKService).SetDetail(detail)
	}
	return output.ToDomain(), nil

}

func (g *Service) api(ctx context.Context, token string) *req.Request {
	return g.client.R().SetContext(ctx).
		SetHeader("token", token)
}

func (g *Service) getToken(ctx context.Context, shop *domain.Shop, allowFromCache bool) (string, bool, *common.Error) {
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
	go g.StoreToken(common.Detach(ctx), deliveryCode, shop, newToken, expireToken)
	return newToken, false, nil
}

func (g *Service) newToken(ctx context.Context, shop *domain.Shop) (string, *common.Error) {
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
		return "", common.ErrSystemError(ctx, detail).SetSource(common.SourceGHTKService)
	}
	return output.Data.Token, nil
}

func (g *Service) success(code string) bool {
	return code == codeSuccess
}

func (g *Service) handleError(ctx context.Context, code string) *common.Error {
	switch code {
	case "ERROR_INVALID_ACCOUNT":
		return common.ErrBadRequest(ctx)
	default:
		detail := fmt.Sprintf("call ghtk error:[%s]", code)
		return common.ErrSystemError(ctx, detail)
	}
}
