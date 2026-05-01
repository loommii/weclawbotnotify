package logic

import (
	"context"

	"weclawbotnotify/services/ilink/internal/svc"
	"weclawbotnotify/services/ilink/pb/weclawbotnotify/services/ilink/ilink"

	"github.com/zeromicro/go-zero/core/logx"
)

type PollQRCodeStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPollQRCodeStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PollQRCodeStatusLogic {
	return &PollQRCodeStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PollQRCodeStatusLogic) PollQRCodeStatus(in *ilink.PollQRCodeStatusReq) (*ilink.PollQRCodeStatusResp, error) {
	l.Infof("[PollQRCodeStatus] qrcode=%s", in.Qrcode)

	resp, err := l.svcCtx.UnauthenticatedCli.PollQRCodeStatus(l.ctx, in.Qrcode)
	if err != nil {
		l.Errorf("[PollQRCodeStatus] failed: %v", err)
		return nil, err
	}

	return &ilink.PollQRCodeStatusResp{
		Status:      resp.Status,
		BotToken:    resp.BotToken,
		IlinkBotId:  resp.ILinkBotID,
		BaseUrl:     resp.BaseURL,
		IlinkUserId: resp.ILinkUserID,
	}, nil
}
