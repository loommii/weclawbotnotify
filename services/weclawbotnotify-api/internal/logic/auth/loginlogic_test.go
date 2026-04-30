package auth

import (
	"context"
	"crypto/rsa"
	"database/sql"
	"errors"
	"testing"
	"time"

	"weclawbotnotify/pkg/jwtx"
	"weclawbotnotify/services/weclawbotnotify-api/internal/model"
	"weclawbotnotify/services/weclawbotnotify-api/internal/svc"
	"weclawbotnotify/services/weclawbotnotify-api/internal/types"

	"golang.org/x/crypto/bcrypt"
)

// hashPassword 对密码进行 bcrypt 哈希
func hashPassword(t *testing.T, password string) string {
	t.Helper()
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("密码哈希失败: %v", err)
	}
	return string(hashed)
}

// TestLoginLogic_Login_Success 测试正常登录流程
// 场景：用户已存在，用户名和密码正确
// 预期：登录成功，返回有效的双令牌和用户信息
func TestLoginLogic_Login_Success(t *testing.T) {
	accessHelper, refreshHelper := generateTestJWTHelpers(t)
	hashedPwd := hashPassword(t, "password123")

	mockModel := &model.MockUsersModel{
		MockFindOneByUsername: func(ctx context.Context, username string) (*model.Users, error) {
			if username == "testuser" {
				return &model.Users{
					Id:           1,
					Username:     "testuser",
					PasswordHash: hashedPwd,
					CreatedAt:    1234567890,
					UpdatedAt:    1234567890,
				}, nil
			}
			return nil, model.ErrNotFound
		},
	}

	svcCtx := &svc.ServiceContext{
		UsersModel:         mockModel,
		RefreshTokensModel: newTestRefreshMock(),
		AccessJWTHelper:    accessHelper,
		RefreshJWTHelper:   refreshHelper,
	}

	logic := NewLoginLogic(context.Background(), svcCtx)
	req := &types.LoginReq{
		Username: "testuser",
		Password: "password123",
	}

	resp, err := logic.Login(req)
	if err != nil {
		t.Fatalf("预期无错误，got %v", err)
	}
	if resp == nil {
		t.Fatal("预期非空响应")
	}
	if resp.Token == "" {
		t.Error("预期非空 Access Token")
	}
	if resp.RefreshToken == "" {
		t.Error("预期非空 Refresh Token")
	}
	if resp.User.Id != 1 {
		t.Errorf("预期 userId=1, got %d", resp.User.Id)
	}
	if resp.User.Username != "testuser" {
		t.Errorf("预期 username=testuser, got %s", resp.User.Username)
	}
}

// TestLoginLogic_Login_EmptyUsername 测试用户名为空
// 预期：返回参数错误
func TestLoginLogic_Login_EmptyUsername(t *testing.T) {
	accessHelper, refreshHelper := generateTestJWTHelpers(t)
	mockModel := &model.MockUsersModel{}

	svcCtx := &svc.ServiceContext{
		UsersModel:         mockModel,
		RefreshTokensModel: newTestRefreshMock(),
		AccessJWTHelper:    accessHelper,
		RefreshJWTHelper:   refreshHelper,
	}

	logic := NewLoginLogic(context.Background(), svcCtx)
	req := &types.LoginReq{
		Username: "",
		Password: "password123",
	}

	_, err := logic.Login(req)
	if err == nil {
		t.Fatal("空用户名应返回错误")
	}
}

// TestLoginLogic_Login_EmptyPassword 测试密码为空
// 预期：返回参数错误
func TestLoginLogic_Login_EmptyPassword(t *testing.T) {
	accessHelper, refreshHelper := generateTestJWTHelpers(t)
	mockModel := &model.MockUsersModel{}

	svcCtx := &svc.ServiceContext{
		UsersModel:         mockModel,
		RefreshTokensModel: newTestRefreshMock(),
		AccessJWTHelper:    accessHelper,
		RefreshJWTHelper:   refreshHelper,
	}

	logic := NewLoginLogic(context.Background(), svcCtx)
	req := &types.LoginReq{
		Username: "testuser",
		Password: "",
	}

	_, err := logic.Login(req)
	if err == nil {
		t.Fatal("空密码应返回错误")
	}
}

// TestLoginLogic_Login_EmptyBoth 测试用户名和密码都为空
// 预期：返回参数错误
func TestLoginLogic_Login_EmptyBoth(t *testing.T) {
	accessHelper, refreshHelper := generateTestJWTHelpers(t)
	mockModel := &model.MockUsersModel{}

	svcCtx := &svc.ServiceContext{
		UsersModel:         mockModel,
		RefreshTokensModel: newTestRefreshMock(),
		AccessJWTHelper:    accessHelper,
		RefreshJWTHelper:   refreshHelper,
	}

	logic := NewLoginLogic(context.Background(), svcCtx)
	req := &types.LoginReq{
		Username: "",
		Password: "",
	}

	_, err := logic.Login(req)
	if err == nil {
		t.Fatal("空用户名和密码应返回错误")
	}
}

