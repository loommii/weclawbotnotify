package auth

import (
	"context"
	"crypto/rsa"
	"database/sql"
	"errors"
	"testing"
	"time"

	"weclawbotnotify/pkg/jwtx"
	"weclawbotnotify/pkg/xerr"
	"weclawbotnotify/services/weclawbotnotify-api/internal/model"
	"weclawbotnotify/services/weclawbotnotify-api/internal/svc"
	"weclawbotnotify/services/weclawbotnotify-api/internal/types"
)

// generateRefreshTestJWTHelpers 生成测试用 Access 和 Refresh JWTHelper
func generateRefreshTestJWTHelpers(t *testing.T) (*jwtx.JWTHelper, *jwtx.JWTHelper) {
	t.Helper()
	privateKey, err := rsa.GenerateKey(nil, 2048)
	if err != nil {
		t.Fatalf("生成 RSA 密钥失败: %v", err)
	}
	accessHelper := jwtx.NewJWTHelper(
		jwtx.WithPrivateKey(privateKey),
		jwtx.WithPublicKey(&privateKey.PublicKey),
		jwtx.WithExpiredTime(15*time.Minute),
	)
	refreshHelper := jwtx.NewJWTHelper(
		jwtx.WithPrivateKey(privateKey),
		jwtx.WithPublicKey(&privateKey.PublicKey),
		jwtx.WithExpiredTime(7*24*time.Hour),
	)
	return accessHelper, refreshHelper
}

// TestRefreshLogic_Refresh_Success 测试正常刷新令牌流程
// 场景：提交有效的 Refresh Token，JTI 存在且未撤销、未过期
// 预期：返回新的 TokenPair，旧 JTI 被撤销
func TestRefreshLogic_Refresh_Success(t *testing.T) {
	accessHelper, refreshHelper := generateRefreshTestJWTHelpers(t)

	refreshToken, err := refreshHelper.GenerateToken(jwtx.JWTClaims{
		UID:       "1",
		TokenType: jwtx.Refresh,
	})
	if err != nil {
		t.Fatalf("GenerateToken 失败: %v", err)
	}

	_, claims, _ := refreshHelper.ValidateToken(refreshToken)

	revokedJti := ""
	mockRefreshModel := &model.MockRefreshTokensModel{
		MockFindOneByJti: func(ctx context.Context, jti string) (*model.RefreshTokens, error) {
			return &model.RefreshTokens{
				Id:        1,
				UserId:    1,
				Jti:       claims.JTI,
				Revoked:   0,
				ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
			}, nil
		},
		MockRevokeByJti: func(ctx context.Context, jti string) error {
			revokedJti = jti
			return nil
		},
		MockInsert: func(ctx context.Context, data *model.RefreshTokens) (sql.Result, error) {
			return &model.MockResult{LastInsertIdVal: 2}, nil
		},
	}

	svcCtx := &svc.ServiceContext{
		AccessJWTHelper:    accessHelper,
		RefreshJWTHelper:   refreshHelper,
		RefreshTokensModel: mockRefreshModel,
	}

	logic := NewRefreshLogic(context.Background(), svcCtx)
	resp, err := logic.Refresh(&types.RefreshReq{RefreshToken: refreshToken})
	if err != nil {
		t.Fatalf("预期无错误，got %v", err)
	}
	if resp.Token == "" {
		t.Error("新 Access Token 不应为空")
	}
	if resp.RefreshToken == "" {
		t.Error("新 Refresh Token 不应为空")
	}
	if resp.RefreshToken == refreshToken {
		t.Error("新 Refresh Token 应与旧的不同（轮换）")
	}
	if revokedJti != claims.JTI {
		t.Error("旧 JTI 应被撤销")
	}
}

// TestRefreshLogic_Refresh_EmptyToken 测试空令牌
// 场景：提交空的 Refresh Token
// 预期：返回 RefreshTokenInvalid 错误
func TestRefreshLogic_Refresh_EmptyToken(t *testing.T) {
	accessHelper, refreshHelper := generateRefreshTestJWTHelpers(t)
	svcCtx := &svc.ServiceContext{
		AccessJWTHelper:    accessHelper,
		RefreshJWTHelper:   refreshHelper,
		RefreshTokensModel: &model.MockRefreshTokensModel{},
	}

	logic := NewRefreshLogic(context.Background(), svcCtx)
	_, err := logic.Refresh(&types.RefreshReq{RefreshToken: ""})
	if err == nil {
		t.Fatal("空令牌应返回错误")
	}
	if !errors.Is(err, xerr.RefreshTokenInvalid) {
		t.Errorf("预期 RefreshTokenInvalid，got %v", err)
	}
}

