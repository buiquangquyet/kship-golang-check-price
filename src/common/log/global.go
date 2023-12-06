package log

import (
	"check-price/src/common"
	"check-price/src/core/enums"
	"context"
	"fmt"
)

var globalLogger *logger

func Info(ctx context.Context, msg string, args ...interface{}) {
	globalLogger.Info(addCtxValue(ctx, msg, args))
}

func Debug(ctx context.Context, msg string, args ...interface{}) {
	globalLogger.Debug(addCtxValue(ctx, msg, args))
}

func Warn(ctx context.Context, msg string, args ...interface{}) {
	globalLogger.Warn(addCtxValue(ctx, msg, args))
}

func Error(ctx context.Context, msg string, args ...interface{}) {
	globalLogger.Error(addCtxValue(ctx, msg, args))
}

func Fatal(msg string, args ...interface{}) {
	globalLogger.Fatal(msg, args...)
}

func IErr(ctx context.Context, err *common.Error) {
	if common.IsInternalError(err) {
		globalLogger.Error(addCtxValue(ctx, err.ToJSon()))
	} else if common.IsClientError(err) {
		globalLogger.Warn(addCtxValue(ctx, err.ToJSon()))
	}
}

func GetLogger() *logger {
	return globalLogger
}

func addCtxValue(ctx context.Context, msg string, args ...interface{}) (string, enums.TraceId, enums.Merchant) {
	msg = getMessage(msg, args)
	traceId := common.GetTraceId(ctx)
	merchant := common.GetMerchant(ctx)
	return msg, enums.TraceId(traceId), enums.Merchant(merchant)
}

func getMessage(template string, fmtArgs []interface{}) string {
	if len(fmtArgs) == 0 {
		return template
	}

	if template != "" {
		return fmt.Sprintf(template, fmtArgs...)
	}

	if len(fmtArgs) == 1 {
		if str, ok := fmtArgs[0].(string); ok {
			return str
		}
	}
	return fmt.Sprint(fmtArgs...)
}
