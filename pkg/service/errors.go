package service

import (
	"net/http"
)

const (
	ServerERROR    = 1000 // 系统错误
	NotFOUND       = 1001 // 401错误 ​
	UnknownERROR   = 1002 // 未知错误
	ParameterERROR = 1003 // 参数错误​
	AuthERROR      = 1004 // 验证错误
)

type APIException struct {
	Code      int    `json:"-"`
	ErrorCode int    `json:"err_code"`
	Message   string `json:"message"`
}

// 实现接口 Error
//func (e *APIException) Error() string {
//	return e.Message
//}

func newAPIException(code int, errorCode int, message string) *APIException {
	return &APIException{
		Code:      code,
		ErrorCode: errorCode,
		Message:   message,
	}
}

// 参数错误
func ParameterError(message string) *APIException {
	return newAPIException(http.StatusBadRequest, ParameterERROR, message)
}
