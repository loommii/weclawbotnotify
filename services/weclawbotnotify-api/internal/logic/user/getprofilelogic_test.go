package user

import (
	"context"
	"testing"

	"weclawbotnotify/pkg/jwtx"
	"weclawbotnotify/pkg/middleware"
	"weclawbotnotify/pkg/xerr"
	"weclawbotnotify/services/weclawbotnotify-api/internal/model"
	"weclawbotnotify/services/weclawbotnotify-api/internal/svc"
	"weclawbotnotify/services/weclawbotnotify-api/internal/types"
)

// TestGetProfileLogic_GetProfile_Success 测试正常获取用户信息
// 场景：JWT claims 有效，数据库中存在对应用户
// 预期：返回正确的用户信息
func TestGetProfileLogic_GetProfile_Success(t *testing.T) {
	claims := &jwtx.JWTClaims{UID: "1", TokenType: jwtx.Access}
	ctx := context.WithValue(context.Background(), middleware.ClaimsKey, claims)

	mockModel := &model.MockUsersModel{
		MockFindOne: func(ctx context.Context, id int64) (*model.Users, error) {
			if id == 1 {
				return &model.Users{
					Id:           1,
					Username:     "testuser",
					PasswordHash: "hashed_pwd",
					CreatedAt:    1234567890,
					UpdatedAt:    1234567890,
				}, nil
			}
			return nil, model.ErrNotFound
		},
	}

	svcCtx := &svc.ServiceContext{
		UsersModel: mockModel,
	}

	logic := NewGetProfileLogic(ctx, svcCtx)
	resp, err := logic.GetProfile()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if resp.User.Id != 1 {
		t.Errorf("expected user id 1, got %d", resp.User.Id)
	}
	if resp.User.Username != "testuser" {
		t.Errorf("expected username testuser, got %s", resp.User.Username)
	}
	if resp.User.CreatedAt != "1234567890" {
		t.Errorf("expected CreatedAt '1234567890', got %s", resp.User.CreatedAt)
	}
}

// TestGetProfileLogic_GetProfile_MissingClaims 测试 JWT claims 缺失
// 场景：context 中没有 JWT claims
// 预期：返回 JwtError
func TestGetProfileLogic_GetProfile_MissingClaims(t *testing.T) {
	ctx := context.Background()

	svcCtx := &svc.ServiceContext{
		UsersModel: &model.MockUsersModel{},
	}

	logic := NewGetProfileLogic(ctx, svcCtx)
	_, err := logic.GetProfile()
	if err == nil {
		t.Fatal("expected error for missing claims, got nil")
	}

	codeErr, ok := err.(*xerr.CodeError)
	if !ok {
		t.Fatalf("expected *xerr.CodeError, got %T", err)
	}
	if codeErr.Code != 100001 {
		t.Errorf("expected JwtError code 100001, got %d", codeErr.Code)
	}
}

// TestGetProfileLogic_GetProfile_InvalidUserId 测试 userId 格式无效
// 场景：JWT claims 中的 UID 不是有效数字
// 预期：返回 JwtError
func TestGetProfileLogic_GetProfile_InvalidUserId(t *testing.T) {
	claims := &jwtx.JWTClaims{UID: "not_a_number", TokenType: jwtx.Access}
	ctx := context.WithValue(context.Background(), middleware.ClaimsKey, claims)

	svcCtx := &svc.ServiceContext{
		UsersModel: &model.MockUsersModel{},
	}

	logic := NewGetProfileLogic(ctx, svcCtx)
	_, err := logic.GetProfile()
	if err == nil {
		t.Fatal("expected error for invalid userId, got nil")
	}

	codeErr, ok := err.(*xerr.CodeError)
	if !ok {
		t.Fatalf("expected *xerr.CodeError, got %T", err)
	}
	if codeErr.Code != 100001 {
		t.Errorf("expected JwtError code 100001, got %d", codeErr.Code)
	}
}

// TestGetProfileLogic_GetProfile_UserNotFound 测试用户不存在
// 场景：JWT claims 有效，但数据库中找不到对应用户
// 预期：返回 JwtError
func TestGetProfileLogic_GetProfile_UserNotFound(t *testing.T) {
	claims := &jwtx.JWTClaims{UID: "999", TokenType: jwtx.Access}
	ctx := context.WithValue(context.Background(), middleware.ClaimsKey, claims)

	mockModel := &model.MockUsersModel{
		MockFindOne: func(ctx context.Context, id int64) (*model.Users, error) {
			return nil, model.ErrNotFound
		},
	}

	svcCtx := &svc.ServiceContext{
		UsersModel: mockModel,
	}

	logic := NewGetProfileLogic(ctx, svcCtx)
	_, err := logic.GetProfile()
	if err == nil {
		t.Fatal("expected error for user not found, got nil")
	}

	codeErr, ok := err.(*xerr.CodeError)
	if !ok {
		t.Fatalf("expected *xerr.CodeError, got %T", err)
	}
	if codeErr.Code != 100001 {
		t.Errorf("expected JwtError code 100001, got %d", codeErr.Code)
	}
}

