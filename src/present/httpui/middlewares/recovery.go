package middlewares

import (
	"check-price/src/common"
	"check-price/src/common/log"
	"check-price/src/core/dto"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http/httputil"
	"runtime/debug"
	"time"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				log.Error(c, "[Recovery from panic]\ntime: [%v]\nerror: [%v] \nrequest: [%v]\nstack: [%v]\n",
					time.Now(), err, string(httpRequest), string(debug.Stack()))
				e := common.ErrSystemError(c, fmt.Sprintf("recovery, err:[%s]", err))
				c.JSON(e.GetHttpStatus(), dto.ConvertErrorToResponse(e))
				c.Abort()
			}
		}()
		c.Next()
	}
}
