package dto

import (
	"check-price/src/common"
)

type ErrorResponse struct {
	Code string `json:"status_code"`
	Data struct {
		ErrorCode string `json:"error_code"`
		Msg       string `json:"msg"`
		Message   string `json:"message"`
	}
	Message string `json:"message"`
	Msg     string `json:"msg"`
	TraceID string `json:"trace_id,omitempty"`
}

func ConvertErrorToResponse(err *common.Error) *ErrorResponse {
	msg, exist := mapErrorMsg[int(err.Code)]
	if !exist {
		msg = ""
	}
	code := err.Code.ToString()
	return &ErrorResponse{
		Code: code,
		Data: struct {
			ErrorCode string `json:"error_code"`
			Msg       string `json:"msg"`
			Message   string `json:"message"`
		}{
			ErrorCode: code,
			Msg:       msg,
			Message:   msg,
		},
		Message: msg,
		TraceID: err.TraceID,
		Msg:     "Có lỗi xảy ra",
	}
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

var mapErrorMsg = map[int]string{
	422:  "Token không hợp lệ.",
	3002: "Hãng vận chuyển không được bật",
	3004: "Shop không được phép check giá",
	3012: "Vui lòng đăng nhập tài khoản GHTK",

	4004: "Vui lòng nhập phường/xã người gửi",
	4006: "Vui lòng nhập phường/xã người nhận",
	4009: "Dịch vụ không tồn tại hoặc không được bật cho hãng vận chuyển",
	4405: "Chỉ nhận thu hộ COD tối đa 20,000,000",
	4406: "Chỉ nhận khai giá tối đa 20,000,000",
}
