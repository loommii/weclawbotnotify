package jwtx

import "github.com/golang-jwt/jwt/v5"

const (
	Access  = "access"  // Access Token 类型标识
	Refresh = "refresh" // Refresh Token 类型标识
)

// JWTClaims 双令牌通用 Claims
// TokenType 区分 access/refresh，JTI 用于 Refresh Token 轮换验证
type JWTClaims struct {
	UID       string `json:"uid"`           // 用户 ID
	TokenType string `json:"token_type"`    // 令牌类型：access / refresh
	JTI       string `json:"jti,omitempty"` // JWT ID，Refresh Token 轮换时使用
	jwt.RegisteredClaims
}

// IsAccessToken 判断是否为 Access Token
func (j JWTClaims) IsAccessToken() bool {
	return j.TokenType == Access
}

// IsRefreshToken 判断是否为 Refresh Token
func (j JWTClaims) IsRefreshToken() bool {
	return j.TokenType == Refresh
}
