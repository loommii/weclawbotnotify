package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ApplicationsModel = (*customApplicationsModel)(nil)

type (
	// ApplicationsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customApplicationsModel.
	ApplicationsModel interface {
		applicationsModel
		withSession(session sqlx.Session) ApplicationsModel
	}

	customApplicationsModel struct {
		*defaultApplicationsModel
	}
)

// NewApplicationsModel returns a model for the database table.
func NewApplicationsModel(conn sqlx.SqlConn) ApplicationsModel {
	return &customApplicationsModel{
		defaultApplicationsModel: newApplicationsModel(conn),
	}
}

func (m *customApplicationsModel) withSession(session sqlx.Session) ApplicationsModel {
	return NewApplicationsModel(sqlx.NewSqlConnFromSession(session))
}
