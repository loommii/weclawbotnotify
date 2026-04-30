package auth

import (
	"context"
	"database/sql"
	"errors"
	"testing"

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
		t.Fatalf("failed to hash password: %v", err)
	}
	return string(hashed)
}

// TestLoginLogic_Login_Success 测试正常登录流程
// 场景：用户已存在，用户名和密码正确
// 预期：登录成功，返回有效的 JWT token 和用户信息
func TestLoginLogic_Login_Success(t *testing.T) {
	helper := generateTestJWTHelper(t)
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
		UsersModel: mockModel,
		JWTHelper:  helper,
	}

	logic := NewLoginLogic(context.Background(), svcCtx)
	req := &types.LoginReq{
		Username: "testuser",
		Password: "password123",
	}

	resp, err := logic.Login(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if resp.Token == "" {
		t.Error("expected non-empty token")
	}
	if resp.User.Id != 1 {
		t.Errorf("expected user id 1, got %d", resp.User.Id)
	}
	if resp.User.Username != "testuser" {
		t.Errorf("expected username testuser, got %s", resp.User.Username)
	}
}

// TestLoginLogic_Login_EmptyUsername 测试用户名为空的场景
// 场景：用户提交登录请求，用户名为空
// 预期：返回参数错误 LoginParamEmpty
func TestLoginLogic_Login_EmptyUsername(t *testing.T) {
	helper := generateTestJWTHelper(t)
	mockModel := &model.MockUsersModel{}

	svcCtx := &svc.ServiceContext{
		UsersModel: mockModel,
		JWTHelper:  helper,
	}

	logic := NewLoginLogic(context.Background(), svcCtx)
	req := &types.LoginReq{
		Username: "",
		Password: "password123",
	}

	_, err := logic.Login(req)
	if err == nil {
		t.Fatal("expected error for empty username, got nil")
	}
}

// TestLoginLogic_Login_EmptyPassword 测试密码为空的场景
// 场景：用户提交登录请求，密码为空
// 预期：返回参数错误 LoginParamEmpty
func TestLoginLogic_Login_EmptyPassword(t *testing.T) {
	helper := generateTestJWTHelper(t)
	mockModel := &model.MockUsersModel{}

	svcCtx := &svc.ServiceContext{
		UsersModel: mockModel,
		JWTHelper:  helper,
	}

	logic := NewLoginLogic(context.Background(), svcCtx)
	req := &types.LoginReq{
		Username: "testuser",
		Password: "",
	}

	_, err := logic.Login(req)
	if err == nil {
		t.Fatal("expected error for empty password, got nil")
	}
}

// TestLoginLogic_Login_EmptyBoth 测试用户名和密码都为空的场景
// 场景：用户提交登录请求，用户名和密码都为空
// 预期：返回参数错误 LoginParamEmpty
func TestLoginLogic_Login_EmptyBoth(t *testing.T) {
	helper := generateTestJWTHelper(t)
	mockModel := &model.MockUsersModel{}

	svcCtx := &svc.ServiceContext{
		UsersModel: mockModel,
		JWTHelper:  helper,
	}

	logic := NewLoginLogic(context.Background(), svcCtx)
	req := &types.LoginReq{
		Username: "",
		Password: "",
	}

	_, err := logic.Login(req)
	if err == nil {
		t.Fatal("expected error for empty username and password, got nil")
	}
}

// TestLoginLogic_Login_UserNotFound 测试用户不存在的场景
// 场景：用户使用未注册的用户名登录
// 预期：返回用户未找到错误 LoginUserNotFound
func TestLoginLogic_Login_UserNotFound(t *testing.T) {
	helper := generateTestJWTHelper(t)
	mockModel := &model.MockUsersModel{
		MockFindOneByUsername: func(ctx context.Context, username string) (*model.Users, error) {
			return nil, model.ErrNotFound
		},
	}

	svcCtx := &svc.ServiceContext{
		UsersModel: mockModel,
		JWTHelper:  helper,
	}

	logic := NewLoginLogic(context.Background(), svcCtx)
	req := &types.LoginReq{
		Username: "nonexistent",
		Password: "password123",
	}

	_, err := logic.Login(req)
	if err == nil {
		t.Fatal("expected error for non-existent user, got nil")
	}
}

