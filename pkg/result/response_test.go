package result

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"weclawbotnotify/pkg/xerr"
)

func TestOkHandler_NilData(t *testing.T) {
	result := OkHandler(context.Background(), nil)
	resp, ok := result.(*Response)
	if !ok {
		t.Fatal("OkHandler should return *Response")
	}
	if resp.Code != 0 {
		t.Errorf("Code = %v, want 0", resp.Code)
	}
	if resp.Msg != "请求成功" {
		t.Errorf("Msg = %v, want 请求成功", resp.Msg)
	}
	if resp.Data == nil {
		t.Error("Data should not be nil, expected empty map")
	}
}

func TestOkHandler_WithData(t *testing.T) {
	data := map[string]string{"key": "value"}
	result := OkHandler(context.Background(), data)
	resp, ok := result.(*Response)
	if !ok {
		t.Fatal("OkHandler should return *Response")
	}
	if resp.Code != 0 {
		t.Errorf("Code = %v, want 0", resp.Code)
	}
	m, ok := resp.Data.(map[string]string)
	if !ok {
		t.Fatal("Data type mismatch")
	}
	if m["key"] != "value" {
		t.Errorf("Data[key] = %v, want value", m["key"])
	}
}

func TestOkHandler_ResponsePassthrough(t *testing.T) {
	original := &Response{Code: 42, Msg: "custom", Data: "hello"}
	result := OkHandler(context.Background(), original)
	resp, ok := result.(*Response)
	if !ok {
		t.Fatal("OkHandler should return *Response")
	}
	if resp != original {
		t.Error("OkHandler should pass through *Response directly")
	}
}

func TestOkHandler_NilPointer(t *testing.T) {
	var s *string
	result := OkHandler(context.Background(), s)
	resp, ok := result.(*Response)
	if !ok {
		t.Fatal("OkHandler should return *Response")
	}
	if resp.Data == nil {
		t.Error("nil pointer should be converted to empty map")
	}
}

func TestErrorHandler_CodeError(t *testing.T) {
	err := xerr.NewCodeError(100000, "请求参数错误")
	httpCode, result := ErrorHandler(context.Background(), err)
	resp, ok := result.(*Response)
	if !ok {
		t.Fatal("ErrorHandler should return *Response")
	}
	if httpCode != http.StatusOK {
		t.Errorf("httpCode = %v, want %v", httpCode, http.StatusOK)
	}
	if resp.Code != 100000 {
		t.Errorf("Code = %v, want 100000", resp.Code)
	}
	if resp.Msg != "请求参数错误" {
		t.Errorf("Msg = %v, want 请求参数错误", resp.Msg)
	}
}

func TestErrorHandler_JwtError(t *testing.T) {
	httpCode, result := ErrorHandler(context.Background(), xerr.JwtError)
	resp, ok := result.(*Response)
	if !ok {
		t.Fatal("ErrorHandler should return *Response")
	}
	if httpCode != http.StatusUnauthorized {
		t.Errorf("httpCode = %v, want %v", httpCode, http.StatusUnauthorized)
	}
	if resp.Code != 100001 {
		t.Errorf("Code = %v, want 100001", resp.Code)
	}
}

func TestErrorHandler_UnknownError(t *testing.T) {
	err := errors.New("something unexpected")
	httpCode, result := ErrorHandler(context.Background(), err)
	resp, ok := result.(*Response)
	if !ok {
		t.Fatal("ErrorHandler should return *Response")
	}
	if httpCode != http.StatusOK {
		t.Errorf("httpCode = %v, want %v", httpCode, http.StatusOK)
	}
	if resp.Code != -1 {
		t.Errorf("Code = %v, want -1", resp.Code)
	}
	if resp.Msg != "something unexpected" {
		t.Errorf("Msg = %v, want something unexpected", resp.Msg)
	}
}

func TestIsNil(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected bool
	}{
		{"nil interface", nil, true},
		{"nil pointer", (*string)(nil), true},
		{"nil slice", ([]int)(nil), true},
		{"nil map", (map[string]int)(nil), true},
		{"nil chan", (chan int)(nil), true},
		{"nil func", (func())(nil), true},
		{"empty slice", []int{}, false},
		{"empty map", map[string]int{}, false},
		{"string", "hello", false},
		{"int", 42, false},
		{"bool", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNil(tt.input); got != tt.expected {
				t.Errorf("IsNil(%v) = %v, want %v", tt.name, got, tt.expected)
			}
		})
	}
}
