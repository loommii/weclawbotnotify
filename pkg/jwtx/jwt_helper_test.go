package jwtx

import (
	"crypto/rsa"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func generateRSAKeys(t *testing.T) (*rsa.PrivateKey, *rsa.PublicKey) {
	t.Helper()
	privateKey, err := rsa.GenerateKey(nil, 2048)
	if err != nil {
		t.Fatalf("failed to generate RSA key: %v", err)
	}
	return privateKey, &privateKey.PublicKey
}

func TestJWTClaims_IsAccessToken(t *testing.T) {
	tests := []struct {
		name      string
		tokenType string
		expected  bool
	}{
		{"access token", Access, true},
		{"refresh token", Refresh, false},
		{"empty type", "", false},
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

func TestJWTClaims_IsRefreshToken(t *testing.T) {
	tests := []struct {
		name      string
		tokenType string
		expected  bool
	}{
		{"refresh token", Refresh, true},
		{"access token", Access, false},
		{"empty type", "", false},
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

func TestWithPrivateKey(t *testing.T) {
	privateKey, _ := generateRSAKeys(t)
	opt := WithPrivateKey(privateKey)
	helper := &JWTHelper{}
	opt(helper)
	if helper.privateKey != privateKey {
		t.Error("WithPrivateKey did not set privateKey")
	}
}

func TestWithPublicKey(t *testing.T) {
	_, publicKey := generateRSAKeys(t)
	opt := WithPublicKey(publicKey)
	helper := &JWTHelper{}
	opt(helper)
	if helper.publicKey != publicKey {
		t.Error("WithPublicKey did not set publicKey")
	}
}

func TestWithExpiredTime(t *testing.T) {
	duration := 2 * time.Hour
	opt := WithExpiredTime(duration)
	helper := &JWTHelper{}
	opt(helper)
	if helper.expiredTime != duration {
		t.Errorf("WithExpiredTime: got %v, want %v", helper.expiredTime, duration)
	}
}

func TestWithSigningMethod(t *testing.T) {
	method := jwt.SigningMethodRS512
	opt := WithSigningMethod(method)
	helper := &JWTHelper{}
	opt(helper)
	if helper.signingMethod != method {
		t.Error("WithSigningMethod did not set signingMethod")
	}
}

func TestNewJWTHelper_Defaults(t *testing.T) {
	helper := NewJWTHelper()
	if helper.signingMethod != jwt.SigningMethodRS256 {
		t.Errorf("default signing method = %v, want RS256", helper.signingMethod)
	}
}

func TestNewJWTHelper_WithOptions(t *testing.T) {
	privateKey, publicKey := generateRSAKeys(t)
	expired := 1 * time.Hour

	helper := NewJWTHelper(
		WithPrivateKey(privateKey),
		WithPublicKey(publicKey),
		WithExpiredTime(expired),
	)

	if helper.privateKey != privateKey {
		t.Error("privateKey not set")
	}
	if helper.publicKey != publicKey {
		t.Error("publicKey not set")
	}
	if helper.expiredTime != expired {
		t.Errorf("expiredTime = %v, want %v", helper.expiredTime, expired)
	}
}

func TestGenerateAndValidateToken(t *testing.T) {
	privateKey, publicKey := generateRSAKeys(t)
	helper := NewJWTHelper(
		WithPrivateKey(privateKey),
		WithPublicKey(publicKey),
		WithExpiredTime(1*time.Hour),
	)

	claims := JWTClaims{
		UID:       "user123",
		TokenType: Access,
	}

	tokenStr, err := helper.GenerateToken(claims)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}
	if tokenStr == "" {
		t.Fatal("GenerateToken returned empty token")
	}

	parsedToken, parsedClaims, err := helper.ValidateToken(tokenStr)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}
	if !parsedToken.Valid {
		t.Error("parsed token is not valid")
	}
	if parsedClaims.UID != "user123" {
		t.Errorf("UID = %v, want user123", parsedClaims.UID)
	}
	if !parsedClaims.IsAccessToken() {
		t.Error("expected access token")
	}
}

func TestValidateToken_Expired(t *testing.T) {
	privateKey, publicKey := generateRSAKeys(t)
	helper := NewJWTHelper(
		WithPrivateKey(privateKey),
		WithPublicKey(publicKey),
		WithExpiredTime(-10*time.Second),
	)

	claims := JWTClaims{UID: "user123", TokenType: Access}
	tokenStr, err := helper.GenerateToken(claims)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	_, _, err = helper.ValidateToken(tokenStr)
	if err == nil {
		t.Error("expected error for expired token, got nil")
	}
}

func TestValidateToken_WrongPublicKey(t *testing.T) {
	privateKey, _ := generateRSAKeys(t)
	_, wrongPublicKey := generateRSAKeys(t)

	signHelper := NewJWTHelper(
		WithPrivateKey(privateKey),
		WithExpiredTime(1*time.Hour),
	)

	claims := JWTClaims{UID: "user123", TokenType: Access}
	tokenStr, err := signHelper.GenerateToken(claims)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	validateHelper := NewJWTHelper(WithPublicKey(wrongPublicKey))
	_, _, err = validateHelper.ValidateToken(tokenStr)
	if err == nil {
		t.Error("expected error for wrong public key, got nil")
	}
}

func TestGenerateToken_RefreshToken(t *testing.T) {
	privateKey, publicKey := generateRSAKeys(t)
	helper := NewJWTHelper(
		WithPrivateKey(privateKey),
		WithPublicKey(publicKey),
		WithExpiredTime(24*time.Hour),
	)

	claims := JWTClaims{UID: "user456", TokenType: Refresh}
	tokenStr, err := helper.GenerateToken(claims)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	_, parsedClaims, err := helper.ValidateToken(tokenStr)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}
	if !parsedClaims.IsRefreshToken() {
		t.Error("expected refresh token")
	}
	if parsedClaims.IsAccessToken() {
		t.Error("should not be access token")
	}
}

func TestGenerateToken_OverwritesTime(t *testing.T) {
	privateKey, publicKey := generateRSAKeys(t)
	helper := NewJWTHelper(
		WithPrivateKey(privateKey),
		WithPublicKey(publicKey),
		WithExpiredTime(1*time.Hour),
	)

	pastTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	claims := JWTClaims{
		UID:              "user789",
		TokenType:        Access,
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(pastTime)},
	}

	tokenStr, err := helper.GenerateToken(claims)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	_, parsedClaims, err := helper.ValidateToken(tokenStr)
	if err != nil {
		t.Fatalf("ValidateToken failed (time should be overwritten): %v", err)
	}
	if parsedClaims.ExpiresAt.Before(time.Now()) {
		t.Error("ExpiresAt should be overwritten to future time")
	}
}
