package external

import (
	"check-price/src/common/configs"
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
	tracer trace.Tracer
	cache  redis.UniversalClient
}

func NewBaseClient(cache redis.UniversalClient) *baseClient {
	return &baseClient{
		tracer: otel.Tracer(configs.Get().Server.Name),
		cache:  cache,
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
