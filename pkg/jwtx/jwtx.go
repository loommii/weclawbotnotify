package jwtx

import "github.com/golang-jwt/jwt/v5"

/*
这里存放的该项目自己用到的一些内容
*/

const (
	Access  = "access"
	Refresh = "refresh"
)

type JWTClaims struct {
	UID       string `json:"uid"`
	TokenType string `json:"token_type"` // 用于区分是刷新Token还是授权Token // access 和 refresh
	// 上方是自己自定义的元素
	jwt.RegisteredClaims // 嵌入标准 claims
}

func (j JWTClaims) IsAccessToken() bool {
	return j.TokenType == Access
}
func (j JWTClaims) IsRefreshToken() bool {
	return j.TokenType == Refresh
}
