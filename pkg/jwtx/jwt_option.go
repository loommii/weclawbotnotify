package jwtx

import (
	"crypto/rsa"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Option JWTHelper 函数选项
type Option func(*JWTHelper)

// WithPrivateKey 设置 RSA 私钥
func WithPrivateKey(privateKey *rsa.PrivateKey) Option {
	return func(helper *JWTHelper) {
		helper.privateKey = privateKey
	}
}

// WithPublicKey 设置 RSA 公钥
func WithPublicKey(publicKey *rsa.PublicKey) Option {
	return func(helper *JWTHelper) {
		helper.publicKey = publicKey
	}
}

// WithExpiredTime 设置 Token 过期时间
func WithExpiredTime(d time.Duration) Option {
	return func(helper *JWTHelper) {
		helper.expire = d
	}
}

// WithSigningMethod 设置签名算法
func WithSigningMethod(method jwt.SigningMethod) Option {
	return func(helper *JWTHelper) {
		helper.signingMethod = method
	}
}
