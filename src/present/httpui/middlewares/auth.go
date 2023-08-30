package middlewares

import (
	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		//xApiKey := c.GetHeader("X-API-KEY")
		//
		//for password, merchantCode := range configs.Get().ApiKey {
		//	if password == xApiKey {
		//		c.Set(constant.MerchantCodeKey, merchantCode)
		//		c.Next()
		//		return
		//	}
		//}
		//err := common.ErrUnauthorized(c).SetMessage("Unauthorized").SetSource(common.SourceAPIService).SetDetail("Invalid X-API-KEY")
		//c.JSON(err.GetHttpStatus(), dto.ConvertErrorToResponse(err))
		//c.Abort()
		c.Next()
		return
	}
}
