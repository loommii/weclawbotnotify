package config

import (
	"time"

	"github.com/zeromicro/go-zero/rest"
)

// AuthConfig JWT 双令牌认证配置
type AuthConfig struct {
	PrivateKeyPath string        // RSA 私钥文件路径
	PublicKeyPath  string        // RSA 公钥文件路径
	AccessExpire   time.Duration // Access Token 过期时间（支持 15m、24h 等格式）
	RefreshExpire  time.Duration // Refresh Token 过期时间（支持 7d、168h 等格式）
}

type Config struct {
	rest.RestConf
	DataSource string     // 数据库连接字符串
	Auth       AuthConfig // JWT 认证配置
}
