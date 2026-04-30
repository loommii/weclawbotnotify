package jwtx

import (
	"crypto/rsa"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// generateRSAKeys 生成测试用 RSA 密钥对
func generateRSAKeys(t *testing.T) (*rsa.PrivateKey, *rsa.PublicKey) {
	t.Helper()
	privateKey, err := rsa.GenerateKey(nil, 2048)
	if err != nil {
		t.Fatalf("生成 RSA 密钥失败: %v", err)
	}
	return privateKey, &privateKey.PublicKey
}

// TestJWTClaims_IsAccessToken 测试 IsAccessToken 判断逻辑
func TestJWTClaims_IsAccessToken(t *testing.T) {
	tests := []struct {
		name      string
		tokenType string
		expected  bool
	}{
		{"access 类型", Access, true},
		{"refresh 类型", Refresh, false},
		{"空类型", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims := JWTClaims{TokenType: tt.tokenType}
			if got := claims.IsAccessToken(); got != tt.expected {
				t.Errorf("IsAccessToken() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestJWTClaims_IsRefreshToken 测试 IsRefreshToken 判断逻辑
func TestJWTClaims_IsRefreshToken(t *testing.T) {
	tests := []struct {
		name      string
		tokenType string
		expected  bool
	}{
		{"refresh 类型", Refresh, true},
		{"access 类型", Access, false},
		{"空类型", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims := JWTClaims{TokenType: tt.tokenType}
			if got := claims.IsRefreshToken(); got != tt.expected {
				t.Errorf("IsRefreshToken() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestWithExpiredTime 测试 Token 过期时间选项
func TestWithExpiredTime(t *testing.T) {
	duration := 15 * time.Minute
	opt := WithExpiredTime(duration)
	helper := &JWTHelper{}
	opt(helper)
	if helper.expire != duration {
		t.Errorf("expire = %v, want %v", helper.expire, duration)
	}
}

// TestWithPrivateKey 测试私钥选项
func TestWithPrivateKey(t *testing.T) {
	privateKey, _ := generateRSAKeys(t)
	opt := WithPrivateKey(privateKey)
	helper := &JWTHelper{}
	opt(helper)
	if helper.privateKey != privateKey {
		t.Error("WithPrivateKey 未正确设置私钥")
	}
}

// TestWithPublicKey 测试公钥选项
func TestWithPublicKey(t *testing.T) {
	_, publicKey := generateRSAKeys(t)
	opt := WithPublicKey(publicKey)
	helper := &JWTHelper{}
	opt(helper)
	if helper.publicKey != publicKey {
		t.Error("WithPublicKey 未正确设置公钥")
	}
}

// TestWithSigningMethod 测试签名算法选项
func TestWithSigningMethod(t *testing.T) {
	method := jwt.SigningMethodRS512
	opt := WithSigningMethod(method)
	helper := &JWTHelper{}
	opt(helper)
	if helper.signingMethod != method {
		t.Error("WithSigningMethod 未正确设置签名算法")
	}
}

// TestNewJWTHelper_Defaults 测试默认构造
func TestNewJWTHelper_Defaults(t *testing.T) {
	helper := NewJWTHelper()
	if helper.signingMethod != jwt.SigningMethodRS256 {
		t.Errorf("默认签名算法 = %v, want RS256", helper.signingMethod)
	}
}

// TestNewJWTHelper_WithOptions 测试带选项构造
func TestNewJWTHelper_WithOptions(t *testing.T) {
	privateKey, publicKey := generateRSAKeys(t)
	expire := 15 * time.Minute

	helper := NewJWTHelper(
		WithPrivateKey(privateKey),
		WithPublicKey(publicKey),
		WithExpiredTime(expire),
	)

	if helper.privateKey != privateKey {
		t.Error("私钥未设置")
	}
	if helper.publicKey != publicKey {
		t.Error("公钥未设置")
	}
	if helper.expire != expire {
		t.Errorf("expire = %v, want %v", helper.expire, expire)
	}
}

// TestGenerateAccessToken 测试 Access Token 生成与验证
func TestGenerateAccessToken(t *testing.T) {
	privateKey, publicKey := generateRSAKeys(t)
	helper := NewJWTHelper(
		WithPrivateKey(privateKey),
		WithPublicKey(publicKey),
		WithExpiredTime(15*time.Minute),
	)

	claims := JWTClaims{UID: "123", TokenType: Access}
	tokenStr, err := helper.GenerateToken(claims)
	if err != nil {
		t.Fatalf("GenerateToken 失败: %v", err)
	}
	if tokenStr == "" {
		t.Fatal("GenerateToken 返回空令牌")
	}

	_, parsedClaims, err := helper.ValidateToken(tokenStr)
	if err != nil {
		t.Fatalf("ValidateToken 失败: %v", err)
	}
	if parsedClaims.UID != "123" {
		t.Errorf("UID = %v, want 123", parsedClaims.UID)
	}
	if !parsedClaims.IsAccessToken() {
		t.Error("应为 Access Token")
	}
	if parsedClaims.JTI != "" {
		t.Error("Access Token 不应包含 JTI")
	}
}

// TestGenerateRefreshToken 测试 Refresh Token 生成与验证（含 JTI）
func TestGenerateRefreshToken(t *testing.T) {
	privateKey, publicKey := generateRSAKeys(t)
	helper := NewJWTHelper(
		WithPrivateKey(privateKey),
		WithPublicKey(publicKey),
		WithExpiredTime(7*24*time.Hour),
	)

	claims := JWTClaims{UID: "456", TokenType: Refresh}
	tokenStr, err := helper.GenerateToken(claims)
	if err != nil {
		t.Fatalf("GenerateToken 失败: %v", err)
	}
	if tokenStr == "" {
		t.Fatal("GenerateToken 返回空令牌")
	}

	_, parsedClaims, err := helper.ValidateToken(tokenStr)
	if err != nil {
		t.Fatalf("ValidateToken 失败: %v", err)
	}
	if parsedClaims.UID != "456" {
		t.Errorf("UID = %v, want 456", parsedClaims.UID)
	}
	if !parsedClaims.IsRefreshToken() {
		t.Error("应为 Refresh Token")
	}
	if parsedClaims.JTI == "" {
		t.Error("Refresh Token 必须包含 JTI")
	}
}

// TestGenerateRefreshToken_JTIUnique 测试每次生成 Refresh Token 的 JTI 唯一
func TestGenerateRefreshToken_JTIUnique(t *testing.T) {
	privateKey, publicKey := generateRSAKeys(t)
	helper := NewJWTHelper(
		WithPrivateKey(privateKey),
		WithPublicKey(publicKey),
		WithExpiredTime(7*24*time.Hour),
	)

	claims := JWTClaims{UID: "789", TokenType: Refresh}
	token1, _ := helper.GenerateToken(claims)
	token2, _ := helper.GenerateToken(claims)

	_, claims1, _ := helper.ValidateToken(token1)
	_, claims2, _ := helper.ValidateToken(token2)

	if claims1.JTI == claims2.JTI {
		t.Error("两次生成的 Refresh Token JTI 不应相同")
	}
}

// TestValidateToken_Expired 测试过期令牌验证
func TestValidateToken_Expired(t *testing.T) {
	privateKey, publicKey := generateRSAKeys(t)
	helper := NewJWTHelper(
		WithPrivateKey(privateKey),
		WithPublicKey(publicKey),
		WithExpiredTime(-10*time.Second),
	)

	claims := JWTClaims{UID: "123", TokenType: Access}
	tokenStr, err := helper.GenerateToken(claims)
	if err != nil {
		t.Fatalf("GenerateToken 失败: %v", err)
	}

	_, _, err = helper.ValidateToken(tokenStr)
	if err == nil {
		t.Error("过期令牌应返回错误")
	}
}

// TestValidateToken_WrongPublicKey 测试错误公钥验签
func TestValidateToken_WrongPublicKey(t *testing.T) {
	privateKey, _ := generateRSAKeys(t)
	_, wrongPublicKey := generateRSAKeys(t)

	signHelper := NewJWTHelper(
		WithPrivateKey(privateKey),
		WithExpiredTime(1*time.Hour),
	)

	claims := JWTClaims{UID: "123", TokenType: Access}
	tokenStr, err := signHelper.GenerateToken(claims)
	if err != nil {
		t.Fatalf("GenerateToken 失败: %v", err)
	}

	validateHelper := NewJWTHelper(WithPublicKey(wrongPublicKey))
	_, _, err = validateHelper.ValidateToken(tokenStr)
	if err == nil {
		t.Error("错误公钥应验签失败")
	}
}

// TestGenerateToken_OverwritesTime 测试签发时覆盖时间字段
func TestGenerateToken_OverwritesTime(t *testing.T) {
	privateKey, publicKey := generateRSAKeys(t)
	helper := NewJWTHelper(
		WithPrivateKey(privateKey),
		WithPublicKey(publicKey),
		WithExpiredTime(1*time.Hour),
	)

	pastTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	claims := JWTClaims{
		UID:              "789",
		TokenType:        Access,
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(pastTime)},
	}

	tokenStr, err := helper.GenerateToken(claims)
	if err != nil {
		t.Fatalf("GenerateToken 失败: %v", err)
	}

	_, parsedClaims, err := helper.ValidateToken(tokenStr)
	if err != nil {
		t.Fatalf("时间应被覆盖，令牌应有效: %v", err)
	}
	if parsedClaims.ExpiresAt.Before(time.Now()) {
		t.Error("ExpiresAt 应被覆盖为未来时间")
	}
}

// TestAccessTokenRejectedAsRefresh 测试 Access Token 不能当 Refresh Token 使用
func TestAccessTokenRejectedAsRefresh(t *testing.T) {
	privateKey, publicKey := generateRSAKeys(t)
	helper := NewJWTHelper(
		WithPrivateKey(privateKey),
		WithPublicKey(publicKey),
		WithExpiredTime(15*time.Minute),
	)

	claims := JWTClaims{UID: "123", TokenType: Access}
	tokenStr, err := helper.GenerateToken(claims)
	if err != nil {
		t.Fatalf("GenerateToken 失败: %v", err)
	}

	_, parsedClaims, err := helper.ValidateToken(tokenStr)
	if err != nil {
		t.Fatalf("ValidateToken 失败: %v", err)
	}
	if parsedClaims.IsRefreshToken() {
		t.Error("Access Token 不应被识别为 Refresh Token")
	}
}
