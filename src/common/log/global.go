package log

import (
	"check-price/src/common"
	"context"
	"fmt"
)

var globalLogger *logger

func Info(ctx context.Context, msg string, args ...interface{}) {
	globalLogger.Info(addCtxValue(ctx, msg), args...)
}

func Debug(ctx context.Context, msg string, args ...interface{}) {
	globalLogger.Debug(addCtxValue(ctx, msg), args...)
}

func Warn(ctx context.Context, msg string, args ...interface{}) {
	globalLogger.Warn(addCtxValue(ctx, msg), args...)
}

func Error(ctx context.Context, msg string, args ...interface{}) {
	globalLogger.Error(addCtxValue(ctx, msg), args...)
}

func Fatal(msg string, args ...interface{}) {
	globalLogger.Fatal(msg, args...)
}

func IErr(ctx context.Context, err *common.Error) {
	if common.IsInternalError(err) {
		globalLogger.Error(addCtxValue(ctx, err.GetDetail()))
	} else if common.IsClientError(err) {
		globalLogger.Warn(addCtxValue(ctx, err.ToJSon()))
	}

}

func GetLogger() *logger {
	return globalLogger
}

func addCtxValue(ctx context.Context, msg string) string {
	traceId := common.GetTraceId(ctx)
	merchant := common.GetMerchant(ctx)
	return fmt.Sprintf("%s, trace_id:[%s], merchant:[%s]", msg, traceId, merchant)
}