// TestLoginLogic_Login_UserNotFound 测试用户不存在
// 预期：返回用户未找到错误
func TestLoginLogic_Login_UserNotFound(t *testing.T) {
	accessHelper, refreshHelper := generateTestJWTHelpers(t)
	mockModel := &model.MockUsersModel{
		MockFindOneByUsername: func(ctx context.Context, username string) (*model.Users, error) {
			return nil, model.ErrNotFound
		},
	}

	svcCtx := &svc.ServiceContext{
		UsersModel:         mockModel,
		RefreshTokensModel: newTestRefreshMock(),
		AccessJWTHelper:    accessHelper,
		RefreshJWTHelper:   refreshHelper,
	}

	logic := NewLoginLogic(context.Background(), svcCtx)
	req := &types.LoginReq{
		Username: "nonexistent",
		Password: "password123",
	}

	_, err := logic.Login(req)
	if err == nil {
		t.Fatal("用户不存在应返回错误")
	}
}

// TestLoginLogic_Login_WrongPassword 测试密码错误
// 预期：返回密码错误
func TestLoginLogic_Login_WrongPassword(t *testing.T) {
	accessHelper, refreshHelper := generateTestJWTHelpers(t)
	hashedPwd := hashPassword(t, "correctpassword")

	mockModel := &model.MockUsersModel{
		MockFindOneByUsername: func(ctx context.Context, username string) (*model.Users, error) {
			if username == "testuser" {
				return &model.Users{
					Id:           1,
					Username:     "testuser",
					PasswordHash: hashedPwd,
					CreatedAt:    1234567890,
					UpdatedAt:    1234567890,
				}, nil
			}
			return nil, model.ErrNotFound
		},
	}

	svcCtx := &svc.ServiceContext{
		UsersModel:         mockModel,
		RefreshTokensModel: newTestRefreshMock(),
		AccessJWTHelper:    accessHelper,
		RefreshJWTHelper:   refreshHelper,
	}

	logic := NewLoginLogic(context.Background(), svcCtx)
	req := &types.LoginReq{
		Username: "testuser",
		Password: "wrongpassword",
	}

	_, err := logic.Login(req)
	if err == nil {
		t.Fatal("密码错误应返回错误")
	}
}

// TestLoginLogic_Login_DBQueryFailed 测试数据库查询失败
// 预期：返回查询失败错误
func TestLoginLogic_Login_DBQueryFailed(t *testing.T) {
	accessHelper, refreshHelper := generateTestJWTHelpers(t)
	mockErr := errors.New("database connection lost")
	mockModel := &model.MockUsersModel{
		MockFindOneByUsername: func(ctx context.Context, username string) (*model.Users, error) {
			return nil, mockErr
		},
	}

	svcCtx := &svc.ServiceContext{
		UsersModel:         mockModel,
		RefreshTokensModel: newTestRefreshMock(),
		AccessJWTHelper:    accessHelper,
		RefreshJWTHelper:   refreshHelper,
	}

	logic := NewLoginLogic(context.Background(), svcCtx)
	req := &types.LoginReq{
		Username: "testuser",
		Password: "password123",
	}

	_, err := logic.Login(req)
	if err == nil {
		t.Fatal("数据库查询失败应返回错误")
	}
}

// TestLoginLogic_Login_TokenGenerationFailed 测试令牌生成失败（JWTHelper 为 nil）
// 预期：触发 panic 或返回错误
func TestLoginLogic_Login_TokenGenerationFailed(t *testing.T) {
	hashedPwd := hashPassword(t, "password123")

	mockModel := &model.MockUsersModel{
		MockFindOneByUsername: func(ctx context.Context, username string) (*model.Users, error) {
			if username == "testuser" {
				return &model.Users{
					Id:           1,
					Username:     "testuser",
					PasswordHash: hashedPwd,
					CreatedAt:    1234567890,
					UpdatedAt:    1234567890,
				}, nil
			}
			return nil, model.ErrNotFound
		},
	}

	svcCtx := &svc.ServiceContext{
		UsersModel:         mockModel,
		RefreshTokensModel: newTestRefreshMock(),
		AccessJWTHelper:    nil,
		RefreshJWTHelper:   nil,
	}

	logic := NewLoginLogic(context.Background(), svcCtx)
	req := &types.LoginReq{
		Username: "testuser",
		Password: "password123",
	}

	defer func() {
		if r := recover(); r == nil {
			t.Log("JWTHelper 为 nil 时应触发 panic 或返回错误")
		}
	}()

	_, err := logic.Login(req)
	if err == nil {
		t.Fatal("JWTHelper 为 nil 应返回错误或 panic")
	}
}

