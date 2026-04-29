package result

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"weclawbotnotify/pkg/xerr"
)

// Response 是返回给前端的统一结构体
type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data,omitempty"`
}

func OkHandler(ctx context.Context, v any) any {
	// 如果 logic 返回的就是一个 *Response 结构，直接返回，避免重复包装
	if resp, ok := v.(*Response); ok {
		return resp
	}

	data := v
	// 如果 logic 返回的数据是 nil，我们给一个默认的空对象，以保证前端 JSON 结构的统一性
	if IsNil(data) {
		data = make(map[string]any)
	}
	return &Response{
		Code: 0,
		Msg:  "请求成功",
		Data: data,
	}
}
func ErrorHandler(ctx context.Context, err error) (int, any) {
	httpCode := http.StatusOK
	response := &Response{}
	switch e := err.(type) {
	// 业务类型的错误
	case *xerr.CodeError:
		if errors.Is(e, xerr.JwtError) { // 特殊错误
			httpCode = http.StatusUnauthorized
		} else {
			httpCode = http.StatusOK
		}
		response = &Response{
			Code: e.Code,
			Msg:  e.Message,
		}
	default:
		// 捕捉的未知错误
		httpCode = http.StatusOK
		response = &Response{
			Code: -1,
			Msg:  err.Error(),
		}
	}
	return httpCode, response
}
func IsNil(i any) bool {
	// 1. 先处理真正的 nil 接口，这是最快、最安全的第一步。
	if i == nil {
		return true
	}
	// 2. 获取 reflect.Value
	val := reflect.ValueOf(i)
	// 3. 【关键的保护步骤】使用 switch 检查 Kind
	//    只在类型是“可为空”的几种类型时，才去调用 IsNil()
	switch val.Kind() {
	case reflect.Ptr, reflect.Map, reflect.Slice, reflect.Chan, reflect.Func, reflect.Interface:
		// 对于这些安全的类型，我们才调用 val.IsNil()
		return val.IsNil()
	default:
		// 4. 对于 int, string, struct 等其他所有类型，
		//    它们本身不可能是 nil，所以直接返回 false。
		return false
	}
}
