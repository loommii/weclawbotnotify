package application

import (
	"context"

	"weclawbotnotify/pkg/jwtx"
	pkgmw "weclawbotnotify/pkg/middleware"
	"weclawbotnotify/pkg/xerr"
	"weclawbotnotify/services/weclawbotnotify-api/internal/model"
	"weclawbotnotify/services/weclawbotnotify-api/internal/svc"
	"weclawbotnotify/services/weclawbotnotify-api/internal/types"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type DeleteApplicationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteApplicationLogic {
	return &DeleteApplicationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteApplicationLogic) DeleteApplication(req *types.DeleteApplicationReq) (resp *types.DeleteApplicationResp, err error) {
	userId, err := jwtx.GetUserIdFromContext(l.ctx, pkgmw.ClaimsKey)
	if err != nil {
		return nil, xerr.RequestParamError
	}

	app, err := l.svcCtx.ApplicationsModel.FindOne(l.ctx, req.Id)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, xerr.ApplicationNotFound
		}
		return nil, xerr.ApplicationQueryFailed
	}

	if app.UserId != userId {
		return nil, xerr.ApplicationNoPermission
	}

	if app.Status == model.AppStatusDeleted {
		return nil, xerr.ApplicationNotFound
	}

	err = l.svcCtx.ApplicationsModel.UpdateStatus(l.ctx, req.Id, model.AppStatusDeleted)
	if err != nil {
		return nil, xerr.ApplicationDeleteFailed
	}

	return &types.DeleteApplicationResp{}, nil
}
