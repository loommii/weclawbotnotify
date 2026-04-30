package model

import (
	"context"
	"database/sql"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// MockRefreshTokensModel RefreshTokensModel 的 Mock 实现
type MockRefreshTokensModel struct {
	MockInsert         func(ctx context.Context, data *RefreshTokens) (sql.Result, error)
	MockFindOne        func(ctx context.Context, id int64) (*RefreshTokens, error)
	MockFindOneByJti   func(ctx context.Context, jti string) (*RefreshTokens, error)
	MockUpdate         func(ctx context.Context, data *RefreshTokens) error
	MockDelete         func(ctx context.Context, id int64) error
	MockRevokeByJti    func(ctx context.Context, jti string) error
	MockRevokeByUserId func(ctx context.Context, userId int64) error
	MockDeleteExpired  func(ctx context.Context) error
}

func (m *MockRefreshTokensModel) Insert(ctx context.Context, data *RefreshTokens) (sql.Result, error) {
	if m.MockInsert != nil {
		return m.MockInsert(ctx, data)
	}
	return &MockResult{}, nil
}

func (m *MockRefreshTokensModel) FindOne(ctx context.Context, id int64) (*RefreshTokens, error) {
	if m.MockFindOne != nil {
		return m.MockFindOne(ctx, id)
	}
	return nil, ErrNotFound
}

func (m *MockRefreshTokensModel) FindOneByJti(ctx context.Context, jti string) (*RefreshTokens, error) {
	if m.MockFindOneByJti != nil {
		return m.MockFindOneByJti(ctx, jti)
	}
	return nil, ErrNotFound
}

func (m *MockRefreshTokensModel) Update(ctx context.Context, data *RefreshTokens) error {
	if m.MockUpdate != nil {
		return m.MockUpdate(ctx, data)
	}
	return nil
}

func (m *MockRefreshTokensModel) Delete(ctx context.Context, id int64) error {
	if m.MockDelete != nil {
		return m.MockDelete(ctx, id)
	}
	return nil
}

// RevokeByJti 按 JTI 撤销（令牌轮换）
func (m *MockRefreshTokensModel) RevokeByJti(ctx context.Context, jti string) error {
	if m.MockRevokeByJti != nil {
		return m.MockRevokeByJti(ctx, jti)
	}
	return nil
}

// RevokeByUserId 按用户 ID 撤销（修改密码/强制下线）
func (m *MockRefreshTokensModel) RevokeByUserId(ctx context.Context, userId int64) error {
	if m.MockRevokeByUserId != nil {
		return m.MockRevokeByUserId(ctx, userId)
	}
	return nil
}

// DeleteExpired 清理过期令牌
func (m *MockRefreshTokensModel) DeleteExpired(ctx context.Context) error {
	if m.MockDeleteExpired != nil {
		return m.MockDeleteExpired(ctx)
	}
	return nil
}

func (m *MockRefreshTokensModel) withSession(session sqlx.Session) RefreshTokensModel {
	return m
}
