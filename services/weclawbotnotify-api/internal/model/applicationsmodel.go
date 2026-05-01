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
		UpdateStatus(ctx context.Context, id int64, status int64) error
		ValidateAppToken(ctx context.Context, token string) (*Applications, error)
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
	query := fmt.Sprintf("select %s from %s where `user_id` = ? AND `status` != -1 order by `created_at` desc", applicationsRows, m.table)
	var resp []*Applications
	err := m.conn.QueryRowsCtx(ctx, &resp, query, userId)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *customApplicationsModel) FindByUserIdWithPagination(ctx context.Context, userId int64, offset, limit int) ([]*Applications, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ? AND `status` != -1 order by `created_at` desc limit ? offset ?", applicationsRows, m.table)
	var resp []*Applications
	err := m.conn.QueryRowsCtx(ctx, &resp, query, userId, limit, offset)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *customApplicationsModel) CountByUserId(ctx context.Context, userId int64) (int64, error) {
	query := fmt.Sprintf("select count(*) from %s where `user_id` = ? AND `status` != -1", m.table)
	var resp int64
	err := m.conn.QueryRowCtx(ctx, &resp, query, userId)
	if err != nil {
		return 0, err
	}
	return resp, nil
}

func (m *customApplicationsModel) UpdateStatus(ctx context.Context, id int64, status int64) error {
	query := fmt.Sprintf("update %s set `status` = ? where `id` = ? AND `status` != -1", m.table)
	_, err := m.conn.ExecCtx(ctx, query, status, id)
	return err
}

func (m *customApplicationsModel) ValidateAppToken(ctx context.Context, token string) (*Applications, error) {
	app, err := m.FindOneByToken(ctx, token)
	if err != nil {
		return nil, err
	}
	if app.Status != AppStatusActive {
		return nil, ErrNotFound
	}
	return app, nil
}
