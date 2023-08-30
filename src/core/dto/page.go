package dto

import "check-price/src/present/httpui/request"

type PageDto struct {
	Page  int
	Limit int
}

func NewPageDto(req *request.Page) *PageDto {
	return &PageDto{
		Page:  req.Page,
		Limit: req.Limit,
	}
}

func (p *PageDto) GetOffset() int {
	return (p.Page - 1) * p.Limit
}

func (p *PageDto) GetLimit() int {
	return p.Limit
}

type PageResponse struct {
	Total int64       `json:"total"`
	Data  interface{} `json:"data"`
}

func NewPageResponse(total int64, data interface{}) *PageResponse {
	return &PageResponse{
		Total: total,
		Data:  data,
	}
}
