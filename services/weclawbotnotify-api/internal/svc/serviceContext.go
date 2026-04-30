package svc

import (
	"fmt"
	"os"
	"path/filepath"

	"weclawbotnotify/pkg/jwtx"
	pkgmw "weclawbotnotify/pkg/middleware"
	"weclawbotnotify/services/weclawbotnotify-api/internal/config"
	localmw "weclawbotnotify/services/weclawbotnotify-api/internal/middleware"
	"weclawbotnotify/services/weclawbotnotify-api/internal/model"

	_ "modernc.org/sqlite"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config             config.Config
	UsersModel         model.UsersModel
	RefreshTokensModel model.RefreshTokensModel
	AccessJWTHelper    *jwtx.JWTHelper
	RefreshJWTHelper   *jwtx.JWTHelper
	ClientAuth         rest.Middleware
	ApplicationAuth    rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := initDB(c.DataSource)
	usersModel := model.NewUsersModel(conn)
	refreshTokensModel := model.NewRefreshTokensModel(conn)
	accessJWTHelper, refreshJWTHelper := initJWT(c.Auth)

	publicKeyPEM, err := os.ReadFile(c.Auth.PublicKeyPath)
	if err != nil {
		logx.Must(err)
	}

	return &ServiceContext{
		Config:             c,
		UsersModel:         usersModel,
		RefreshTokensModel: refreshTokensModel,
		AccessJWTHelper:    accessJWTHelper,
		RefreshJWTHelper:   refreshJWTHelper,
		ClientAuth:         pkgmw.NewJWTMiddleware(publicKeyPEM).Handle,
		ApplicationAuth:    localmw.NewApplicationAuthMiddleware().Handle,
	}
}

// initDB 初始化数据库连接并自动建表
func initDB(dataSource string) sqlx.SqlConn {
	dir := filepath.Dir(dataSource)
	if dir != "." {
		os.MkdirAll(dir, 0755)
	}

	conn := sqlx.NewSqlConn("sqlite", dataSource+"?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)&_pragma=busy_timeout(5000)")

	// 依次执行建表 SQL
	tables := []string{
		"sql/table_users.sqlite.sql",
		"sql/table_refresh_tokens.sqlite.sql",
	}
	for _, table := range tables {
		sqlBytes, err := os.ReadFile(table)
		if err != nil {
			logx.Must(err)
		}
		if _, err = conn.ExecCtx(nil, string(sqlBytes)); err != nil {
			logx.Must(err)
		}
	}

	return conn
}

// initJWT 初始化双令牌 JWTHelper（RSA 签名）
func initJWT(c config.AuthConfig) (*jwtx.JWTHelper, *jwtx.JWTHelper) {
	if c.RefreshExpire <= c.AccessExpire {
		logx.Must(fmt.Errorf("RefreshExpire (%v) 必须大于 AccessExpire (%v)", c.RefreshExpire, c.AccessExpire))
	}

	privateKey, _, err := jwtx.ParseRSAPrivateKeyFromPath(c.PrivateKeyPath)
	if err != nil {
		logx.Must(err)
	}
	publicKey, _, err := jwtx.ParseRSAPublicKeyFromPath(c.PublicKeyPath)
	if err != nil {
		logx.Must(err)
	}

	accessJWTHelper := jwtx.NewJWTHelper(
		jwtx.WithPrivateKey(privateKey),
		jwtx.WithPublicKey(publicKey),
		jwtx.WithExpiredTime(c.AccessExpire),
	)

	refreshJWTHelper := jwtx.NewJWTHelper(
		jwtx.WithPrivateKey(privateKey),
		jwtx.WithPublicKey(publicKey),
		jwtx.WithExpiredTime(c.RefreshExpire),
	)

	return accessJWTHelper, refreshJWTHelper
}
