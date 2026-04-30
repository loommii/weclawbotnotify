package svc

import (
	"os"
	"path/filepath"
	"time"

	"weclawbotnotify/pkg/jwtx"
	"weclawbotnotify/services/weclawbotnotify-api/internal/config"
	"weclawbotnotify/services/weclawbotnotify-api/internal/middleware"
	"weclawbotnotify/services/weclawbotnotify-api/internal/model"

	_ "modernc.org/sqlite"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config          config.Config
	UsersModel      model.UsersModel
	JWTHelper       *jwtx.JWTHelper
	ClientAuth      rest.Middleware
	ApplicationAuth rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := initDB(c.DataSource)
	usersModel := model.NewUsersModel(conn)
	jwtHelper := initJWT(c.Auth)

	return &ServiceContext{
		Config:          c,
		UsersModel:      usersModel,
		JWTHelper:       jwtHelper,
		ClientAuth:      middleware.NewClientAuthMiddleware().Handle,
		ApplicationAuth: middleware.NewApplicationAuthMiddleware().Handle,
	}
}

// initDB 初始化数据库连接并自动建表
// 注意：SQLite 是文件型数据库，无需账号密码，安全性由文件系统权限控制
func initDB(dataSource string) sqlx.SqlConn {
	dir := filepath.Dir(dataSource)
	if dir != "." {
		os.MkdirAll(dir, 0755)
	}

	conn := sqlx.NewSqlConn("sqlite", dataSource+"?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)&_pragma=busy_timeout(5000)")

	sqlBytes, err := os.ReadFile("sql/table_users.sqlite.sql")
	if err != nil {
		logx.Must(err)
	}
	_, err = conn.ExecCtx(nil, string(sqlBytes))
	if err != nil {
		logx.Must(err)
	}

	return conn
}

// initJWT 初始化 JWT 工具（RSA 签名）
func initJWT(c config.AuthConfig) *jwtx.JWTHelper {
	privateKey, _, err := jwtx.ParseRSAPrivateKeyFromPath(c.PrivateKeyPath)
	if err != nil {
		logx.Must(err)
	}
	publicKey, _, err := jwtx.ParseRSAPublicKeyFromPath(c.PublicKeyPath)
	if err != nil {
		logx.Must(err)
	}

	return jwtx.NewJWTHelper(
		jwtx.WithPrivateKey(privateKey),
		jwtx.WithPublicKey(publicKey),
		jwtx.WithExpiredTime(time.Duration(c.AccessExpire)*time.Second),
	)
}
