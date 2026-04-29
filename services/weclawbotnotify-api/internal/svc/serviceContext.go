// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"github.com/zeromicro/go-zero/rest"
	"weclawbotnotify/services/weclawbotnotify-api/internal/config"
	"weclawbotnotify/services/weclawbotnotify-api/internal/middleware"
)

type ServiceContext struct {
	Config          config.Config
	ClientAuth      rest.Middleware
	ApplicationAuth rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:          c,
		ClientAuth:      middleware.NewClientAuthMiddleware().Handle,
		ApplicationAuth: middleware.NewApplicationAuthMiddleware().Handle,
	}
}