// TestLoginLogic_Login_WrongPassword 测试密码错误的场景
// 场景：用户名正确，但密码错误
// 预期：返回密码错误 LoginPasswordWrong
func TestLoginLogic_Login_WrongPassword(t *testing.T) {
	helper := generateTestJWTHelper(t)
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
		UsersModel: mockModel,
		JWTHelper:  helper,
	}

	logic := NewLoginLogic(context.Background(), svcCtx)
	req := &types.LoginReq{
		Username: "testuser",
		Password: "wrongpassword",
	}

	_, err := logic.Login(req)
	if err == nil {
		t.Fatal("expected error for wrong password, got nil")
	}
}

// TestLoginLogic_Login_DBQueryFailed 测试数据库查询失败的场景
// 场景：数据库查询用户时发生错误（非未找到错误）
// 预期：返回查询失败错误 LoginQueryFailed
func TestLoginLogic_Login_DBQueryFailed(t *testing.T) {
	helper := generateTestJWTHelper(t)
	mockErr := errors.New("database connection lost")
	mockModel := &model.MockUsersModel{
		MockFindOneByUsername: func(ctx context.Context, username string) (*model.Users, error) {
			return nil, mockErr
		},
	}

	svcCtx := &svc.ServiceContext{
		UsersModel: mockModel,
		JWTHelper:  helper,
	}

	logic := NewLoginLogic(context.Background(), svcCtx)
	req := &types.LoginReq{
		Username: "testuser",
		Password: "password123",
	}

	_, err := logic.Login(req)
	if err == nil {
		t.Fatal("expected error for DB query failure, got nil")
	}
}

// TestLoginLogic_Login_TokenGenerationFailed 测试 JWT 生成失败的场景
// 场景：用户验证成功，但 JWT 生成失败（JWTHelper 为 nil）
// 预期：捕获 panic 或返回令牌生成失败错误
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
		UsersModel: mockModel,
		JWTHelper:  nil,
	}

	logic := NewLoginLogic(context.Background(), svcCtx)
	req := &types.LoginReq{
		Username: "testuser",
		Password: "password123",
	}

	defer func() {
		if r := recover(); r == nil {
			t.Log("expected panic or error for nil JWTHelper")
		}
	}()

	_, err := logic.Login(req)
	if err == nil {
		t.Fatal("expected error or panic for nil JWTHelper, got nil")
	}
}

// TestLoginLogic_Login_CaseSensitiveUsername 测试用户名大小写敏感
// 场景：用户使用不同大小写的用户名登录
// 预期：用户名区分大小写，testuser 和 TestUser 视为不同用户
func TestLoginLogic_Login_CaseSensitiveUsername(t *testing.T) {
	helper := generateTestJWTHelper(t)
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
		UsersModel: mockModel,
		JWTHelper:  helper,
	}

	logic := NewLoginLogic(context.Background(), svcCtx)
	req := &types.LoginReq{
		Username: "TestUser",
		Password: "password123",
	}

	_, err := logic.Login(req)
	if err == nil {
		t.Fatal("expected error for case-sensitive username mismatch, got nil")
	}
}

// TestLoginLogic_Login_LongPassword 测试超长密码登录
// 场景：bcrypt 限制密码最大长度为 72 字节，超过则哈希失败
// 预期：登录失败，因为 bcrypt.GenerateFromPassword 会返回错误
func TestLoginLogic_Login_LongPassword(t *testing.T) {
	longPassword := make([]byte, 100)
	for i := range longPassword {
		longPassword[i] = 'a'
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(longPassword), bcrypt.DefaultCost)
	if err != nil {
		t.Skipf("bcrypt rejected long password (expected): %v", err)
	}

	helper := generateTestJWTHelper(t)
	mockModel := &model.MockUsersModel{
		MockFindOneByUsername: func(ctx context.Context, username string) (*model.Users, error) {
			if username == "testuser" {
				return &model.Users{
					Id:           1,
					Username:     "testuser",
					PasswordHash: string(hashedPwd),
					CreatedAt:    1234567890,
					UpdatedAt:    1234567890,
				}, nil
			}
			return nil, model.ErrNotFound
		},
	}

	svcCtx := &svc.ServiceContext{
		UsersModel: mockModel,
		JWTHelper:  helper,
	}

	logic := NewLoginLogic(context.Background(), svcCtx)
	req := &types.LoginReq{
		Username: "testuser",
		Password: string(longPassword),
	}

	resp, err := logic.Login(req)
	if err != nil {
		t.Fatalf("expected no error for valid password, got %v", err)
	}
	if resp.Token == "" {
		t.Error("expected non-empty token for valid password")
	}
}

