// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package client

import (
	"context"

	"weclawbotnotify/services/weclawbotnotify-api/internal/svc"
	"weclawbotnotify/services/weclawbotnotify-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListClientsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListClientsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListClientsLogic {
	return &ListClientsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListClientsLogic) ListClients(req *types.ListClientsReq) (resp *types.ListClientsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
