// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package middleware

import (
	"net/http"

	pkgmw "weclawbotnotify/pkg/middleware"
)

type ClientAuthMiddleware struct {
	jwtMiddleware *pkgmw.JWTMiddleware
}

func NewClientAuthMiddleware(publicKeyPEM []byte) *ClientAuthMiddleware {
	return &ClientAuthMiddleware{
		jwtMiddleware: pkgmw.NewJWTMiddleware(publicKeyPEM),
	}
}

func (m *ClientAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return m.jwtMiddleware.Handle(next)
}
