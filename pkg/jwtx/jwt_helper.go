package jwtx

import (
	"context"
	"crypto/rsa"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/golang-jwt/jwt/v5/request"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// JWTHelper JWT 令牌签发与验证工具
type JWTHelper struct {
	publicKey     *rsa.PublicKey    // RSA 公钥，用于验签
	privateKey    *rsa.PrivateKey   // RSA 私钥，用于签名
	expire        time.Duration     // Token 过期时间
	signingMethod jwt.SigningMethod // 签名算法，默认 RS256
}

// NewJWTHelper 构造 JWTHelper，支持函数选项模式
func NewJWTHelper(opts ...Option) *JWTHelper {
	helper := &JWTHelper{
		signingMethod: jwt.SigningMethodRS256,
	}
	for _, opt := range opts {
		opt(helper)
	}
	return helper
}

// GenerateToken 生成 JWT Token（使用配置的过期时间）
func (j *JWTHelper) GenerateToken(claims JWTClaims) (string, error) {
	if claims.IsRefreshToken() && claims.JTI == "" {
		claims.JTI = uuid.New().String()
	}
	return j.generateToken(claims, j.expire)
}

// generateToken 核心签发逻辑：覆盖时间字段，签名返回
func (j *JWTHelper) generateToken(claims JWTClaims, expire time.Duration) (string, error) {
	now := time.Now()
	claims.ExpiresAt = jwt.NewNumericDate(now.Add(expire))
	claims.IssuedAt = jwt.NewNumericDate(now)
	claims.NotBefore = jwt.NewNumericDate(now)

	token := jwt.NewWithClaims(j.signingMethod, claims)
	signedToken, err := token.SignedString(j.privateKey)
	if err != nil {
		return "", errors.Wrap(err, "签名令牌失败")
	}
	return signedToken, nil
}

// ValidateToken 验证 JWT 令牌签名和 Claims
func (j *JWTHelper) ValidateToken(tokenString string) (*jwt.Token, *JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.Errorf("非预期签名算法: %v", token.Header["alg"])
		}
		return j.publicKey, nil
	}, jwt.WithLeeway(5*time.Second))

	if err != nil {
		return nil, nil, errors.Wrap(err, "解析令牌失败")
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return token, claims, nil
	}

	return nil, nil, errors.New("无效令牌")
}

// ValidateRequest 从 HTTP 请求 Authorization 头提取并验证 JWT
func (j *JWTHelper) ValidateRequest(req *http.Request) (*jwt.Token, *JWTClaims, error) {
	token, err := request.ParseFromRequest(req, request.AuthorizationHeaderExtractor, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.Errorf("非预期签名算法: %v", token.Header["alg"])
		}
		return j.publicKey, nil
	}, request.WithClaims(&JWTClaims{}))

	if err != nil {
		return nil, nil, errors.Wrap(err, "解析令牌失败")
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return token, claims, nil
	}

	return nil, nil, errors.New("无效令牌")
}

// ParseRSAPrivateKeyFromPath 从 PEM 文件解析 RSA 私钥
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

// ParseRSAPublicKeyFromPath 从 PEM 文件解析 RSA 公钥
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

// GetUserIdFromContext 从 context 中的 JWT Claims 解析用户 ID
// 需要配合 JWTMiddleware 使用，中间件会将 claims 存入 context
func GetUserIdFromContext(ctx context.Context, claimsKey any) (int64, error) {
	claims, ok := ctx.Value(claimsKey).(*JWTClaims)
	if !ok || claims == nil {
		return 0, fmt.Errorf("jwt claims not found in context")
	}

	userId, err := strconv.ParseInt(claims.UID, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid user id format: %s", claims.UID)
	}

	return userId, nil
}