// TestLoginLogic_Login_CaseSensitiveUsername 测试用户名大小写敏感
// 预期：用户名区分大小写
func TestLoginLogic_Login_CaseSensitiveUsername(t *testing.T) {
	accessHelper, refreshHelper := generateTestJWTHelpers(t)
	hashedPwd := hashPassword(t, "password123")

	mockModel := &model.MockUsersModel{
		MockFindOneByUsername: func(ctx context.Context, username string) (*model.Users, error) {
			if username == "testuser" {
				return &model.Users{
					Id:           1,
					Username:     "testuser",
					PasswordHash: hashedPwd,
					CreatedAt:    1234567890,
					UpdatedAt:    1234567890,
				}, nil
			}
			return nil, model.ErrNotFound
		},
	}

	svcCtx := &svc.ServiceContext{
		UsersModel:         mockModel,
		RefreshTokensModel: newTestRefreshMock(),
		AccessJWTHelper:    accessHelper,
		RefreshJWTHelper:   refreshHelper,
	}

	logic := NewLoginLogic(context.Background(), svcCtx)
	req := &types.LoginReq{
		Username: "TestUser",
		Password: "password123",
	}

	_, err := logic.Login(req)
	if err == nil {
		t.Fatal("大小写不匹配应返回错误")
	}
}

// TestLoginLogic_Login_SpecialCharactersInPassword 测试密码包含特殊字符
// 预期：登录流程正常处理特殊字符密码
func TestLoginLogic_Login_SpecialCharactersInPassword(t *testing.T) {
	accessHelper, refreshHelper := generateTestJWTHelpers(t)
	specialPassword := "P@$$w0rd!#$%^&*()_+-=[]{}|;':\",.<>?/~`"
	hashedPwd := hashPassword(t, specialPassword)

	mockModel := &model.MockUsersModel{
		MockFindOneByUsername: func(ctx context.Context, username string) (*model.Users, error) {
			if username == "testuser" {
				return &model.Users{
					Id:           1,
					Username:     "testuser",
					PasswordHash: hashedPwd,
					CreatedAt:    1234567890,
					UpdatedAt:    1234567890,
				}, nil
			}
			return nil, model.ErrNotFound
		},
	}

	svcCtx := &svc.ServiceContext{
		UsersModel:         mockModel,
		RefreshTokensModel: newTestRefreshMock(),
		AccessJWTHelper:    accessHelper,
		RefreshJWTHelper:   refreshHelper,
	}

	logic := NewLoginLogic(context.Background(), svcCtx)
	req := &types.LoginReq{
		Username: "testuser",
		Password: specialPassword,
	}

	resp, err := logic.Login(req)
	if err != nil {
		t.Fatalf("特殊字符密码应正常登录，got %v", err)
	}
	if resp.Token == "" {
		t.Error("预期非空 Access Token")
	}
	if resp.RefreshToken == "" {
		t.Error("预期非空 Refresh Token")
	}
}

// TestLoginLogic_Login_WhitespaceUsername 测试用户名包含空白字符
// 预期：系统不会自动 trim，按原样匹配
func TestLoginLogic_Login_WhitespaceUsername(t *testing.T) {
	accessHelper, refreshHelper := generateTestJWTHelpers(t)
	hashedPwd := hashPassword(t, "password123")

	mockModel := &model.MockUsersModel{
		MockFindOneByUsername: func(ctx context.Context, username string) (*model.Users, error) {
			if username == "testuser" {
				return &model.Users{
					Id:           1,
					Username:     "testuser",
					PasswordHash: hashedPwd,
					CreatedAt:    1234567890,
					UpdatedAt:    1234567890,
				}, nil
			}
			return nil, model.ErrNotFound
		},
	}

	svcCtx := &svc.ServiceContext{
		UsersModel:         mockModel,
		RefreshTokensModel: newTestRefreshMock(),
		AccessJWTHelper:    accessHelper,
		RefreshJWTHelper:   refreshHelper,
	}

	logic := NewLoginLogic(context.Background(), svcCtx)
	req := &types.LoginReq{
		Username: " testuser ",
		Password: "password123",
	}

	_, err := logic.Login(req)
	if err == nil {
		t.Fatal("带空格的用户名应返回错误")
	}
}

