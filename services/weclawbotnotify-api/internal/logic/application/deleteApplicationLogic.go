// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package application

import (
	"context"

	"weclawbotnotify/services/weclawbotnotify-api/internal/svc"
	"weclawbotnotify/services/weclawbotnotify-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteApplicationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteApplicationLogic {
	return &DeleteApplicationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteApplicationLogic) DeleteApplication(req *types.DeleteApplicationReq) (resp *types.DeleteApplicationResp, err error) {
	// todo: add your logic here and delete this line

	return
}
