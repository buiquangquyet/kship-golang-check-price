package controllers

import (
	"check-price/src/common"
	"check-price/src/core/constant"
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
	retailerId, exist := c.Get(constant.MerchantIdKey)
	if !exist {
		ierr := common.ErrSystemError(c, "retailer id not exist ")
		m.ErrorData(c, ierr)
		return
	}
	versionLocation, exist := c.Get(constant.VersionLocation)
	if !exist {
		ierr := common.ErrSystemError(c, "version location not exist ")
		m.ErrorData(c, ierr)
		return
	}
	req.RetailerId = retailerId.(int64)
	req.VersionLocation = versionLocation.(int)
	data, err := m.priceService.GetPrice(c, req)
	if err != nil {
		m.ErrorData(c, err)
		return
	}
	m.Success(c, data)
}