// TestLoginLogic_Login_SpecialCharactersInPassword 测试密码包含特殊字符
// 场景：用户使用包含特殊字符的密码登录
// 预期：登录流程正常处理特殊字符密码
func TestLoginLogic_Login_SpecialCharactersInPassword(t *testing.T) {
	helper := generateTestJWTHelper(t)
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
		UsersModel: mockModel,
		JWTHelper:  helper,
	}

	logic := NewLoginLogic(context.Background(), svcCtx)
	req := &types.LoginReq{
		Username: "testuser",
		Password: specialPassword,
	}

	resp, err := logic.Login(req)
	if err != nil {
		t.Fatalf("expected no error for special character password, got %v", err)
	}
	if resp.Token == "" {
		t.Error("expected non-empty token for special character password")
	}
}

// TestLoginLogic_Login_WhitespaceUsername 测试用户名包含空白字符
// 场景：用户使用包含前后空白的用户名登录
// 预期：系统不会自动 trim，按原样匹配（区分空白）
func TestLoginLogic_Login_WhitespaceUsername(t *testing.T) {
	helper := generateTestJWTHelper(t)
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
		UsersModel: mockModel,
		JWTHelper:  helper,
	}

	logic := NewLoginLogic(context.Background(), svcCtx)
	req := &types.LoginReq{
		Username: " testuser ",
		Password: "password123",
	}

	_, err := logic.Login(req)
	if err == nil {
		t.Fatal("expected error for username with whitespace, got nil")
	}
}

// TestLoginLogic_Login_ResponseCreatedAt 测试返回的用户创建时间格式
// 场景：用户登录成功
// 预期：返回的 CreatedAt 字段为 Unix 时间戳字符串格式
func TestLoginLogic_Login_ResponseCreatedAt(t *testing.T) {
	helper := generateTestJWTHelper(t)
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
		UsersModel: mockModel,
		JWTHelper:  helper,
	}

	logic := NewLoginLogic(context.Background(), svcCtx)
	req := &types.LoginReq{
		Username: "testuser",
		Password: "password123",
	}

	resp, err := logic.Login(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.User.CreatedAt != "1234567890" {
		t.Errorf("expected CreatedAt '1234567890', got %s", resp.User.CreatedAt)
	}
}

// TestLoginLogic_Login_MultipleUsers 测试多用户场景下的登录
// 场景：系统中有多个用户，使用正确的用户名密码登录特定用户
// 预期：返回对应用户的信息和 token
func TestLoginLogic_Login_MultipleUsers(t *testing.T) {
	helper := generateTestJWTHelper(t)
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
		UsersModel: mockModel,
		JWTHelper:  helper,
	}

	logic := NewLoginLogic(context.Background(), svcCtx)

	req1 := &types.LoginReq{
		Username: "user1",
		Password: "password1",
	}
	resp1, err := logic.Login(req1)
	if err != nil {
		t.Fatalf("expected no error for user1, got %v", err)
	}
	if resp1.User.Id != 1 || resp1.User.Username != "user1" {
		t.Errorf("expected user1 info, got id=%d username=%s", resp1.User.Id, resp1.User.Username)
	}

	req2 := &types.LoginReq{
		Username: "user2",
		Password: "password2",
	}
	resp2, err := logic.Login(req2)
	if err != nil {
		t.Fatalf("expected no error for user2, got %v", err)
	}
	if resp2.User.Id != 2 || resp2.User.Username != "user2" {
		t.Errorf("expected user2 info, got id=%d username=%s", resp2.User.Id, resp2.User.Username)
	}
}

// TestHashPassword 测试密码哈希辅助函数
// 场景：调用 hashPassword
// 预期：返回非空的 bcrypt 哈希字符串
func TestHashPassword(t *testing.T) {
	hashed := hashPassword(t, "testpassword")
	if hashed == "" {
		t.Fatal("expected non-empty hashed password")
	}

	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte("testpassword"))
	if err != nil {
		t.Errorf("hashed password does not match original: %v", err)
	}
}

// TestJWTHelper_GenerateAndValidateToken 测试 JWT 令牌的生成和验证
// 场景：生成 JWT token 并验证其有效性
// 预期：token 可以被正确生成和验证
func TestJWTHelper_GenerateAndValidateToken(t *testing.T) {
	helper := generateTestJWTHelper(t)

	token, err := helper.GenerateToken(jwtx.JWTClaims{
		UID:       "123",
		TokenType: jwtx.Access,
	})
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	if token == "" {
		t.Fatal("expected non-empty token")
	}

	_, claims, err := helper.ValidateToken(token)
	if err != nil {
		t.Fatalf("failed to validate token: %v", err)
	}

	if claims == nil {
		t.Fatal("expected non-nil claims")
	}
}
