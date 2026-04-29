// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package client

import (
	"context"

	"weclawbotnotify/services/weclawbotnotify-api/internal/svc"
	"weclawbotnotify/services/weclawbotnotify-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PollClientStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPollClientStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PollClientStatusLogic {
	return &PollClientStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PollClientStatusLogic) PollClientStatus(req *types.PollClientStatusReq) (resp *types.PollClientStatusResp, err error) {
	// todo: add your logic here and delete this line

	return
}
