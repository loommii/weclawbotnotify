package jwtx

import (
	"crypto/rsa"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5" // 使用真正的 JWT 库
	"github.com/golang-jwt/jwt/v5/request"
	"github.com/pkg/errors"
)

// JWTHelper 结构体，封装 JWT 的操作
type JWTHelper struct {
	publicKey     *rsa.PublicKey
	privateKey    *rsa.PrivateKey
	expiredTime   time.Duration     // 过期时间
	signingMethod jwt.SigningMethod // 签名算法
}

// NewJWTHelper 使用选项模式创建 JWTHelper 实例
func NewJWTHelper(opts ...Option) *JWTHelper {
	helper := &JWTHelper{
		signingMethod: jwt.SigningMethodRS256, // 默认使用 RS256 算法
	}

	// 应用选项
	for _, opt := range opts {
		opt(helper)
	}

	return helper
}

// GenerateToken 生成 JWT token
/*
	注意
		1、这函数的入参值拷贝，不会修改原先的claims
		2、该函数会覆盖时间，传入的时间将无效

*/
func (j *JWTHelper) GenerateToken(claims JWTClaims) (string, error) {
	now := time.Now()

	claims.ExpiresAt = jwt.NewNumericDate(now.Add(j.expiredTime)) // 过期时间
	claims.IssuedAt = jwt.NewNumericDate(now)
	claims.NotBefore = jwt.NewNumericDate(now)

	token := jwt.NewWithClaims(j.signingMethod, claims)
	signedToken, err := token.SignedString(j.privateKey)
	if err != nil {
		return "", errors.Wrap(err, "failed to sign token")
	}

	return signedToken, nil
}

// ValidateToken 验证 JWT token
func (j *JWTHelper) ValidateToken(tokenString string) (*jwt.Token, *JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (any, error) {
		// 检查签名方法
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.publicKey, nil
	}, jwt.WithLeeway(5*time.Second))

	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to parse token")
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return token, claims, nil
	}

	return nil, nil, errors.New("invalid token")
}

// ValidateRequest 验证 JWT 请求
func (j *JWTHelper) ValidateRequest(req *http.Request) (*jwt.Token, *JWTClaims, error) {
	token, err := request.ParseFromRequest(req, request.AuthorizationHeaderExtractor, func(token *jwt.Token) (any, error) {
		// 检查签名方法
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.publicKey, nil
	}, request.WithClaims(&JWTClaims{}))

	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to parse token")
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return token, claims, nil
	}

	return nil, nil, errors.New("invalid token")
}

// ParseRSAPrivateKeyFromPath 从 文件地址解析 RSA 私钥
func ParseRSAPrivateKeyFromPath(path string) (*rsa.PrivateKey, []byte, error) {
	privateKeyPEM, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		return nil, nil, err

	}
	return privateKey, privateKeyPEM, err
}

// ParseRSAPublicKeyFromPath 从 文件地址解析 RSA 公钥
func ParseRSAPublicKeyFromPath(path string) (*rsa.PublicKey, []byte, error) {
	publicKeyPEM, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyPEM)
	if err != nil {
		return nil, nil, err

	}
	return publicKey, publicKeyPEM, err
}