// TestLoginLogic_Login_ResponseCreatedAt 测试返回的用户创建时间格式
// 预期：返回的 CreatedAt 字段为 Unix 时间戳字符串格式
func TestLoginLogic_Login_ResponseCreatedAt(t *testing.T) {
	accessHelper, refreshHelper := generateTestJWTHelpers(t)
	hashedPwd := hashPassword(t, "password123")
	createdAt := int64(1234567890)

	mockModel := &model.MockUsersModel{
		MockFindOneByUsername: func(ctx context.Context, username string) (*model.Users, error) {
			if username == "testuser" {
				return &model.Users{
					Id:           1,
					Username:     "testuser",
					PasswordHash: hashedPwd,
					CreatedAt:    createdAt,
					UpdatedAt:    createdAt,
				}, nil
			}
			return nil, model.ErrNotFound
		},
	}

	svcCtx := &svc.ServiceContext{
		UsersModel:         mockModel,
		RefreshTokensModel: newTestRefreshMock(),
		AccessJWTHelper:    accessHelper,
		RefreshJWTHelper:   refreshHelper,
	}

	logic := NewLoginLogic(context.Background(), svcCtx)
	req := &types.LoginReq{
		Username: "testuser",
		Password: "password123",
	}

	resp, err := logic.Login(req)
	if err != nil {
		t.Fatalf("预期无错误，got %v", err)
	}
	if resp.User.CreatedAt != "1234567890" {
		t.Errorf("预期 CreatedAt='1234567890', got %s", resp.User.CreatedAt)
	}
}

// TestLoginLogic_Login_MultipleUsers 测试多用户场景下的登录
// 预期：返回对应用户的信息和双令牌
func TestLoginLogic_Login_MultipleUsers(t *testing.T) {
	accessHelper, refreshHelper := generateTestJWTHelpers(t)
	hashedPwd1 := hashPassword(t, "password1")
	hashedPwd2 := hashPassword(t, "password2")

	mockModel := &model.MockUsersModel{
		MockFindOneByUsername: func(ctx context.Context, username string) (*model.Users, error) {
			switch username {
			case "user1":
				return &model.Users{
					Id:           1,
					Username:     "user1",
					PasswordHash: hashedPwd1,
					CreatedAt:    1234567890,
					UpdatedAt:    1234567890,
				}, nil
			case "user2":
				return &model.Users{
					Id:           2,
					Username:     "user2",
					PasswordHash: hashedPwd2,
					CreatedAt:    1234567891,
					UpdatedAt:    1234567891,
				}, nil
			}
			return nil, model.ErrNotFound
		},
		MockInsert: func(ctx context.Context, data *model.Users) (sql.Result, error) {
			return nil, nil
		},
	}

	svcCtx := &svc.ServiceContext{
		UsersModel:         mockModel,
		RefreshTokensModel: newTestRefreshMock(),
		AccessJWTHelper:    accessHelper,
		RefreshJWTHelper:   refreshHelper,
	}

	logic := NewLoginLogic(context.Background(), svcCtx)

	resp1, err := logic.Login(&types.LoginReq{Username: "user1", Password: "password1"})
	if err != nil {
		t.Fatalf("user1 登录失败: %v", err)
	}
	if resp1.User.Id != 1 || resp1.User.Username != "user1" {
		t.Errorf("预期 user1 信息，got id=%d username=%s", resp1.User.Id, resp1.User.Username)
	}
	if resp1.RefreshToken == "" {
		t.Error("user1 预期非空 Refresh Token")
	}

	resp2, err := logic.Login(&types.LoginReq{Username: "user2", Password: "password2"})
	if err != nil {
		t.Fatalf("user2 登录失败: %v", err)
	}
	if resp2.User.Id != 2 || resp2.User.Username != "user2" {
		t.Errorf("预期 user2 信息，got id=%d username=%s", resp2.User.Id, resp2.User.Username)
	}
	if resp2.RefreshToken == "" {
		t.Error("user2 预期非空 Refresh Token")
	}
}

// TestHashPassword 测试密码哈希辅助函数
func TestHashPassword(t *testing.T) {
	hashed := hashPassword(t, "testpassword")
	if hashed == "" {
		t.Fatal("预期非空哈希值")
	}

	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte("testpassword"))
	if err != nil {
		t.Errorf("哈希密码与原始密码不匹配: %v", err)
	}
}

// generateTestJWTHelpers 生成测试用 Access 和 Refresh JWTHelper
func generateTestJWTHelpers(t *testing.T) (*jwtx.JWTHelper, *jwtx.JWTHelper) {
	t.Helper()
	privateKey, err := rsa.GenerateKey(nil, 2048)
	if err != nil {
		t.Fatalf("生成 RSA 密钥失败: %v", err)
	}
	accessHelper := jwtx.NewJWTHelper(
		jwtx.WithPrivateKey(privateKey),
		jwtx.WithPublicKey(&privateKey.PublicKey),
		jwtx.WithExpiredTime(1*time.Hour),
	)
	refreshHelper := jwtx.NewJWTHelper(
		jwtx.WithPrivateKey(privateKey),
		jwtx.WithPublicKey(&privateKey.PublicKey),
		jwtx.WithExpiredTime(24*time.Hour),
	)
	return accessHelper, refreshHelper
}
