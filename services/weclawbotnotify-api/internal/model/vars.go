package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var ErrNotFound = sqlx.ErrNotFound

// Application 状态常量（Go 避免 0 值原则）
const (
	AppStatusActive    = 1  // 正常：可查询、可推送
	AppStatusDisabled  = 2  // 禁用：临时停用，Token 失效，可恢复
	AppStatusDeleted   = -1 // 已删除：软删除，Token 失效，不可恢复
)
