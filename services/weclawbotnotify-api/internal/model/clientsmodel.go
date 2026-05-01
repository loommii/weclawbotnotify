package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ClientsModel = (*customClientsModel)(nil)

type (
	ClientsModel interface {
		clientsModel
		withSession(session sqlx.Session) ClientsModel
		FindByUserIdWithPagination(ctx context.Context, userId int64, offset, limit int) ([]*Clients, error)
		CountByUserId(ctx context.Context, userId int64) (int64, error)
	}

	customClientsModel struct {
		*defaultClientsModel
	}
)

func NewClientsModel(conn sqlx.SqlConn) ClientsModel {
	return &customClientsModel{
		defaultClientsModel: newClientsModel(conn),
	}
}

func (m *customClientsModel) withSession(session sqlx.Session) ClientsModel {
	return NewClientsModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customClientsModel) FindByUserIdWithPagination(ctx context.Context, userId int64, offset, limit int) ([]*Clients, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ? order by `created_at` desc limit ? offset ?", clientsRows, m.table)
	var resp []*Clients
	err := m.conn.QueryRowsCtx(ctx, &resp, query, userId, limit, offset)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *customClientsModel) CountByUserId(ctx context.Context, userId int64) (int64, error) {
	query := fmt.Sprintf("select count(*) from %s where `user_id` = ?", m.table)
	var resp int64
	err := m.conn.QueryRowCtx(ctx, &resp, query, userId)
	if err != nil {
		return 0, err
	}
	return resp, nil
}
