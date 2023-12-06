package middlewares

import (
	"bytes"
	"check-price/src/common/log"
	"github.com/gin-gonic/gin"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func Log() gin.HandlerFunc {
	return func(c *gin.Context) {
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		c.Next()
		log.Info(c, "response: [%s], path: [%v], status: [%v]", blw.body.String(), c.Request.URL.Path, c.Writer.Status())
	}
}
