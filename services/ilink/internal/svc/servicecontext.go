package svc

import (
	"weclawbotnotify/services/ilink/internal/config"
	"weclawbotnotify/services/ilink/internal/ilinkhttp"

	"github.com/zeromicro/go-zero/core/logx"
)

type ServiceContext struct {
	Config              config.Config
	UnauthenticatedCli  *ilinkhttp.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	unauthenticatedCli := ilinkhttp.NewUnauthenticatedClient(c.Ilink.BaseURL)
	logx.Infof("[ilink-svc] initialized with BaseURL: %s", c.Ilink.BaseURL)

	return &ServiceContext{
		Config:             c,
		UnauthenticatedCli: unauthenticatedCli,
	}
}
