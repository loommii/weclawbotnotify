package model

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "modernc.org/sqlite"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// newRefreshTestDB 创建测试用 SQLite 数据库
func newRefreshTestDB(t *testing.T) (sqlx.SqlConn, func()) {
	t.Helper()
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")
	conn := sqlx.NewSqlConn("sqlite", dbPath)
	_, err := conn.ExecCtx(context.Background(), `
		CREATE TABLE IF NOT EXISTS refresh_tokens (
			id            INTEGER      PRIMARY KEY AUTOINCREMENT,
			user_id       INT          NOT NULL,
			jti           VARCHAR(36)  NOT NULL UNIQUE,
			revoked       INT          NOT NULL DEFAULT 0,
			expires_at    INT          NOT NULL,
			created_at    INT          NOT NULL DEFAULT (strftime('%s', 'now'))
		)
	`)
	if err != nil {
		t.Fatalf("创建表失败: %v", err)
	}
	return conn, func() { os.Remove(dbPath) }
}

// TestRefreshTokensModel_InsertAndFind 测试插入和按 JTI 查询
func TestRefreshTokensModel_InsertAndFind(t *testing.T) {
	conn, _ := newRefreshTestDB(t)
	m := NewRefreshTokensModel(conn)

	token := &RefreshTokens{
		UserId:    1,
		Jti:       "test-jti-001",
		Revoked:   0,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
	}

	result, err := m.Insert(context.Background(), token)
	if err != nil {
		t.Fatalf("Insert 失败: %v", err)
	}

	id, _ := result.LastInsertId()
	if id != 1 {
		t.Errorf("预期 id=1, got %d", id)
	}

	found, err := m.FindOneByJti(context.Background(), "test-jti-001")
	if err != nil {
		t.Fatalf("FindOneByJti 失败: %v", err)
	}
	if found.Jti != "test-jti-001" {
		t.Errorf("JTI = %s, want test-jti-001", found.Jti)
	}
	if found.Revoked != 0 {
		t.Errorf("Revoked = %d, want 0", found.Revoked)
	}
	if found.UserId != 1 {
		t.Errorf("UserId = %d, want 1", found.UserId)
	}
}

// TestRefreshTokensModel_FindOneByJti_NotFound 测试 JTI 不存在
func TestRefreshTokensModel_FindOneByJti_NotFound(t *testing.T) {
	conn, _ := newRefreshTestDB(t)
	m := NewRefreshTokensModel(conn)

	_, err := m.FindOneByJti(context.Background(), "nonexistent-jti")
	if err != ErrNotFound {
		t.Errorf("预期 ErrNotFound, got %v", err)
	}
}

// TestRefreshTokensModel_RevokeByJti 测试按 JTI 撤销（令牌轮换）
func TestRefreshTokensModel_RevokeByJti(t *testing.T) {
	conn, _ := newRefreshTestDB(t)
	m := NewRefreshTokensModel(conn)

	token := &RefreshTokens{
		UserId:    1,
		Jti:       "jti-to-revoke",
		Revoked:   0,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	m.Insert(context.Background(), token)

	err := m.RevokeByJti(context.Background(), "jti-to-revoke")
	if err != nil {
		t.Fatalf("RevokeByJti 失败: %v", err)
	}

	found, err := m.FindOneByJti(context.Background(), "jti-to-revoke")
	if err != nil {
		t.Fatalf("FindOneByJti 失败: %v", err)
	}
	if found.Revoked != 1 {
		t.Errorf("撤销后 Revoked = %d, want 1", found.Revoked)
	}
}

// TestRefreshTokensModel_RevokeByUserId 测试按用户 ID 撤销（修改密码/强制下线）
func TestRefreshTokensModel_RevokeByUserId(t *testing.T) {
	conn, _ := newRefreshTestDB(t)
	m := NewRefreshTokensModel(conn)

	// 插入同一用户的多个 Refresh Token
	for i := 0; i < 3; i++ {
		m.Insert(context.Background(), &RefreshTokens{
			UserId:    1,
			Jti:       fmt.Sprintf("jti-user1-%d", i),
			Revoked:   0,
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
		})
	}
	// 插入另一个用户的 Refresh Token
	m.Insert(context.Background(), &RefreshTokens{
		UserId:    2,
		Jti:       "jti-user2-0",
		Revoked:   0,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
	})

	err := m.RevokeByUserId(context.Background(), 1)
	if err != nil {
		t.Fatalf("RevokeByUserId 失败: %v", err)
	}

	// 用户 1 的所有令牌应被撤销
	found, err := m.FindOneByJti(context.Background(), "jti-user1-0")
	if err != nil {
		t.Fatalf("FindOneByJti 失败: %v", err)
	}
	if found.Revoked != 1 {
		t.Errorf("用户 1 的令牌 Revoked = %d, want 1", found.Revoked)
	}

	// 用户 2 的令牌应保持有效
	found2, err := m.FindOneByJti(context.Background(), "jti-user2-0")
	if err != nil {
		t.Fatalf("用户 2 的令牌应保留: %v", err)
	}
	if found2.Revoked != 0 {
		t.Errorf("用户 2 的令牌 Revoked = %d, want 0", found2.Revoked)
	}
	if found2.UserId != 2 {
		t.Errorf("保留的令牌 UserId = %d, want 2", found2.UserId)
	}
}

// TestRefreshTokensModel_DeleteExpired 测试清理过期令牌
func TestRefreshTokensModel_DeleteExpired(t *testing.T) {
	conn, _ := newRefreshTestDB(t)
	m := NewRefreshTokensModel(conn)

	// 插入已过期的令牌
	m.Insert(context.Background(), &RefreshTokens{
		UserId:    1,
		Jti:       "expired-jti",
		Revoked:   0,
		ExpiresAt: time.Now().Add(-1 * time.Hour).Unix(),
	})

	// 插入未过期的令牌
	m.Insert(context.Background(), &RefreshTokens{
		UserId:    1,
		Jti:       "valid-jti",
		Revoked:   0,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
	})

	err := m.DeleteExpired(context.Background())
	if err != nil {
		t.Fatalf("DeleteExpired 失败: %v", err)
	}

	// 过期令牌应被删除
	_, err = m.FindOneByJti(context.Background(), "expired-jti")
	if err != ErrNotFound {
		t.Errorf("过期令牌应被删除")
	}

	// 未过期令牌应保留
	_, err = m.FindOneByJti(context.Background(), "valid-jti")
	if err != nil {
		t.Fatalf("未过期令牌应保留: %v", err)
	}
}

// TestRefreshTokensModel_JTIUnique 测试 JTI 唯一约束
func TestRefreshTokensModel_JTIUnique(t *testing.T) {
	conn, _ := newRefreshTestDB(t)
	m := NewRefreshTokensModel(conn)

	token1 := &RefreshTokens{
		UserId:    1,
		Jti:       "same-jti",
		Revoked:   0,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	_, err := m.Insert(context.Background(), token1)
	if err != nil {
		t.Fatalf("第一次 Insert 失败: %v", err)
	}

	token2 := &RefreshTokens{
		UserId:    1,
		Jti:       "same-jti",
		Revoked:   0,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	_, err = m.Insert(context.Background(), token2)
	if err == nil {
		t.Error("重复 JTI 应返回错误")
	}
}