// TestRefreshLogic_Refresh_InvalidToken 测试无效签名令牌
// 场景：提交签名无效的令牌
// 预期：返回 RefreshTokenInvalid 错误
func TestRefreshLogic_Refresh_InvalidToken(t *testing.T) {
	accessHelper, refreshHelper := generateRefreshTestJWTHelpers(t)
	svcCtx := &svc.ServiceContext{
		AccessJWTHelper:    accessHelper,
		RefreshJWTHelper:   refreshHelper,
		RefreshTokensModel: &model.MockRefreshTokensModel{},
	}

	logic := NewRefreshLogic(context.Background(), svcCtx)
	_, err := logic.Refresh(&types.RefreshReq{RefreshToken: "invalid.token.here"})
	if err == nil {
		t.Fatal("无效令牌应返回错误")
	}
	if !errors.Is(err, xerr.RefreshTokenInvalid) {
		t.Errorf("预期 RefreshTokenInvalid，got %v", err)
	}
}

// TestRefreshLogic_Refresh_AccessTokenRejected 测试用 Access Token 刷新
// 场景：提交 Access Token 而非 Refresh Token
// 预期：返回 RefreshTokenInvalid 错误（TokenType 不匹配）
func TestRefreshLogic_Refresh_AccessTokenRejected(t *testing.T) {
	accessHelper, refreshHelper := generateRefreshTestJWTHelpers(t)

	accessToken, err := accessHelper.GenerateToken(jwtx.JWTClaims{
		UID:       "1",
		TokenType: jwtx.Access,
	})
	if err != nil {
		t.Fatalf("GenerateToken 失败: %v", err)
	}

	svcCtx := &svc.ServiceContext{
		AccessJWTHelper:    accessHelper,
		RefreshJWTHelper:   refreshHelper,
		RefreshTokensModel: &model.MockRefreshTokensModel{},
	}

	logic := NewRefreshLogic(context.Background(), svcCtx)
	_, err = logic.Refresh(&types.RefreshReq{RefreshToken: accessToken})
	if err == nil {
		t.Fatal("Access Token 不应用于刷新")
	}
	if !errors.Is(err, xerr.RefreshTokenInvalid) {
		t.Errorf("预期 RefreshTokenInvalid，got %v", err)
	}
}

// TestRefreshLogic_Refresh_RevokedJTI 测试已撤销的 JTI
// 场景：Refresh Token 签名有效但 JTI 不在数据库中（已轮换/撤销）
// 预期：返回 RefreshTokenRevoked 错误
func TestRefreshLogic_Refresh_RevokedJTI(t *testing.T) {
	accessHelper, refreshHelper := generateRefreshTestJWTHelpers(t)

	refreshToken, err := refreshHelper.GenerateToken(jwtx.JWTClaims{
		UID:       "1",
		TokenType: jwtx.Refresh,
	})
	if err != nil {
		t.Fatalf("GenerateToken 失败: %v", err)
	}

	mockRefreshModel := &model.MockRefreshTokensModel{
		MockFindOneByJti: func(ctx context.Context, jti string) (*model.RefreshTokens, error) {
			return nil, model.ErrNotFound
		},
	}

	svcCtx := &svc.ServiceContext{
		AccessJWTHelper:    accessHelper,
		RefreshJWTHelper:   refreshHelper,
		RefreshTokensModel: mockRefreshModel,
	}

	logic := NewRefreshLogic(context.Background(), svcCtx)
	_, err = logic.Refresh(&types.RefreshReq{RefreshToken: refreshToken})
	if err == nil {
		t.Fatal("已撤销 JTI 应返回错误")
	}
	if !errors.Is(err, xerr.RefreshTokenRevoked) {
		t.Errorf("预期 RefreshTokenRevoked，got %v", err)
	}
}

