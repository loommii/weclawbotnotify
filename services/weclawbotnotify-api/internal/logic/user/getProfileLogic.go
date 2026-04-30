package user

import (
	"context"
	"fmt"
	"strconv"

	"weclawbotnotify/pkg/jwtx"
	"weclawbotnotify/pkg/middleware"
	"weclawbotnotify/pkg/xerr"
	"weclawbotnotify/services/weclawbotnotify-api/internal/model"
	"weclawbotnotify/services/weclawbotnotify-api/internal/svc"
	"weclawbotnotify/services/weclawbotnotify-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProfileLogic {
	return &GetProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProfileLogic) GetProfile() (resp *types.ProfileResp, err error) {
	claims, ok := l.ctx.Value(middleware.ClaimsKey).(*jwtx.JWTClaims)
	if !ok || claims == nil {
		l.Errorf("获取 JWT claims 失败")
		return nil, xerr.JwtError
	}

	l.Infof("获取用户信息请求: userId=%s", claims.UID)

	userId, err := strconv.ParseInt(claims.UID, 10, 64)
	if err != nil {
		l.Errorf("解析 userId 失败: %v", err)
		return nil, xerr.JwtError
	}

	user, err := l.svcCtx.UsersModel.FindOne(l.ctx, userId)
	if err != nil {
		if err == model.ErrNotFound {
			l.Errorf("用户不存在: userId=%d", userId)
			return nil, xerr.JwtError
		}
		l.Errorf("查询用户失败: userId=%d, err=%v", userId, err)
		return nil, xerr.RegisterQueryFailed
	}

	l.Infof("获取用户信息成功: userId=%d, username=%s", user.Id, user.Username)

	return &types.ProfileResp{
		User: types.UserInfo{
			Id:        user.Id,
			Username:  user.Username,
			CreatedAt: fmt.Sprintf("%d", user.CreatedAt),
		},
	}, nil
}
