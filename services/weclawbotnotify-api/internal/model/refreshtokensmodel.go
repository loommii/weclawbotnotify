package model

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ RefreshTokensModel = (*customRefreshTokensModel)(nil)

type (
	// RefreshTokensModel Refresh Token 数据访问接口
	RefreshTokensModel interface {
		refreshTokensModel
		withSession(session sqlx.Session) RefreshTokensModel
		RevokeByJti(ctx context.Context, jti string) error      // 按 JTI 撤销（令牌轮换）
		RevokeByUserId(ctx context.Context, userId int64) error // 按用户撤销（修改密码/强制下线）
		DeleteExpired(ctx context.Context) error                // 清理过期令牌
	}

	customRefreshTokensModel struct {
		*defaultRefreshTokensModel
	}
)

func NewRefreshTokensModel(conn sqlx.SqlConn) RefreshTokensModel {
	return &customRefreshTokensModel{
		defaultRefreshTokensModel: newRefreshTokensModel(conn),
	}
}

func (m *customRefreshTokensModel) withSession(session sqlx.Session) RefreshTokensModel {
	return NewRefreshTokensModel(sqlx.NewSqlConnFromSession(session))
}

// RevokeByJti 按 JTI 撤销 Refresh Token（令牌轮换时使用）
func (m *customRefreshTokensModel) RevokeByJti(ctx context.Context, jti string) error {
	query := fmt.Sprintf("update %s set `revoked` = 1 where `jti` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, jti)
	return err
}

// RevokeByUserId 按用户 ID 撤销所有 Refresh Token（修改密码/强制下线时使用）
func (m *customRefreshTokensModel) RevokeByUserId(ctx context.Context, userId int64) error {
	query := fmt.Sprintf("update %s set `revoked` = 1 where `user_id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, userId)
	return err
}

// DeleteExpired 清理过期的 Refresh Token（兼容 SQLite，应用层计算时间）
func (m *customRefreshTokensModel) DeleteExpired(ctx context.Context) error {
	now := time.Now().Unix()
	query := fmt.Sprintf("delete from %s where `expires_at` < ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, now)
	return err
}
