package model

import (
	"context"
	"database/sql"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// MockUsersModel is a mock implementation of UsersModel for testing
type MockUsersModel struct {
	MockCount             func(ctx context.Context) (int64, error)
	MockInsert            func(ctx context.Context, data *Users) (sql.Result, error)
	MockFindOne           func(ctx context.Context, id int64) (*Users, error)
	MockFindOneByUsername func(ctx context.Context, username string) (*Users, error)
	MockUpdate            func(ctx context.Context, data *Users) error
	MockDelete            func(ctx context.Context, id int64) error
	MockWithSession       func(session sqlx.Session) UsersModel
}

func (m *MockUsersModel) Count(ctx context.Context) (int64, error) {
	if m.MockCount != nil {
		return m.MockCount(ctx)
	}
	return 0, nil
}

func (m *MockUsersModel) Insert(ctx context.Context, data *Users) (sql.Result, error) {
	if m.MockInsert != nil {
		return m.MockInsert(ctx, data)
	}
	return nil, nil
}

func (m *MockUsersModel) FindOne(ctx context.Context, id int64) (*Users, error) {
	if m.MockFindOne != nil {
		return m.MockFindOne(ctx, id)
	}
	return nil, ErrNotFound
}

func (m *MockUsersModel) FindOneByUsername(ctx context.Context, username string) (*Users, error) {
	if m.MockFindOneByUsername != nil {
		return m.MockFindOneByUsername(ctx, username)
	}
	return nil, ErrNotFound
}

func (m *MockUsersModel) Update(ctx context.Context, data *Users) error {
	if m.MockUpdate != nil {
		return m.MockUpdate(ctx, data)
	}
	return nil
}

func (m *MockUsersModel) Delete(ctx context.Context, id int64) error {
	if m.MockDelete != nil {
		return m.MockDelete(ctx, id)
	}
	return nil
}

func (m *MockUsersModel) withSession(session sqlx.Session) UsersModel {
	if m.MockWithSession != nil {
		return m.MockWithSession(session)
	}
	return m
}

// MockResult is a mock implementation of sql.Result for testing
type MockResult struct {
	LastInsertIdVal int64
	LastInsertIdErr error
	RowsAffectedVal int64
	RowsAffectedErr error
}

func (m *MockResult) LastInsertId() (int64, error) {
	return m.LastInsertIdVal, m.LastInsertIdErr
}

func (m *MockResult) RowsAffected() (int64, error) {
	return m.RowsAffectedVal, m.RowsAffectedErr
}
