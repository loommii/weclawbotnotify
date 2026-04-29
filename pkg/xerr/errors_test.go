package xerr

import (
	"testing"
)

func TestCodeError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *CodeError
		expected string
	}{
		{
			name:     "basic error",
			err:      &CodeError{Code: 100001, Message: "Unauthorized"},
			expected: "Code: 100001, Message: Unauthorized",
		},
		{
			name:     "success code",
			err:      &CodeError{Code: 0, Message: "请求成功"},
			expected: "Code: 0, Message: 请求成功",
		},
		{
			name:     "with http code",
			err:      &CodeError{HttpCode: 400, Code: 100000, Message: "请求参数错误"},
			expected: "Code: 100000, Message: 请求参数错误",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("CodeError.Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestNewCodeError(t *testing.T) {
	err := NewCodeError(100001, "Unauthorized")
	ce, ok := err.(*CodeError)
	if !ok {
		t.Fatal("NewCodeError should return *CodeError")
	}
	if ce.Code != 100001 {
		t.Errorf("Code = %v, want 100001", ce.Code)
	}
	if ce.Message != "Unauthorized" {
		t.Errorf("Message = %v, want Unauthorized", ce.Message)
	}
}

func TestNewErrMsg(t *testing.T) {
	err := NewErrMsg("something went wrong")
	ce, ok := err.(*CodeError)
	if !ok {
		t.Fatal("NewErrMsg should return *CodeError")
	}
	if ce.Code != -1 {
		t.Errorf("Code = %v, want -1", ce.Code)
	}
	if ce.Message != "something went wrong" {
		t.Errorf("Message = %v, want something went wrong", ce.Message)
	}
}

func TestPredefinedErrors(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		code     int
		message  string
	}{
		{"success", Success, 0, "请求成功"},
		{"request param error", RequestParamError, 100000, "请求参数错误"},
		{"jwt error", JwtError, 100001, "Unauthorized"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ce, ok := tt.err.(*CodeError)
			if !ok {
				t.Fatalf("%s should be *CodeError", tt.name)
			}
			if ce.Code != tt.code {
				t.Errorf("Code = %v, want %v", ce.Code, tt.code)
			}
			if ce.Message != tt.message {
				t.Errorf("Message = %v, want %v", ce.Message, tt.message)
			}
		})
	}
}