// TestRefreshLogic_Refresh_TokenMarkedRevoked 测试 revoked=1 的令牌
// 场景：Refresh Token 在数据库中已被标记为 revoked=1
// 预期：返回 RefreshTokenRevoked 错误
func TestRefreshLogic_Refresh_TokenMarkedRevoked(t *testing.T) {
	accessHelper, refreshHelper := generateRefreshTestJWTHelpers(t)

	refreshToken, err := refreshHelper.GenerateToken(jwtx.JWTClaims{
		UID:       "1",
		TokenType: jwtx.Refresh,
	})
	if err != nil {
		t.Fatalf("GenerateToken 失败: %v", err)
	}

	_, claims, _ := refreshHelper.ValidateToken(refreshToken)

	mockRefreshModel := &model.MockRefreshTokensModel{
		MockFindOneByJti: func(ctx context.Context, jti string) (*model.RefreshTokens, error) {
			return &model.RefreshTokens{
				Id:        1,
				UserId:    1,
				Jti:       claims.JTI,
				Revoked:   1,
				ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
			}, nil
		},
	}

	svcCtx := &svc.ServiceContext{
		AccessJWTHelper:    accessHelper,
		RefreshJWTHelper:   refreshHelper,
		RefreshTokensModel: mockRefreshModel,
	}

	logic := NewRefreshLogic(context.Background(), svcCtx)
	_, err = logic.Refresh(&types.RefreshReq{RefreshToken: refreshToken})
	if err == nil {
		t.Fatal("已撤销令牌应返回错误")
	}
	if !errors.Is(err, xerr.RefreshTokenRevoked) {
		t.Errorf("预期 RefreshTokenRevoked，got %v", err)
	}
}

// TestRefreshLogic_Refresh_ExpiredToken 测试过期令牌
// 场景：Refresh Token 在数据库中已过期
// 预期：返回 RefreshTokenExpired 错误
func TestRefreshLogic_Refresh_ExpiredToken(t *testing.T) {
	accessHelper, refreshHelper := generateRefreshTestJWTHelpers(t)

	refreshToken, err := refreshHelper.GenerateToken(jwtx.JWTClaims{
		UID:       "1",
		TokenType: jwtx.Refresh,
	})
	if err != nil {
		t.Fatalf("GenerateToken 失败: %v", err)
	}

	_, claims, _ := refreshHelper.ValidateToken(refreshToken)

	mockRefreshModel := &model.MockRefreshTokensModel{
		MockFindOneByJti: func(ctx context.Context, jti string) (*model.RefreshTokens, error) {
			return &model.RefreshTokens{
				Id:        1,
				UserId:    1,
				Jti:       claims.JTI,
				Revoked:   0,
				ExpiresAt: time.Now().Add(-1 * time.Hour).Unix(),
			}, nil
		},
	}

	svcCtx := &svc.ServiceContext{
		AccessJWTHelper:    accessHelper,
		RefreshJWTHelper:   refreshHelper,
		RefreshTokensModel: mockRefreshModel,
	}

	logic := NewRefreshLogic(context.Background(), svcCtx)
	_, err = logic.Refresh(&types.RefreshReq{RefreshToken: refreshToken})
	if err == nil {
		t.Fatal("过期令牌应返回错误")
	}
	if !errors.Is(err, xerr.RefreshTokenExpired) {
		t.Errorf("预期 RefreshTokenExpired，got %v", err)
	}
}

// TestRefreshLogic_Refresh_RevokeOldFailed 测试撤销旧 JTI 失败
// 场景：撤销旧 JTI 时数据库报错
// 预期：返回 RefreshTokenInvalid 错误
func TestRefreshLogic_Refresh_RevokeOldFailed(t *testing.T) {
	accessHelper, refreshHelper := generateRefreshTestJWTHelpers(t)

	refreshToken, err := refreshHelper.GenerateToken(jwtx.JWTClaims{
		UID:       "1",
		TokenType: jwtx.Refresh,
	})
	if err != nil {
		t.Fatalf("GenerateToken 失败: %v", err)
	}

	_, claims, _ := refreshHelper.ValidateToken(refreshToken)

	mockRefreshModel := &model.MockRefreshTokensModel{
		MockFindOneByJti: func(ctx context.Context, jti string) (*model.RefreshTokens, error) {
			return &model.RefreshTokens{
				Id:        1,
				UserId:    1,
				Jti:       claims.JTI,
				Revoked:   0,
				ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
			}, nil
		},
		MockRevokeByJti: func(ctx context.Context, jti string) error {
			return errors.New("数据库错误")
		},
	}

	svcCtx := &svc.ServiceContext{
		AccessJWTHelper:    accessHelper,
		RefreshJWTHelper:   refreshHelper,
		RefreshTokensModel: mockRefreshModel,
	}

	logic := NewRefreshLogic(context.Background(), svcCtx)
	_, err = logic.Refresh(&types.RefreshReq{RefreshToken: refreshToken})
	if err == nil {
		t.Fatal("撤销失败应返回错误")
	}
}
