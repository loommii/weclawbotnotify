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

// generateRegisterTestJWTHelpers 生成测试用的 Access 和 Refresh JWTHelper
func generateRegisterTestJWTHelpers(t *testing.T) (*jwtx.JWTHelper, *jwtx.JWTHelper) {
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

// newTestRefreshMock 创建默认的 MockRefreshTokensModel
func newTestRefreshMock() *model.MockRefreshTokensModel {
	return &model.MockRefreshTokensModel{
		MockInsert: func(ctx context.Context, data *model.RefreshTokens) (sql.Result, error) {
			return &model.MockResult{LastInsertIdVal: 1}, nil
		},
	}
}

// TestRegisterLogic_Register_Success 测试正常注册流程
// 场景：首次注册用户，系统中无用户，用户名不重复，所有操作正常
// 预期：注册成功，返回有效的双令牌和用户信息
func TestRegisterLogic_Register_Success(t *testing.T) {
	accessHelper, refreshHelper := generateRegisterTestJWTHelpers(t)
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
		UsersModel:         mockModel,
		RefreshTokensModel: newTestRefreshMock(),
		AccessJWTHelper:    accessHelper,
		RefreshJWTHelper:   refreshHelper,
	}

	logic := NewRegisterLogic(context.Background(), svcCtx)
	req := &types.RegisterReq{
		Username: "testuser",
		Password: "password123",
	}

	resp, err := logic.Register(req)
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

// TestRegisterLogic_Register_EmptyUsername 测试用户名为空的参数校验
// 场景：用户注册时用户名为空字符串
// 预期：返回 RegisterParamEmpty 错误
func TestRegisterLogic_Register_EmptyUsername(t *testing.T) {
	accessHelper, refreshHelper := generateRegisterTestJWTHelpers(t)
	mockModel := &model.MockUsersModel{}

	svcCtx := &svc.ServiceContext{
		UsersModel:         mockModel,
		RefreshTokensModel: newTestRefreshMock(),
		AccessJWTHelper:    accessHelper,
		RefreshJWTHelper:   refreshHelper,
	}

	logic := NewRegisterLogic(context.Background(), svcCtx)
	req := &types.RegisterReq{
		Username: "",
		Password: "password123",
	}

	resp, err := logic.Register(req)
	if err == nil {
		t.Fatal("空用户名应返回错误")
	}
	if resp != nil {
		t.Error("预期 nil 响应")
	}
	if !errors.Is(err, xerr.RegisterParamEmpty) {
		t.Errorf("预期 RegisterParamEmpty, got %v", err)
	}
}

// TestRegisterLogic_Register_EmptyPassword 测试密码为空的参数校验
// 场景：用户注册时密码为空字符串
// 预期：返回 RegisterParamEmpty 错误
func TestRegisterLogic_Register_EmptyPassword(t *testing.T) {
	accessHelper, refreshHelper := generateRegisterTestJWTHelpers(t)
	mockModel := &model.MockUsersModel{}

	svcCtx := &svc.ServiceContext{
		UsersModel:         mockModel,
		RefreshTokensModel: newTestRefreshMock(),
		AccessJWTHelper:    accessHelper,
		RefreshJWTHelper:   refreshHelper,
	}

	logic := NewRegisterLogic(context.Background(), svcCtx)
	req := &types.RegisterReq{
		Username: "testuser",
		Password: "",
	}

	resp, err := logic.Register(req)
	if err == nil {
		t.Fatal("空密码应返回错误")
	}
	if resp != nil {
		t.Error("预期 nil 响应")
	}
	if !errors.Is(err, xerr.RegisterParamEmpty) {
		t.Errorf("预期 RegisterParamEmpty, got %v", err)
	}
}

// TestRegisterLogic_Register_AlreadyExists 测试 V1 版本单用户限制
// 场景：系统中已存在用户（Count > 0），再次尝试注册新用户
// 预期：返回 RegisterClosed 错误
func TestRegisterLogic_Register_AlreadyExists(t *testing.T) {
	accessHelper, refreshHelper := generateRegisterTestJWTHelpers(t)
	mockModel := &model.MockUsersModel{
		MockCount: func(ctx context.Context) (int64, error) {
			return 1, nil
		},
	}

	svcCtx := &svc.ServiceContext{
		UsersModel:         mockModel,
		RefreshTokensModel: newTestRefreshMock(),
		AccessJWTHelper:    accessHelper,
		RefreshJWTHelper:   refreshHelper,
	}

	logic := NewRegisterLogic(context.Background(), svcCtx)
	req := &types.RegisterReq{
		Username: "testuser",
		Password: "password123",
	}

	resp, err := logic.Register(req)
	if err == nil {
		t.Fatal("已存在用户应返回错误")
	}
	if resp != nil {
		t.Error("预期 nil 响应")
	}
	if !errors.Is(err, xerr.RegisterClosed) {
		t.Errorf("预期 RegisterClosed, got %v", err)
	}
}

// TestRegisterLogic_Register_UsernameAlreadyExists 测试用户名重复检测
// 场景：系统中无用户，但尝试注册的用户名已存在
// 预期：返回 RegisterUsernameExist 错误
func TestRegisterLogic_Register_UsernameAlreadyExists(t *testing.T) {
	accessHelper, refreshHelper := generateRegisterTestJWTHelpers(t)
	mockModel := &model.MockUsersModel{
		MockCount: func(ctx context.Context) (int64, error) {
			return 0, nil
		},
		MockFindOneByUsername: func(ctx context.Context, username string) (*model.Users, error) {
			return &model.Users{Id: 1, Username: username}, nil
		},
	}

	svcCtx := &svc.ServiceContext{
		UsersModel:         mockModel,
		RefreshTokensModel: newTestRefreshMock(),
		AccessJWTHelper:    accessHelper,
		RefreshJWTHelper:   refreshHelper,
	}

	logic := NewRegisterLogic(context.Background(), svcCtx)
	req := &types.RegisterReq{
		Username: "testuser",
		Password: "password123",
	}

	resp, err := logic.Register(req)
	if err == nil {
		t.Fatal("用户名重复应返回错误")
	}
	if resp != nil {
		t.Error("预期 nil 响应")
	}
	if !errors.Is(err, xerr.RegisterUsernameExist) {
		t.Errorf("预期 RegisterUsernameExist, got %v", err)
	}
}

// TestRegisterLogic_Register_CountError 测试查询用户数量失败
// 场景：数据库查询用户数量时发生错误
// 预期：返回 RegisterQueryFailed 错误
func TestRegisterLogic_Register_CountError(t *testing.T) {
	accessHelper, refreshHelper := generateRegisterTestJWTHelpers(t)
	mockModel := &model.MockUsersModel{
		MockCount: func(ctx context.Context) (int64, error) {
			return 0, errors.New("db error")
		},
	}

	svcCtx := &svc.ServiceContext{
		UsersModel:         mockModel,
		RefreshTokensModel: newTestRefreshMock(),
		AccessJWTHelper:    accessHelper,
		RefreshJWTHelper:   refreshHelper,
	}

	logic := NewRegisterLogic(context.Background(), svcCtx)
	req := &types.RegisterReq{
		Username: "testuser",
		Password: "password123",
	}

	resp, err := logic.Register(req)
	if err == nil {
		t.Fatal("查询失败应返回错误")
	}
	if resp != nil {
		t.Error("预期 nil 响应")
	}
	if !errors.Is(err, xerr.RegisterQueryFailed) {
		t.Errorf("预期 RegisterQueryFailed, got %v", err)
	}
}

// TestRegisterLogic_Register_FindOneError 测试查询用户名时发生非 NotFound 错误
// 场景：查询用户名是否存在时，数据库返回非预期错误
// 预期：返回 RegisterQueryFailed 错误
func TestRegisterLogic_Register_FindOneError(t *testing.T) {
	accessHelper, refreshHelper := generateRegisterTestJWTHelpers(t)
	mockModel := &model.MockUsersModel{
		MockCount: func(ctx context.Context) (int64, error) {
			return 0, nil
		},
		MockFindOneByUsername: func(ctx context.Context, username string) (*model.Users, error) {
			return nil, errors.New("db error")
		},
	}

	svcCtx := &svc.ServiceContext{
		UsersModel:         mockModel,
		RefreshTokensModel: newTestRefreshMock(),
		AccessJWTHelper:    accessHelper,
		RefreshJWTHelper:   refreshHelper,
	}

	logic := NewRegisterLogic(context.Background(), svcCtx)
	req := &types.RegisterReq{
		Username: "testuser",
		Password: "password123",
	}

	resp, err := logic.Register(req)
	if err == nil {
		t.Fatal("查询失败应返回错误")
	}
	if resp != nil {
		t.Error("预期 nil 响应")
	}
	if !errors.Is(err, xerr.RegisterQueryFailed) {
		t.Errorf("预期 RegisterQueryFailed, got %v", err)
	}
}

// TestRegisterLogic_Register_InsertError 测试插入用户记录失败
// 场景：参数校验通过、用户名不重复，但插入用户时数据库报错
// 预期：返回 RegisterInsertFailed 错误
func TestRegisterLogic_Register_InsertError(t *testing.T) {
	accessHelper, refreshHelper := generateRegisterTestJWTHelpers(t)
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
		UsersModel:         mockModel,
		RefreshTokensModel: newTestRefreshMock(),
		AccessJWTHelper:    accessHelper,
		RefreshJWTHelper:   refreshHelper,
	}

	logic := NewRegisterLogic(context.Background(), svcCtx)
	req := &types.RegisterReq{
		Username: "testuser",
		Password: "password123",
	}

	resp, err := logic.Register(req)
	if err == nil {
		t.Fatal("插入失败应返回错误")
	}
	if resp != nil {
		t.Error("预期 nil 响应")
	}
	if !errors.Is(err, xerr.RegisterInsertFailed) {
		t.Errorf("预期 RegisterInsertFailed, got %v", err)
	}
}

