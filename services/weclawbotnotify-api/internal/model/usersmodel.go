package model

import (
	"context"
	"database/sql"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UsersModel = (*customUsersModel)(nil)

type (
	// UsersModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUsersModel.
	UsersModel interface {
		usersModel
		withSession(session sqlx.Session) UsersModel
		Count(ctx context.Context) (int64, error)
	}

	customUsersModel struct {
		*defaultUsersModel
	}
)

// NewUsersModel returns a model for the database table.
func NewUsersModel(conn sqlx.SqlConn) UsersModel {
	return &customUsersModel{
		defaultUsersModel: newUsersModel(conn),
	}
}

func (m *customUsersModel) withSession(session sqlx.Session) UsersModel {
	return NewUsersModel(sqlx.NewSqlConnFromSession(session))
}

// Count 返回用户总数
func (m *customUsersModel) Count(ctx context.Context) (int64, error) {
	var count int64
	query := "SELECT COUNT(*) FROM " + m.table
	err := m.conn.QueryRowCtx(ctx, &count, query)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	return count, nil
}
