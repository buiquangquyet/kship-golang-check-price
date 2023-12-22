package controllers

import "github.com/gin-gonic/gin"

type LogController struct {
	*BaseController
}

func NewLogController(baseController *BaseController) *LogController {
	return &LogController{
		BaseController: baseController,
	}
}

func (m *LogController) GetTable(c *gin.Context) {
	return
}
