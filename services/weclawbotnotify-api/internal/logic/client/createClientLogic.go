// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package client

import (
	"context"

	"weclawbotnotify/services/weclawbotnotify-api/internal/svc"
	"weclawbotnotify/services/weclawbotnotify-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateClientLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateClientLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateClientLogic {
	return &CreateClientLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateClientLogic) CreateClient(req *types.CreateClientReq) (resp *types.CreateClientResp, err error) {
	// todo: add your logic here and delete this line

	return
}
