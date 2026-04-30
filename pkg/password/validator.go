package password

import "weclawbotnotify/pkg/xerr"

// ValidateStrength 验证密码强度
// 要求：
//   - 长度至少8位
//   - 包含大写字母
//   - 包含小写字母
//   - 包含数字
func ValidateStrength(password string) error {
	if len(password) < 8 {
		return xerr.PasswordTooWeak
	}

	var (
		hasUpper bool
		hasLower bool
		hasDigit bool
	)

	for _, c := range password {
		switch {
		case c >= 'A' && c <= 'Z':
			hasUpper = true
		case c >= 'a' && c <= 'z':
			hasLower = true
		case c >= '0' && c <= '9':
			hasDigit = true
		}
	}

	if !hasUpper || !hasLower || !hasDigit {
		return xerr.PasswordTooWeak
	}

	return nil
}
