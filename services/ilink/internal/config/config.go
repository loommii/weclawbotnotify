package config

import "github.com/zeromicro/go-zero/zrpc"

type IlinkConfig struct {
	BaseURL string
}

type Config struct {
	zrpc.RpcServerConf
	Ilink IlinkConfig
}
