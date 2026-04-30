package config

import "github.com/zeromicro/go-zero/rest"

type AuthConfig struct {
	PrivateKeyPath string
	PublicKeyPath  string
	AccessExpire   int64
}

type Config struct {
	rest.RestConf
	DataSource string
	Auth       AuthConfig
}
