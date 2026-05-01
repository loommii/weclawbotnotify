package logic

import (
	"context"

	"weclawbotnotify/services/ilink/internal/svc"
	"weclawbotnotify/services/ilink/pb/weclawbotnotify/services/ilink/ilink"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetQRCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetQRCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetQRCodeLogic {
	return &GetQRCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetQRCodeLogic) GetQRCode(in *ilink.GetQRCodeReq) (*ilink.GetQRCodeResp, error) {
	l.Infof("[GetQRCode] bot_type=%d", in.BotType)

	resp, err := l.svcCtx.UnauthenticatedCli.GetQRCode(l.ctx, in.BotType)
	if err != nil {
		l.Errorf("[GetQRCode] failed: %v", err)
		return nil, err
	}

	return &ilink.GetQRCodeResp{
		Qrcode:           resp.QRCode,
		QrcodeImgContent: resp.QRCodeImgContent,
	}, nil
}
