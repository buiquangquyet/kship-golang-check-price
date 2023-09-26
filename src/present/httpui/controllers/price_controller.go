package controllers

import (
	"check-price/src/core/service"
	"check-price/src/present/httpui/request"
	"github.com/gin-gonic/gin"
)

type PriceController struct {
	*BaseController
	priceService *service.PriceService
}

func NewPriceController(
	baseController *BaseController,
	priceService *service.PriceService,
) *PriceController {
	return &PriceController{
		BaseController: baseController,
		priceService:   priceService,
	}
}

func (m *PriceController) GetPrice(c *gin.Context) {
	req := new(request.GetPriceReRequest)
	if ierr := m.BindAndValidateRequest(c, req); ierr != nil {
		m.ErrorData(c, ierr)
		return
	}
	tokenInfo, ierr := m.getTokenInfo(c)
	if ierr != nil {
		m.ErrorData(c, ierr)
		return
	}
	data, err := m.priceService.GetPrice(c, req, tokenInfo)
	if err != nil {
		m.ErrorData(c, err)
		return
	}
	m.Success(c, data)
}
