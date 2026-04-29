// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package application

import (
	"context"

	"weclawbotnotify/services/weclawbotnotify-api/internal/svc"
	"weclawbotnotify/services/weclawbotnotify-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateApplicationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateApplicationLogic {
	return &CreateApplicationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateApplicationLogic) CreateApplication(req *types.CreateApplicationReq) (resp *types.CreateApplicationResp, err error) {
	// todo: add your logic here and delete this line

	return
}
