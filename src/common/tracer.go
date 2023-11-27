package common

import (
	"check-price/src/common/configs"
	"check-price/src/core/constant"
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func InitTracer(lc fx.Lifecycle, logger *zap.SugaredLogger) {
	cf := configs.Get()
	if !cf.Tracer.Enabled {
		return
	}
	res, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(cf.Server.Name),
			semconv.ServiceVersion("v0.1.0"),
			attribute.String("environment", cf.Mode),
		),
	)
	opts := make([]sdktrace.TracerProviderOption, 0)
	opts = append(opts, sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()))

	if cf.Tracer.Jaeger.Active {
		expJaeger, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(cf.Tracer.Jaeger.Endpoint)))
		if err != nil {
			logger.Fatal("connect jaeger error, err:[%v]", err.Error())
		}
		opts = append(opts, sdktrace.WithBatcher(expJaeger))
	}
	tp := sdktrace.NewTracerProvider(opts...)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}))
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			logger.Info("trace provider shutdown")
			return tp.Shutdown(ctx)
		},
	})
	return
}

func GetTraceId(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	traceId := ""
	if ctx.Value(constant.TraceIdName) != nil {
		traceId = ctx.Value(constant.TraceIdName).(string)
	}
	return traceId
}

func GetMerchant(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	merchant := ""
	if ctx.Value(constant.Merchant) != nil {
		merchant = ctx.Value(constant.Merchant).(string)
	}
	return merchant
}
