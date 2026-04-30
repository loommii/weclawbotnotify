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

// generateTestJWTHelper 生成测试用的 JWTHelper
func generateTestJWTHelper(t *testing.T) *jwtx.JWTHelper {
	t.Helper()
	privateKey, err := rsa.GenerateKey(nil, 2048)
	if err != nil {
		t.Fatalf("failed to generate RSA private key: %v", err)
	}
	return jwtx.NewJWTHelper(
		jwtx.WithPrivateKey(privateKey),
		jwtx.WithPublicKey(&privateKey.PublicKey),
		jwtx.WithExpiredTime(1*time.Hour),
	)
}

// TestRegisterLogic_Register_Success 测试正常注册流程
// 场景：首次注册用户，系统中无用户，用户名不重复，所有操作正常
// 预期：注册成功，返回有效的 JWT token 和用户信息
func TestRegisterLogic_Register_Success(t *testing.T) {
	helper := generateTestJWTHelper(t)
	mockModel := &model.MockUsersModel{
		MockCount: func(ctx context.Context) (int64, error) {
			return 0, nil
		},
		MockFindOneByUsername: func(ctx context.Context, username string) (*model.Users, error) {
			return nil, model.ErrNotFound
		},
		MockInsert: func(ctx context.Context, data *model.Users) (sql.Result, error) {
			return &model.MockResult{LastInsertIdVal: 1}, nil
		},
	}

	svcCtx := &svc.ServiceContext{
		UsersModel: mockModel,
		JWTHelper:  helper,
	}

	logic := NewRegisterLogic(context.Background(), svcCtx)
	req := &types.RegisterReq{
		Username: "testuser",
		Password: "password123",
	}

	resp, err := logic.Register(req)
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

// TestRegisterLogic_Register_EmptyUsername 测试用户名为空的参数校验
// 场景：用户注册时用户名为空字符串
// 预期：返回 RegisterParamEmpty 错误，响应为 nil
func TestRegisterLogic_Register_EmptyUsername(t *testing.T) {
	helper := generateTestJWTHelper(t)
	mockModel := &model.MockUsersModel{}

	svcCtx := &svc.ServiceContext{
		UsersModel: mockModel,
		JWTHelper:  helper,
	}

	logic := NewRegisterLogic(context.Background(), svcCtx)
	req := &types.RegisterReq{
		Username: "",
		Password: "password123",
	}

	resp, err := logic.Register(req)
	if err == nil {
		t.Fatal("expected error for empty username")
	}
	if resp != nil {
		t.Error("expected nil response")
	}
	if !errors.Is(err, xerr.RegisterParamEmpty) {
		t.Errorf("expected RegisterParamEmpty error, got %v", err)
	}
}

// TestRegisterLogic_Register_EmptyPassword 测试密码为空的参数校验
// 场景：用户注册时密码为空字符串
// 预期：返回 RegisterParamEmpty 错误，响应为 nil
func TestRegisterLogic_Register_EmptyPassword(t *testing.T) {
	helper := generateTestJWTHelper(t)
	mockModel := &model.MockUsersModel{}

	svcCtx := &svc.ServiceContext{
		UsersModel: mockModel,
		JWTHelper:  helper,
	}

	logic := NewRegisterLogic(context.Background(), svcCtx)
	req := &types.RegisterReq{
		Username: "testuser",
		Password: "",
	}

	resp, err := logic.Register(req)
	if err == nil {
		t.Fatal("expected error for empty password")
	}
	if resp != nil {
		t.Error("expected nil response")
	}
	if !errors.Is(err, xerr.RegisterParamEmpty) {
		t.Errorf("expected RegisterParamEmpty error, got %v", err)
	}
}

// TestRegisterLogic_Register_AlreadyExists 测试 V1 版本单用户限制
// 场景：系统中已存在用户（Count > 0），再次尝试注册新用户
// 预期：返回 RegisterClosed 错误，拒绝注册
func TestRegisterLogic_Register_AlreadyExists(t *testing.T) {
	helper := generateTestJWTHelper(t)
	mockModel := &model.MockUsersModel{
		MockCount: func(ctx context.Context) (int64, error) {
			return 1, nil
		},
	}

	svcCtx := &svc.ServiceContext{
		UsersModel: mockModel,
		JWTHelper:  helper,
	}

	logic := NewRegisterLogic(context.Background(), svcCtx)
	req := &types.RegisterReq{
		Username: "testuser",
		Password: "password123",
	}

	resp, err := logic.Register(req)
	if err == nil {
		t.Fatal("expected error when user already exists")
	}
	if resp != nil {
		t.Error("expected nil response")
	}
	if !errors.Is(err, xerr.RegisterClosed) {
		t.Errorf("expected RegisterClosed error, got %v", err)
	}
}

// TestRegisterLogic_Register_UsernameAlreadyExists 测试用户名重复检测
// 场景：系统中无用户，但尝试注册的用户名已存在（FindOneByUsername 返回用户）
// 预期：返回 RegisterUsernameExist 错误，响应为 nil
func TestRegisterLogic_Register_UsernameAlreadyExists(t *testing.T) {
	helper := generateTestJWTHelper(t)
	mockModel := &model.MockUsersModel{
		MockCount: func(ctx context.Context) (int64, error) {
			return 0, nil
		},
		MockFindOneByUsername: func(ctx context.Context, username string) (*model.Users, error) {
			return &model.Users{Id: 1, Username: username}, nil
		},
	}

	svcCtx := &svc.ServiceContext{
		UsersModel: mockModel,
		JWTHelper:  helper,
	}

	logic := NewRegisterLogic(context.Background(), svcCtx)
	req := &types.RegisterReq{
		Username: "testuser",
		Password: "password123",
	}

	resp, err := logic.Register(req)
	if err == nil {
		t.Fatal("expected error when username already exists")
	}
	if resp != nil {
		t.Error("expected nil response")
	}
	if !errors.Is(err, xerr.RegisterUsernameExist) {
		t.Errorf("expected RegisterUsernameExist error, got %v", err)
	}
}

// TestRegisterLogic_Register_CountError 测试查询用户数量失败
// 场景：数据库查询用户数量时发生错误
// 预期：返回 RegisterQueryFailed 错误，响应为 nil
func TestRegisterLogic_Register_CountError(t *testing.T) {
	helper := generateTestJWTHelper(t)
	mockModel := &model.MockUsersModel{
		MockCount: func(ctx context.Context) (int64, error) {
			return 0, errors.New("db error")
		},
	}

	svcCtx := &svc.ServiceContext{
		UsersModel: mockModel,
		JWTHelper:  helper,
	}

	logic := NewRegisterLogic(context.Background(), svcCtx)
	req := &types.RegisterReq{
		Username: "testuser",
		Password: "password123",
	}

	resp, err := logic.Register(req)
	if err == nil {
		t.Fatal("expected error when count fails")
	}
	if resp != nil {
		t.Error("expected nil response")
	}
	if !errors.Is(err, xerr.RegisterQueryFailed) {
		t.Errorf("expected RegisterQueryFailed error, got %v", err)
	}
}

// TestRegisterLogic_Register_FindOneError 测试查询用户名时发生非 NotFound 错误
// 场景：查询用户名是否存在时，数据库返回非预期错误（非 ErrNotFound）
// 预期：返回 RegisterQueryFailed 错误，响应为 nil
func TestRegisterLogic_Register_FindOneError(t *testing.T) {
	helper := generateTestJWTHelper(t)
	mockModel := &model.MockUsersModel{
		MockCount: func(ctx context.Context) (int64, error) {
			return 0, nil
		},
		MockFindOneByUsername: func(ctx context.Context, username string) (*model.Users, error) {
			return nil, errors.New("db error")
		},
	}

	svcCtx := &svc.ServiceContext{
		UsersModel: mockModel,
		JWTHelper:  helper,
	}

	logic := NewRegisterLogic(context.Background(), svcCtx)
	req := &types.RegisterReq{
		Username: "testuser",
		Password: "password123",
	}

	resp, err := logic.Register(req)
	if err == nil {
		t.Fatal("expected error when FindOneByUsername fails")
	}
	if resp != nil {
		t.Error("expected nil response")
	}
	if !errors.Is(err, xerr.RegisterQueryFailed) {
		t.Errorf("expected RegisterQueryFailed error, got %v", err)
	}
}

// TestRegisterLogic_Register_InsertError 测试插入用户记录失败
// 场景：参数校验通过、用户名不重复，但插入用户时数据库报错
// 预期：返回 RegisterInsertFailed 错误，响应为 nil
func TestRegisterLogic_Register_InsertError(t *testing.T) {
	helper := generateTestJWTHelper(t)
	mockModel := &model.MockUsersModel{
		MockCount: func(ctx context.Context) (int64, error) {
			return 0, nil
		},
		MockFindOneByUsername: func(ctx context.Context, username string) (*model.Users, error) {
			return nil, model.ErrNotFound
		},
		MockInsert: func(ctx context.Context, data *model.Users) (sql.Result, error) {
			return nil, errors.New("db error")
		},
	}

	svcCtx := &svc.ServiceContext{
		UsersModel: mockModel,
		JWTHelper:  helper,
	}

	logic := NewRegisterLogic(context.Background(), svcCtx)
	req := &types.RegisterReq{
		Username: "testuser",
		Password: "password123",
	}

	resp, err := logic.Register(req)
	if err == nil {
		t.Fatal("expected error when insert fails")
	}
	if resp != nil {
		t.Error("expected nil response")
	}
	if !errors.Is(err, xerr.RegisterInsertFailed) {
		t.Errorf("expected RegisterInsertFailed error, got %v", err)
	}
}

// TestRegisterLogic_Register_GetIdError 测试获取自增 ID 失败
// 场景：用户插入成功，但获取 LastInsertId 时发生错误
// 预期：返回 RegisterGetIdFailed 错误，响应为 nil
func TestRegisterLogic_Register_GetIdError(t *testing.T) {
	helper := generateTestJWTHelper(t)
	mockModel := &model.MockUsersModel{
		MockCount: func(ctx context.Context) (int64, error) {
			return 0, nil
		},
		MockFindOneByUsername: func(ctx context.Context, username string) (*model.Users, error) {
			return nil, model.ErrNotFound
		},
		MockInsert: func(ctx context.Context, data *model.Users) (sql.Result, error) {
			return &model.MockResult{LastInsertIdErr: errors.New("db error")}, nil
		},
	}

	svcCtx := &svc.ServiceContext{
		UsersModel: mockModel,
		JWTHelper:  helper,
	}

	logic := NewRegisterLogic(context.Background(), svcCtx)
	req := &types.RegisterReq{
		Username: "testuser",
		Password: "password123",
	}

	resp, err := logic.Register(req)
	if err == nil {
		t.Fatal("expected error when get id fails")
	}
	if resp != nil {
		t.Error("expected nil response")
	}
	if !errors.Is(err, xerr.RegisterGetIdFailed) {
		t.Errorf("expected RegisterGetIdFailed error, got %v", err)
	}
}

// TestRegisterLogic_Register_JWTError 测试 JWT 生成失败（JWTHelper 为 nil）
// 场景：用户插入成功，但 JWTHelper 未初始化（nil），调用 GenerateToken 时 panic
// 预期：触发 panic，验证代码对空指针的防御机制
func TestRegisterLogic_Register_JWTError(t *testing.T) {
	mockModel := &model.MockUsersModel{
		MockCount: func(ctx context.Context) (int64, error) {
			return 0, nil
		},
		MockFindOneByUsername: func(ctx context.Context, username string) (*model.Users, error) {
			return nil, model.ErrNotFound
		},
		MockInsert: func(ctx context.Context, data *model.Users) (sql.Result, error) {
			return &model.MockResult{LastInsertIdVal: 1}, nil
		},
	}

	svcCtx := &svc.ServiceContext{
		UsersModel: mockModel,
		JWTHelper:  nil,
	}

	logic := NewRegisterLogic(context.Background(), svcCtx)
	req := &types.RegisterReq{
		Username: "testuser",
		Password: "password123",
	}

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic when JWTHelper is nil")
		}
	}()
	logic.Register(req)
}

// TestIsNotFound 测试 isNotFound 辅助函数的判断逻辑
// 场景 1：传入 model.ErrNotFound，预期返回 true
// 场景 2：传入普通错误，预期返回 false
// 场景 3：传入 nil，预期返回 false
func TestIsNotFound(t *testing.T) {
	if !isNotFound(model.ErrNotFound) {
		t.Error("expected isNotFound to return true for ErrNotFound")
	}
	if isNotFound(errors.New("some other error")) {
		t.Error("expected isNotFound to return false for other errors")
	}
	if isNotFound(nil) {
		t.Error("expected isNotFound to return false for nil")
	}
}
