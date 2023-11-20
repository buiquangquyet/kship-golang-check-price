package aieliminating

import (
	"check-price/src/common"
	"check-price/src/common/configs"
	"check-price/src/common/log"
	"check-price/src/helpers"
	"check-price/src/infra/external"
	"context"
	"fmt"
	"github.com/imroc/req/v3"
	"github.com/redis/go-redis/v9"
	"time"
)

const (
	timeout     = 5 * time.Second
	expireToken = 24 * time.Hour
	keyCache    = "ai-address-eliminating-redundancy-token"

	loginPath  = "/address_eliminating_redundancy/login"
	redundancy = "/address_eliminating_redundancy"
)

type Service struct {
	*external.BaseClient
	client *req.Client
	cf     *configs.AIEliminating
}

func NewService(base *external.BaseClient) *Service {
	cf := configs.Get().ExtService.AIEliminating
	cli := req.C().SetBaseURL(cf.Host).SetTimeout(timeout)
	cli.SetCommonHeaders(map[string]string{
		"Content-Type": "application/json",
	})
	base.SetTracer(cli)
	return &Service{
		BaseClient: base,
		client:     cli,
		cf:         cf,
	}
}

func (g *Service) api(ctx context.Context) *req.Request {
	return g.client.R().SetContext(ctx)
}

func (g *Service) Redundancy(ctx context.Context, address, ward, district, province string) (string, *common.Error) {
	token, fromCache, err := g.getToken(ctx, true)
	if err != nil {
		return "", err
	}
	result, ierr := g.redundancy(ctx, token, address, ward, district, province)
	if ierr != nil {
		if fromCache && helpers.IsUnauthorizedError(ierr) {
			// retry once
			newToken, _, err := g.getToken(ctx, false)
			if err != nil {
				return "", err
			}
			return g.redundancy(ctx, newToken, address, ward, district, province)
		}
		return "", ierr
	}
	return result, nil
}

func (g *Service) redundancy(ctx context.Context, token string, address, ward, district, province string) (string, *common.Error) {
	var output RedundancyOutput
	var outputErr OutputError
	resp, err := g.api(ctx).
		SetBearerAuthToken(token).
		SetFormData(map[string]string{
			"address":       address,
			"ward_name":     ward,
			"district_name": district,
			"province_name": province,
		}).
		SetSuccessResult(&output).
		SetErrorResult(&outputErr).
		Post(redundancy)
	if err != nil {
		return "", common.ErrSystemError(ctx, err.Error())
	}

	if resp.IsErrorState() {
		log.Debug(ctx, "Call AI Eliminating failed with body: %+v", output)
		detail := fmt.Sprintf("http: [%d], resp: [%s]", resp.StatusCode, resp.String())
		return "", common.ErrSystemError(ctx, detail).SetSource(common.SourceGHTKService)
	}
	return output.Data.AddressNew, nil
}

func (g *Service) getToken(ctx context.Context, allowFromCache bool) (string, bool, *common.Error) {
	if allowFromCache {
		token, err := g.GetDataCache(ctx, keyCache)
		if err == nil && token != "" {
			return token, true, nil
		}
		if err != nil && err != redis.Nil {
			log.Warn(ctx, err.Error())
		}
	}
	newToken, err := g.newToken(ctx)
	if err != nil {
		return "", false, err
	}
	go g.StoreData(common.Detach(ctx), keyCache, newToken, expireToken)
	return newToken, false, nil
}

func (g *Service) newToken(ctx context.Context) (string, *common.Error) {
	var output LoginOutput
	var outputErr OutputError
	resp, err := g.api(ctx).
		SetFormData(map[string]string{
			"username": g.cf.Username,
			"password": g.cf.Password,
		}).
		SetSuccessResult(&output).
		SetErrorResult(&outputErr).
		Post(loginPath)
	if err != nil {
		return "", common.ErrSystemError(ctx, err.Error())
	}

	if resp.IsErrorState() {
		log.Debug(ctx, "Call AI Eliminating failed with body: %+v", output)
		detail := fmt.Sprintf("http: [%d], resp: [%s]", resp.StatusCode, resp.String())
		return "", common.ErrSystemError(ctx, detail).SetSource(common.SourceGHTKService)
	}
	return output.Data.AccessToken, nil
}

func (g *Service) handleError(ctx context.Context) *common.Error {
	//Todo code

	return nil
}