// TestGetProfileLogic_GetProfile_DBError 测试数据库查询异常
// 场景：JWT claims 有效，但数据库查询发生未知错误
// 预期：返回查询失败错误
func TestGetProfileLogic_GetProfile_DBError(t *testing.T) {
	claims := &jwtx.JWTClaims{UID: "1", TokenType: jwtx.Access}
	ctx := context.WithValue(context.Background(), middleware.ClaimsKey, claims)

	dbErr := model.ErrNotFound
	mockModel := &model.MockUsersModel{
		MockFindOne: func(ctx context.Context, id int64) (*model.Users, error) {
			return nil, dbErr
		},
	}

	svcCtx := &svc.ServiceContext{
		UsersModel: mockModel,
	}

	logic := NewGetProfileLogic(ctx, svcCtx)
	_, err := logic.GetProfile()
	if err == nil {
		t.Fatal("expected error for DB error, got nil")
	}

	codeErr, ok := err.(*xerr.CodeError)
	if !ok {
		t.Fatalf("expected *xerr.CodeError, got %T", err)
	}
	if codeErr.Code != 100001 {
		t.Errorf("expected JwtError code 100001, got %d", codeErr.Code)
	}
}

// TestGetProfileLogic_GetProfile_NilClaims 测试 claims 为 nil
// 场景：context 中 ClaimsKey 对应的值为 nil
// 预期：返回 JwtError
func TestGetProfileLogic_GetProfile_NilClaims(t *testing.T) {
	ctx := context.WithValue(context.Background(), middleware.ClaimsKey, nil)

	svcCtx := &svc.ServiceContext{
		UsersModel: &model.MockUsersModel{},
	}

	logic := NewGetProfileLogic(ctx, svcCtx)
	_, err := logic.GetProfile()
	if err == nil {
		t.Fatal("expected error for nil claims, got nil")
	}
}

// TestGetProfileLogic_GetProfile_WrongClaimType 测试 claims 类型不匹配
// 场景：context 中 ClaimsKey 对应的值不是 *jwtx.JWTClaims 类型
// 预期：返回 JwtError
func TestGetProfileLogic_GetProfile_WrongClaimType(t *testing.T) {
	ctx := context.WithValue(context.Background(), middleware.ClaimsKey, "not_a_claim")

	svcCtx := &svc.ServiceContext{
		UsersModel: &model.MockUsersModel{},
	}

	logic := NewGetProfileLogic(ctx, svcCtx)
	_, err := logic.GetProfile()
	if err == nil {
		t.Fatal("expected error for wrong claim type, got nil")
	}
}

// TestGetProfileLogic_GetProfile_LongUsername 测试长用户名
// 场景：用户名较长，但格式合法
// 预期：正常返回用户信息
func TestGetProfileLogic_GetProfile_LongUsername(t *testing.T) {
	longUsername := "a_very_long_username_that_is_still_valid_for_testing"
	claims := &jwtx.JWTClaims{UID: "1", TokenType: jwtx.Access}
	ctx := context.WithValue(context.Background(), middleware.ClaimsKey, claims)

	mockModel := &model.MockUsersModel{
		MockFindOne: func(ctx context.Context, id int64) (*model.Users, error) {
			return &model.Users{
				Id:           1,
				Username:     longUsername,
				PasswordHash: "hashed_pwd",
				CreatedAt:    1234567890,
				UpdatedAt:    1234567890,
			}, nil
		},
	}

	svcCtx := &svc.ServiceContext{
		UsersModel: mockModel,
	}

	logic := NewGetProfileLogic(ctx, svcCtx)
	resp, err := logic.GetProfile()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.User.Username != longUsername {
		t.Errorf("expected username %s, got %s", longUsername, resp.User.Username)
	}
}

// TestGetProfileLogic_GetProfile_ResponseType 测试返回值类型和结构
// 场景：正常获取用户信息
// 预期：返回类型是 *types.ProfileResp，结构与 UML 中 UserInfo 定义一致
func TestGetProfileLogic_GetProfile_ResponseType(t *testing.T) {
	claims := &jwtx.JWTClaims{UID: "42", TokenType: jwtx.Access}
	ctx := context.WithValue(context.Background(), middleware.ClaimsKey, claims)

	mockModel := &model.MockUsersModel{
		MockFindOne: func(ctx context.Context, id int64) (*model.Users, error) {
			return &model.Users{
				Id:           42,
				Username:     "profile_test_user",
				PasswordHash: "hashed_pwd",
				CreatedAt:    9876543210,
				UpdatedAt:    9876543210,
			}, nil
		},
	}

	svcCtx := &svc.ServiceContext{
		UsersModel: mockModel,
	}

	logic := NewGetProfileLogic(ctx, svcCtx)
	result, err := logic.GetProfile()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if _, ok := interface{}(result).(*types.ProfileResp); !ok {
		t.Fatal("expected response type to be *types.ProfileResp")
	}

	expected := &types.ProfileResp{
		User: types.UserInfo{
			Id:        42,
			Username:  "profile_test_user",
			CreatedAt: "9876543210",
		},
	}

	if result.User.Id != expected.User.Id {
		t.Errorf("expected user id %d, got %d", expected.User.Id, result.User.Id)
	}
	if result.User.Username != expected.User.Username {
		t.Errorf("expected username %s, got %s", expected.User.Username, result.User.Username)
	}
	if result.User.CreatedAt != expected.User.CreatedAt {
		t.Errorf("expected CreatedAt %s, got %s", expected.User.CreatedAt, result.User.CreatedAt)
	}
}
