package model

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func newTestDB(t *testing.T) (sqlx.SqlConn, func()) {
	t.Helper()
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")
	conn := sqlx.NewSqlConn("sqlite", dbPath)
	_, err := conn.ExecCtx(context.Background(), `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username VARCHAR(50) NOT NULL UNIQUE,
			password_hash CHAR(60) NOT NULL,
			created_at INT NOT NULL DEFAULT (strftime('%s', 'now')),
			updated_at INT NOT NULL DEFAULT (strftime('%s', 'now'))
		)
	`)
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}
	return conn, func() { os.Remove(dbPath) }
}

func TestCustomUsersModel_Count(t *testing.T) {
	conn, _ := newTestDB(t)
	m := NewUsersModel(conn)

	count, err := m.Count(context.Background())
	if err != nil {
		t.Fatalf("Count failed: %v", err)
	}
	if count != 0 {
		t.Errorf("expected count 0, got %d", count)
	}

	_, err = m.Insert(context.Background(), &Users{
		Username:     "testuser1",
		PasswordHash: "hash123",
	})
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	count, err = m.Count(context.Background())
	if err != nil {
		t.Fatalf("Count failed: %v", err)
	}
	if count != 1 {
		t.Errorf("expected count 1, got %d", count)
	}
}

func TestCustomUsersModel_InsertAndFindOne(t *testing.T) {
	conn, _ := newTestDB(t)
	m := NewUsersModel(conn)

	result, err := m.Insert(context.Background(), &Users{
		Username:     "testuser2",
		PasswordHash: "hash123",
	})
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("LastInsertId failed: %v", err)
	}
	if id != 1 {
		t.Errorf("expected id 1, got %d", id)
	}

	user, err := m.FindOneByUsername(context.Background(), "testuser2")
	if err != nil {
		t.Fatalf("FindOneByUsername failed: %v", err)
	}
	if user.Username != "testuser2" {
		t.Errorf("expected username testuser2, got %s", user.Username)
	}
	if user.PasswordHash != "hash123" {
		t.Errorf("expected password_hash hash123, got %s", user.PasswordHash)
	}

	user, err = m.FindOne(context.Background(), id)
	if err != nil {
		t.Fatalf("FindOne failed: %v", err)
	}
	if user.Id != id {
		t.Errorf("expected id %d, got %d", id, user.Id)
	}

	_, err = m.FindOneByUsername(context.Background(), "nonexistent")
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestCustomUsersModel_Update(t *testing.T) {
	conn, _ := newTestDB(t)
	m := NewUsersModel(conn)

	result, err := m.Insert(context.Background(), &Users{
		Username:     "testuser3",
		PasswordHash: "oldhash",
	})
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	id, _ := result.LastInsertId()

	err = m.Update(context.Background(), &Users{
		Id:           id,
		Username:     "newuser3",
		PasswordHash: "newhash",
	})
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	user, err := m.FindOne(context.Background(), id)
	if err != nil {
		t.Fatalf("FindOne failed: %v", err)
	}
	if user.Username != "newuser3" {
		t.Errorf("expected username newuser3, got %s", user.Username)
	}
	if user.PasswordHash != "newhash" {
		t.Errorf("expected password_hash newhash, got %s", user.PasswordHash)
	}
}

func TestCustomUsersModel_Delete(t *testing.T) {
	conn, _ := newTestDB(t)
	m := NewUsersModel(conn)

	result, err := m.Insert(context.Background(), &Users{
		Username:     "testuser4",
		PasswordHash: "hash123",
	})
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	id, _ := result.LastInsertId()

	err = m.Delete(context.Background(), id)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = m.FindOne(context.Background(), id)
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestCustomUsersModel_WithSession(t *testing.T) {
	conn, _ := newTestDB(t)
	m := NewUsersModel(conn)
	mWithSession := m.withSession(nil)
	if mWithSession == nil {
		t.Error("withSession returned nil")
	}
}

func TestMockUsersModel(t *testing.T) {
	mock := &MockUsersModel{
		MockCount: func(ctx context.Context) (int64, error) {
			return 42, nil
		},
		MockFindOneByUsername: func(ctx context.Context, username string) (*Users, error) {
			return &Users{Id: 1, Username: username}, nil
		},
	}

	count, err := mock.Count(context.Background())
	if err != nil {
		t.Fatalf("Mock Count failed: %v", err)
	}
	if count != 42 {
		t.Errorf("expected count 42, got %d", count)
	}

	user, err := mock.FindOneByUsername(context.Background(), "test")
	if err != nil {
		t.Fatalf("Mock FindOneByUsername failed: %v", err)
	}
	if user.Username != "test" {
		t.Errorf("expected username test, got %s", user.Username)
	}
}

func TestMockResult(t *testing.T) {
	mock := &MockResult{
		LastInsertIdVal: 123,
		RowsAffectedVal: 1,
	}

	id, err := mock.LastInsertId()
	if err != nil {
		t.Fatalf("Mock LastInsertId failed: %v", err)
	}
	if id != 123 {
		t.Errorf("expected id 123, got %d", id)
	}

	rows, err := mock.RowsAffected()
	if err != nil {
		t.Fatalf("Mock RowsAffected failed: %v", err)
	}
	if rows != 1 {
		t.Errorf("expected rows 1, got %d", rows)
	}
}
