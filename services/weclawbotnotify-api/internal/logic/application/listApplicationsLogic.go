// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package application

import (
	"context"

	"weclawbotnotify/services/weclawbotnotify-api/internal/svc"
	"weclawbotnotify/services/weclawbotnotify-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListApplicationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListApplicationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListApplicationsLogic {
	return &ListApplicationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListApplicationsLogic) ListApplications() (resp *types.ListApplicationsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
