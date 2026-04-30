package middleware

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"weclawbotnotify/pkg/jwtx"
)

func generateTestKeys(t *testing.T) (*rsa.PrivateKey, []byte) {
	t.Helper()
	privateKey, err := rsa.GenerateKey(nil, 2048)
	if err != nil {
		t.Fatalf("failed to generate RSA key: %v", err)
	}
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		t.Fatalf("failed to marshal public key: %v", err)
	}
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyBytes,
	})
	return privateKey, publicKeyPEM
}

func TestNewJWTMiddleware(t *testing.T) {
	_, publicKeyPEM := generateTestKeys(t)
	m := NewJWTMiddleware(publicKeyPEM)
	if m == nil {
		t.Fatal("NewJWTMiddleware returned nil")
	}
	if m.helper == nil {
		t.Fatal("helper should not be nil")
	}
}

func TestJWTMiddleware_ValidAccessToken(t *testing.T) {
	privateKey, publicKeyPEM := generateTestKeys(t)
	m := NewJWTMiddleware(publicKeyPEM)

	helper := jwtx.NewJWTHelper(
		jwtx.WithPrivateKey(privateKey),
		jwtx.WithExpiredTime(1*time.Hour),
	)

	claims := jwtx.JWTClaims{UID: "user123", TokenType: jwtx.Access}
	tokenStr, err := helper.GenerateToken(claims)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStr)
	rec := httptest.NewRecorder()

	m.Handle(next).ServeHTTP(rec, req)

	if !called {
		t.Error("next handler was not called for valid access token")
	}
}

func TestJWTMiddleware_RefreshTokenRejected(t *testing.T) {
	privateKey, publicKeyPEM := generateTestKeys(t)
	m := NewJWTMiddleware(publicKeyPEM)

	helper := jwtx.NewJWTHelper(
		jwtx.WithPrivateKey(privateKey),
		jwtx.WithExpiredTime(24*time.Hour),
	)

	claims := jwtx.JWTClaims{UID: "user123", TokenType: jwtx.Refresh}
	tokenStr, err := helper.GenerateToken(claims)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStr)
	rec := httptest.NewRecorder()

	m.Handle(next).ServeHTTP(rec, req)

	if called {
		t.Error("next handler should not be called for refresh token")
	}
}

func TestJWTMiddleware_NoAuthHeader(t *testing.T) {
	_, publicKeyPEM := generateTestKeys(t)
	m := NewJWTMiddleware(publicKeyPEM)

	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	m.Handle(next).ServeHTTP(rec, req)

	if called {
		t.Error("next handler should not be called without auth header")
	}
}

func TestJWTMiddleware_InvalidToken(t *testing.T) {
	_, publicKeyPEM := generateTestKeys(t)
	m := NewJWTMiddleware(publicKeyPEM)

	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")
	rec := httptest.NewRecorder()

	m.Handle(next).ServeHTTP(rec, req)

	if called {
		t.Error("next handler should not be called for invalid token")
	}
}

func TestJWTMiddleware_ClaimsInContext(t *testing.T) {
	privateKey, publicKeyPEM := generateTestKeys(t)
	m := NewJWTMiddleware(publicKeyPEM)

	helper := jwtx.NewJWTHelper(
		jwtx.WithPrivateKey(privateKey),
		jwtx.WithExpiredTime(1*time.Hour),
	)

	claims := jwtx.JWTClaims{UID: "user123", TokenType: jwtx.Access}
	tokenStr, err := helper.GenerateToken(claims)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	var gotClaims *jwtx.JWTClaims
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if val := r.Context().Value(ClaimsKey); val != nil {
			gotClaims = val.(*jwtx.JWTClaims)
		}
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStr)
	rec := httptest.NewRecorder()

	m.Handle(next).ServeHTTP(rec, req)

	if gotClaims == nil {
		t.Fatal("claims not found in context")
	}
	if gotClaims.UID != "user123" {
		t.Errorf("UID = %v, want user123", gotClaims.UID)
	}
}
