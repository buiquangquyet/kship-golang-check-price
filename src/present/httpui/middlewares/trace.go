package middlewares

import (
	"check-price/src/core/constant"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
	"math/rand"
)

func Tracer() gin.HandlerFunc {
	return func(c *gin.Context) {
		spanContext := trace.SpanContextFromContext(c.Request.Context())
		span := trace.SpanFromContext(c.Request.Context())
		span.SetName(fmt.Sprintf("[%s] %s", c.Request.Method, c.FullPath()))
		var traceId string
		if spanContext.SpanID().IsValid() {
			traceId = spanContext.TraceID().String()
		} else {
			traceIdByte := make([]byte, 16)
			rand.Read(traceIdByte)
			traceId = hex.EncodeToString(traceIdByte[:])
		}
		c.Set(constant.TraceIdName, traceId)
		c.Next()
	}
}
