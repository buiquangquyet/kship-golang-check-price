package dto

import (
	"check-price/src/common"
	"check-price/src/core/constant"
	"net/http"
)

type ErrorResponse struct {
	Code       int    `json:"status_code"`
	Message    string `json:"message"`
	Msg        string `json:"msg"`
	TraceID    string `json:"trace_id,omitempty"`
	Detail     string `json:"detail,omitempty"`
	Source     string `json:"source"`
	HTTPStatus int    `json:"http_status"`
}

func ConvertErrorToResponse(err *common.Error) *ErrorResponse {
	detail := ""
	if !isInternalError(err) || !constant.IsProdEnv() {
		detail = err.Detail
	}
	return &ErrorResponse{
		Code:       int(err.Code),
		Message:    err.Message,
		TraceID:    err.TraceID,
		Detail:     detail,
		Source:     string(err.Source),
		Msg:        string(err.Message),
		HTTPStatus: err.HTTPStatus,
	}
}

func isInternalError(err *common.Error) bool {
	return err.GetHttpStatus() >= http.StatusInternalServerError
}

type SuccessResponse struct {
	StatusCode int         `json:"status_code"`
	Data       interface{} `json:"data"`
	Msg        string      `json:"msg"`
	Message    string      `json:"message"`
}

func NewSuccessResponse(data interface{}) *SuccessResponse {
	return &SuccessResponse{
		StatusCode: 200,
		Data:       data,
		Msg:        "Thành công",
		Message:    "Thành công",
	}
}
