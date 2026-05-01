package logic

import (
	"context"

	"weclawbotnotify/services/ilink/internal/ilinkhttp"
	"weclawbotnotify/services/ilink/internal/svc"
	"weclawbotnotify/services/ilink/pb/weclawbotnotify/services/ilink/ilink"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUpdatesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUpdatesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUpdatesLogic {
	return &GetUpdatesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUpdatesLogic) GetUpdates(in *ilink.GetUpdatesReq) (*ilink.GetUpdatesResp, error) {
	l.Infof("[GetUpdates] buf_len=%d", len(in.GetUpdatesBuf))

	cli := ilinkhttp.NewClient(in.BotToken, in.BaseUrl)
	resp, err := cli.GetUpdates(l.ctx, in.GetUpdatesBuf)
	if err != nil {
		l.Errorf("[GetUpdates] failed: %v", err)
		return nil, err
	}

	pbMsgs := make([]*ilink.WeixinMessage, 0, len(resp.Msgs))
	for _, msg := range resp.Msgs {
		pbMsg := &ilink.WeixinMessage{
			Seq:          int32(msg.Seq),
			MessageId:    msg.MessageID,
			FromUserId:   msg.FromUserID,
			ToUserId:     msg.ToUserID,
			MessageType:  int32(msg.MessageType),
			MessageState: int32(msg.MessageState),
			ContextToken: msg.ContextToken,
		}
		pbItems := make([]*ilink.MessageItem, 0, len(msg.ItemList))
		for _, item := range msg.ItemList {
			pbItem := &ilink.MessageItem{
				Type: int32(item.Type),
			}
			if item.TextItem != nil {
				pbItem.TextItem = &ilink.TextItem{
					Text: item.TextItem.Text,
				}
			}
			pbItems = append(pbItems, pbItem)
		}
		pbMsg.ItemList = pbItems
		pbMsgs = append(pbMsgs, pbMsg)
	}

	return &ilink.GetUpdatesResp{
		Ret:                int32(resp.Ret),
		ErrCode:            int32(resp.ErrCode),
		ErrMsg:             resp.ErrMsg,
		GetUpdatesBuf:      resp.GetUpdatesBuf,
		LongPollingTimeoutMs: int32(resp.LongPollingTimeoutMs),
		Msgs:               pbMsgs,
	}, nil
}