// TestRegisterLogic_Register_GetIdError 测试获取自增 ID 失败
// 场景：用户插入成功，但获取 LastInsertId 时发生错误
// 预期：返回 RegisterGetIdFailed 错误
func TestRegisterLogic_Register_GetIdError(t *testing.T) {
	accessHelper, refreshHelper := generateRegisterTestJWTHelpers(t)
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
		UsersModel:         mockModel,
		RefreshTokensModel: newTestRefreshMock(),
		AccessJWTHelper:    accessHelper,
		RefreshJWTHelper:   refreshHelper,
	}

	logic := NewRegisterLogic(context.Background(), svcCtx)
	req := &types.RegisterReq{
		Username: "testuser",
		Password: "password123",
	}

	resp, err := logic.Register(req)
	if err == nil {
		t.Fatal("获取 ID 失败应返回错误")
	}
	if resp != nil {
		t.Error("预期 nil 响应")
	}
	if !errors.Is(err, xerr.RegisterGetIdFailed) {
		t.Errorf("预期 RegisterGetIdFailed, got %v", err)
	}
}

// TestRegisterLogic_Register_JWTError 测试 JWT 生成失败（JWTHelper 为 nil）
// 场景：用户插入成功，但 JWTHelper 未初始化（nil），调用 GenerateToken 时 panic
// 预期：触发 panic
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
		UsersModel:         mockModel,
		RefreshTokensModel: newTestRefreshMock(),
		AccessJWTHelper:    nil,
		RefreshJWTHelper:   nil,
	}

	logic := NewRegisterLogic(context.Background(), svcCtx)
	req := &types.RegisterReq{
		Username: "testuser",
		Password: "password123",
	}

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("JWTHelper 为 nil 时应触发 panic")
		}
	}()
	logic.Register(req)
}

// TestIsNotFound 测试 isNotFound 辅助函数
func TestIsNotFound(t *testing.T) {
	if !isNotFound(model.ErrNotFound) {
		t.Error("ErrNotFound 应返回 true")
	}
	if isNotFound(errors.New("some other error")) {
		t.Error("其他错误应返回 false")
	}
	if isNotFound(nil) {
		t.Error("nil 应返回 false")
	}
}
