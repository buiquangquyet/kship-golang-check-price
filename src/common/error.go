package common

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ErrorCode int

const (
	//internal
	Success               ErrorCode = 0
	ErrorCodeBadRequest   ErrorCode = 400
	ErrorCodeNotFound     ErrorCode = 404
	ErrorCodeUnauthorized ErrorCode = 401
	ErrorCodeSystemError  ErrorCode = 500

	ErrorCodeMessageHasSent ErrorCode = 1001
	ErrorCodeMessageIsNew   ErrorCode = 1002

	//external
	ErrorZALOUnauthorized ErrorCode = 1401

	ErrorInvalidParam            ErrorCode = 1004
	ErrorPhoneNumberInvalid      ErrorCode = 1108
	ErrorZaloAccountNotExisted   ErrorCode = 1404
	ErrorTemplateIdInvalid       ErrorCode = 1400
	ErrorOADoesNotHavePermission ErrorCode = 1403
	ErrorDailyQuotaExceeded      ErrorCode = 1144

	ErrorTemplateNotApproved ErrorCode = 1131
	ErrorUserRefused         ErrorCode = 1141
	ErrorOutOfQuota          ErrorCode = 1115
	ErrorBodyDataEmpty       ErrorCode = 1121
	ErrorUnknownErrorZalo    ErrorCode = 1500

	ErrorKWalletNotFound                              ErrorCode = 2000
	ErrorKWalletNotEnoughMoney                        ErrorCode = 2001
	ErrorKWalletTotalRefundExceedsOriginalTransaction ErrorCode = 2002
	ErrorCodeSystemErrorKWallet                       ErrorCode = 2500
)

type Source string

const (
	SourceAPIService     Source = "API_Service"
	SourceZALOService    Source = "ZALO_service"
	SourceInfraService   Source = "Infra_Service"
	SourceKWalletService Source = "KWallet_Service"
)

type Error struct {
	Code       ErrorCode `json:"code"`
	Message    string    `json:"message"`
	TraceID    string    `json:"trace_id,omitempty"`
	Detail     string    `json:"detail"`
	Source     Source    `json:"source"`
	HTTPStatus int       `json:"http_status"`
}

func NewError(ctx context.Context, code ErrorCode, message string, httpStatus int) *Error {
	traceId := GetTraceId(ctx)
	return &Error{
		Code:       code,
		Message:    message,
		TraceID:    traceId,
		HTTPStatus: httpStatus,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("code:[%d], message:[%s], detail:[%s], source:[%s]", e.Code, e.Message, e.Detail, e.Source)
}

func (e *Error) GetHttpStatus() int {
	return e.HTTPStatus
}

func (e *Error) GetCode() ErrorCode {
	return e.Code
}

func (e *Error) GetMessage() string {
	return e.Message
}

func (e *Error) SetCode(code ErrorCode) *Error {
	e.Code = code
	return e
}

func (e *Error) SetTraceId(traceId string) *Error {
	e.TraceID = fmt.Sprintf("%s:%d", traceId, time.Now().Unix())
	return e
}

func (e *Error) SetHTTPStatus(status int) *Error {
	e.HTTPStatus = status
	return e
}

func (e *Error) SetMessage(msg string) *Error {
	e.Message = msg
	return e
}

func (e *Error) SetDetail(detail string) *Error {
	e.Detail = detail
	return e
}

func (e *Error) GetDetail() string {
	return e.Detail
}

func (e *Error) SetSource(source Source) *Error {
	e.Source = source
	return e
}

func (e *Error) ToJSon() string {
	data, err := json.Marshal(e)
	if err != nil {
		return "marshal error failed"
	}
	return string(data)
}

var (
	// Status 4xx ********

	ErrUnauthorized = func(ctx context.Context) *Error {
		traceId := GetTraceId(ctx)
		return &Error{
			Code:       ErrorCodeUnauthorized,
			Message:    DefaultUnauthorizedMessage,
			TraceID:    traceId,
			Source:     SourceAPIService,
			HTTPStatus: http.StatusUnauthorized,
		}
	}

	ErrBadRequest = func(ctx context.Context) *Error {
		traceId := GetTraceId(ctx)
		return &Error{
			Code:       ErrorCodeBadRequest,
			Message:    DefaultBadRequestMessage,
			TraceID:    traceId,
			HTTPStatus: http.StatusBadRequest,
			Source:     SourceAPIService,
		}
	}

	ErrNotFound = func(ctx context.Context, object, status string) *Error {
		traceId := GetTraceId(ctx)
		return &Error{
			Code:       ErrorCodeNotFound,
			Message:    fmt.Sprintf("%s %s", object, status),
			TraceID:    traceId,
			HTTPStatus: http.StatusNotFound,
		}
	}

	// Status 5xx *******

	ErrSystemError = func(ctx context.Context, detail string) *Error {
		traceId := GetTraceId(ctx)
		return &Error{
			Code:       ErrorCodeSystemError,
			Message:    DefaultServerErrorMessage,
			TraceID:    traceId,
			HTTPStatus: http.StatusInternalServerError,
			Source:     SourceAPIService,
			Detail:     detail,
		}
	}
)

const (
	DefaultServerErrorMessage  = "Something has gone wrong, please contact admin"
	DefaultBadRequestMessage   = "Invalid request"
	DefaultUnauthorizedMessage = "Token invalid"
)

func IsClientError(err *Error) bool {
	if err == nil {
		return false
	}
	if http.StatusBadRequest <= err.GetHttpStatus() && err.GetHttpStatus() < http.StatusInternalServerError {
		return true
	}
	return false
}

func IsInternalError(err *Error) bool {
	if err == nil {
		return false
	}
	if err.GetHttpStatus() >= http.StatusInternalServerError {
		return true
	}
	return false
}
