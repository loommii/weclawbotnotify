package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ApplicationsModel = (*customApplicationsModel)(nil)

type (
	ApplicationsModel interface {
		applicationsModel
		withSession(session sqlx.Session) ApplicationsModel
		FindByUserId(ctx context.Context, userId int64) ([]*Applications, error)
		FindByUserIdWithPagination(ctx context.Context, userId int64, offset, limit int) ([]*Applications, error)
		CountByUserId(ctx context.Context, userId int64) (int64, error)
	}

	customApplicationsModel struct {
		*defaultApplicationsModel
	}
)

func NewApplicationsModel(conn sqlx.SqlConn) ApplicationsModel {
	return &customApplicationsModel{
		defaultApplicationsModel: newApplicationsModel(conn),
	}
}

func (m *customApplicationsModel) withSession(session sqlx.Session) ApplicationsModel {
	return NewApplicationsModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customApplicationsModel) FindByUserId(ctx context.Context, userId int64) ([]*Applications, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ? order by `created_at` desc", applicationsRows, m.table)
	var resp []*Applications
	err := m.conn.QueryRowsCtx(ctx, &resp, query, userId)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *customApplicationsModel) FindByUserIdWithPagination(ctx context.Context, userId int64, offset, limit int) ([]*Applications, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ? order by `created_at` desc limit ? offset ?", applicationsRows, m.table)
	var resp []*Applications
	err := m.conn.QueryRowsCtx(ctx, &resp, query, userId, limit, offset)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *customApplicationsModel) CountByUserId(ctx context.Context, userId int64) (int64, error) {
	query := fmt.Sprintf("select count(*) from %s where `user_id` = ?", m.table)
	var resp int64
	err := m.conn.QueryRowCtx(ctx, &resp, query, userId)
	if err != nil {
		return 0, err
	}
	return resp, nil
}
