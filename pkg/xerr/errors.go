package xerr

import "fmt"

// CodeError 自定义错误结构体，用于业务错误
type CodeError struct {
	HttpCode int    `json:"http_code"`
	Code     int    `json:"code"`
	Message  string `json:"message"`
}

// Error 实现 error 接口
func (e *CodeError) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s", e.Code, e.Message)
}

// NewCodeError 创建一个新的 CodeError
func NewCodeError(code int, msg string) error {
	return &CodeError{Code: code, Message: msg}
}

// NewErrMsg 创建一个只含消息的 CodeError（使用默认错误码）
func NewErrMsg(msg string) error {
	return &CodeError{Code: -1, Message: msg} // -1 作为通用业务错误码
}
