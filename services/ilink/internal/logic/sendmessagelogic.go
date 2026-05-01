package logic

import (
	"context"

	"weclawbotnotify/services/ilink/internal/ilinkhttp"
	"weclawbotnotify/services/ilink/internal/svc"
	"weclawbotnotify/services/ilink/pb/weclawbotnotify/services/ilink/ilink"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMessageLogic {
	return &SendMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendMessageLogic) SendMessage(in *ilink.SendMessageReq) (*ilink.SendMessageResp, error) {
	l.Infof("[SendMessage] to_user_id=%s, text_len=%d, context_token=%t", in.ToUserId, len(in.Text), in.ContextToken != "")

	cli := ilinkhttp.NewClient(in.BotToken, in.BaseUrl)
	resp, err := cli.SendMessage(l.ctx, in.ToUserId, in.Text, in.ContextToken)
	if err != nil {
		l.Errorf("[SendMessage] failed: to_user_id=%s, err=%v", in.ToUserId, err)
		return nil, err
	}

	return &ilink.SendMessageResp{
		Ret:    int32(resp.Ret),
		ErrMsg: resp.ErrMsg,
	}, nil
}
