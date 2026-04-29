// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package middleware

import (
	"context"
	"net/http"
	"weclawbotnotify/pkg/jwtx"
	"weclawbotnotify/pkg/xerr"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type JWTMiddleware struct {
	helper *jwtx.JWTHelper
	// publicKeyPEM []byte
}

// 定义一个私有类型作为 context 的 key，防止与其他包冲突
type contextKey string

// 定义具体的 key 常量
const ClaimsKey = contextKey("claims")

func NewJWTMiddleware(publicKeyPEM []byte) *JWTMiddleware {
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyPEM)
	if err != nil {
		panic(err)
	}
	helper := jwtx.NewJWTHelper(
		jwtx.WithPublicKey(publicKey), // 设置公钥
	)
	return &JWTMiddleware{
		helper: helper,
	}
}

func (m *JWTMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logx.WithContext(r.Context()).Debug("jwt middleware,r", r)
		token, claims, err := m.helper.ValidateRequest(r)
		if err != nil {
			// 401 返回
			httpx.ErrorCtx(r.Context(), w, xerr.JwtError)
			return
		}
		logx.Debug("token", token)
		logx.Debug("claims", claims)
		if !claims.IsAccessToken() { // 必须是授权Token
			httpx.ErrorCtx(r.Context(), w, xerr.JwtError)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), ClaimsKey, claims))
		next(w, r)
	}
}
