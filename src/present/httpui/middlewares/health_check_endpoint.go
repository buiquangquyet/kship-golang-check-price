package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func HealthCheckEndpoint(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
