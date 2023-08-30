package external

import (
	"check-price/src/common/configs"
	"check-price/src/common/log"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/base64"
	"github.com/go-redis/redis/v8"
	"github.com/imroc/req/v3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const (
	keyTokenFormat          = "%s.%s.%s.token"
	keyZOALoginStatusFormat = "%s.%s.%s.zoa_status"
)

type baseClient struct {
	tracer        trace.Tracer
	cache         redis.UniversalClient
	rsaPrivateKey *rsa.PrivateKey
}

func NewBaseClient(cache redis.UniversalClient, rsaPrivateKey *rsa.PrivateKey) *baseClient {
	return &baseClient{
		tracer:        otel.Tracer(configs.Get().Server.Name),
		cache:         cache,
		rsaPrivateKey: rsaPrivateKey,
	}
}

func (b *baseClient) SetTracer(c *req.Client) {
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

func (b *baseClient) SignBody(ctx context.Context, body []byte) (string, error) {
	hashed := sha1.Sum(body)
	sign, err := rsa.SignPKCS1v15(rand.Reader, b.rsaPrivateKey, crypto.SHA1, hashed[:])
	if err != nil {
		log.Error(ctx, "sign body error, body:[%s]", string(body))
		return "", err
	}
	return base64.StdEncoding.EncodeToString(sign), nil
}
