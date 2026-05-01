// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package message

import (
	"context"
	"database/sql"
	"time"

	"weclawbotnotify/pkg/xerr"
	"weclawbotnotify/services/weclawbotnotify-api/internal/model"
	"weclawbotnotify/services/weclawbotnotify-api/internal/svc"
	"weclawbotnotify/services/weclawbotnotify-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateMessageLogic {
	return &CreateMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateMessageLogic) CreateMessage(req *types.CreateMessageReq) (resp *types.CreateMessageResp, err error) {
	if req.AppToken == "" {
		return nil, xerr.ApplicationTokenInvalid
	}

	app, err := l.svcCtx.ApplicationsModel.ValidateAppToken(l.ctx, req.AppToken)
	if err != nil {
		l.Errorf("validate app token failed: %v", err)
		return nil, xerr.ApplicationTokenInvalid
	}

	l.Infof("推送消息: appId=%d, appName=%s, title=%s", app.Id, app.Name, req.Title)

	now := time.Now().Unix()
	app.LastUsedAt = sql.NullInt64{Int64: now, Valid: true}
	if err := l.svcCtx.ApplicationsModel.Update(l.ctx, app); err != nil {
		l.Errorf("update application last_used_at failed: %v", err)
	}

	clients, err := l.svcCtx.ClientsModel.FindByUserIdWithPagination(l.ctx, app.UserId, 0, 100)
	if err != nil {
		l.Errorf("query clients failed: %v", err)
		return nil, xerr.MessageQueryFailed
	}

	var boundClients []*model.Clients
	for _, c := range clients {
		if c.Status == "bound" {
			boundClients = append(boundClients, c)
		}
	}

	if len(boundClients) == 0 {
		return nil, xerr.MessageNoClientAvailable
	}

	l.Infof("找到 %d 个可用客户端，开始推送", len(boundClients))

	return nil, nil
}
