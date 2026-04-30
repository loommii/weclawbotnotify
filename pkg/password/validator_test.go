package password

import (
	"testing"

	"weclawbotnotify/pkg/xerr"
)

// TestValidateStrength_WeakPasswords 测试弱密码
func TestValidateStrength_WeakPasswords(t *testing.T) {
	weakPasswords := []string{
		"",
		"short",
		"nouppercase1",
		"NOLOWERCASE1",
		"NoDigitsHere",
		"12345678",
		"abcdefgh",
		"ABCDEFGH",
	}

	for _, pwd := range weakPasswords {
		err := ValidateStrength(pwd)
		if err == nil {
			t.Errorf("弱密码 '%s' 应返回错误", pwd)
		}

		codeErr, ok := err.(*xerr.CodeError)
		if !ok {
			t.Fatalf("预期 *xerr.CodeError, got %T", err)
		}
		if codeErr.Code != 110400 {
			t.Errorf("弱密码 '%s' 预期错误码 110400, got %d", pwd, codeErr.Code)
		}
	}
}

// TestValidateStrength_StrongPasswords 测试强密码
func TestValidateStrength_StrongPasswords(t *testing.T) {
	strongPasswords := []string{
		"Password1",
		"Strong123!",
		"Complex@Pass99",
		"Abcdefg1",
		"Test1234",
		"MinLength1",
	}

	for _, pwd := range strongPasswords {
		err := ValidateStrength(pwd)
		if err != nil {
			t.Errorf("强密码 '%s' 不应返回错误, got %v", pwd, err)
		}
	}
}
