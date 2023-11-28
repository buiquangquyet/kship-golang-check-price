package external

import (
	"check-price/src/common/configs"
	"check-price/src/common/log"
	"check-price/src/core/domain"
	"context"
	"fmt"
	"github.com/imroc/req/v3"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"strings"
	"time"
)

const (
	keyTokenFormat = "shipping:%s-token-%s"
)

type BaseClient struct {
	tracer trace.Tracer
	cache  redis.UniversalClient
}

func NewBaseClient(cache redis.UniversalClient) *BaseClient {
	return &BaseClient{
		tracer: otel.Tracer(configs.Get().Server.Name),
		cache:  cache,
	}
}

func (b *BaseClient) SetTracer(c *req.Client) {
	c.WrapRoundTripFunc(func(rt req.RoundTripper) req.RoundTripFunc {
		return func(req *req.Request) (resp *req.Response, err error) {
			apiName := req.URL.Path
			_, span := b.tracer.Start(req.Context(), apiName)
			defer span.End()
			span.SetAttributes(
				attribute.String("http.url", req.URL.String()),
				attribute.String("http.method", req.Method),
			)
			if len(req.Body) > 0 {
				span.SetAttributes(
					attribute.String("http.req.body", string(req.Body)),
				)
			}
			resp, err = rt.RoundTrip(req)
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			}
			if resp.Response == nil {
				return resp, nil
			}
			span.SetAttributes(
				attribute.Int("http.status_code", resp.StatusCode),
			)
			if !resp.IsSuccessState() {
				span.SetAttributes(
					attribute.String("http.resp.header", resp.HeaderToString()),
					attribute.String("http.resp.body", resp.String()),
				)
			}
			return
		}
	})
}

func (b *BaseClient) GetTokenFromCache(ctx context.Context, deliveryCode string, shop *domain.Shop) (string, error) {
	keyOfToken := fmt.Sprintf(keyTokenFormat, strings.ToLower(deliveryCode), shop.Code)
	var token string
	if err := b.cache.Get(ctx, keyOfToken).Scan(&token); err != nil {
		return "", err
	}
	return token, nil
}

func (b *BaseClient) GetDataCache(ctx context.Context, key string) (string, error) {
	var data string
	if err := b.cache.Get(ctx, key).Scan(&data); err != nil {
		return "", err
	}
	return data, nil
}

func (b *BaseClient) StoreToken(ctx context.Context, deliveryCode string, shop *domain.Shop, newToken string, expire time.Duration) {
	keyOfToken := fmt.Sprintf(keyTokenFormat, strings.ToLower(deliveryCode), shop.Code)
	if err := b.cache.Set(ctx, keyOfToken, newToken, expire).Err(); err != nil {
		log.Warn(ctx, "Cache Token failed, delivery: [%s], err: [%s]", deliveryCode, err.Error())
	}
	return
}

func (b *BaseClient) StoreData(ctx context.Context, key string, newToken string, expire time.Duration) {
	if err := b.cache.Set(ctx, key, newToken, expire).Err(); err != nil {
		log.Warn(ctx, "Cache Token failed, err: [%s]", err.Error())
	}
	return
}
