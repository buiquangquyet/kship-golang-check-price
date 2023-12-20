package helpers

import (
	"check-price/src/common"
	"net/http"
)

func IsClientError(err *common.Error) bool {
	if err == nil {
		return false
	}
	if http.StatusBadRequest <= err.GetHttpStatus() && err.GetHttpStatus() < http.StatusInternalServerError {
		return true
	}
	return false
}

func IsInternalError(err *common.Error) bool {
	if err == nil {
		return false
	}
	if err.GetHttpStatus() >= http.StatusInternalServerError {
		return true
	}
	return false
}

func IsUnauthorizedError(err *common.Error) bool {
	if err == nil {
		return false
	}
	if err.GetHttpStatus() == http.StatusUnauthorized {
		return true
	}
	return false
}

func IsNotFoundError(err *common.Error) bool {
	if err == nil {
		return false
	}
	if err.GetHttpStatus() == http.StatusNotFound {
		return true
	}
	return false
}
