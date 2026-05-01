package config

import (
	"time"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type AuthConfig struct {
	PrivateKeyPath string
	PublicKeyPath  string
	AccessExpire   time.Duration
	RefreshExpire  time.Duration
}

type Config struct {
	rest.RestConf
	DataSource string
	Auth       AuthConfig
	IlinkRpc   zrpc.RpcClientConf
}
