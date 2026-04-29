// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package health

import (
	"context"
	"time"

	"weclawbotnotify/services/weclawbotnotify-api/internal/svc"
	"weclawbotnotify/services/weclawbotnotify-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// HealthLogic 健康检查业务逻辑
type HealthLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHealthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HealthLogic {
	return &HealthLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Health 健康检查接口，返回服务运行状态
func (l *HealthLogic) Health() (resp *types.HealthResp, err error) {
	return &types.HealthResp{
		ServiceName: "weclawbotnotify-api",
		Timestamp:   time.Now().Unix(),
		Message:     "ok",
	}, nil
}
