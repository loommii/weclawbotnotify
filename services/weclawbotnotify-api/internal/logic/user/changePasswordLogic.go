package user

import (
	"context"
	"fmt"

	"weclawbotnotify/pkg/jwtx"
	pkgmw "weclawbotnotify/pkg/middleware"
	"weclawbotnotify/pkg/xerr"
	"weclawbotnotify/services/weclawbotnotify-api/internal/svc"
	"weclawbotnotify/services/weclawbotnotify-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type ChangePasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePasswordLogic {
	return &ChangePasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ChangePassword 修改密码：验证旧密码 → 更新密码 → 撤销所有 Refresh Token
func (l *ChangePasswordLogic) ChangePassword(req *types.ChangePasswordReq) (resp *types.ChangePasswordResp, err error) {
	claims, ok := l.ctx.Value(pkgmw.ClaimsKey).(*jwtx.JWTClaims)
	if !ok || claims == nil {
		return nil, xerr.JwtError
	}

	userId, err := parseUserId(claims.UID)
	if err != nil {
		return nil, xerr.JwtError
	}

	user, err := l.svcCtx.UsersModel.FindOne(l.ctx, userId)
	if err != nil {
		l.Errorf("查询用户失败: %v", err)
		return nil, xerr.JwtError
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)); err != nil {
		l.Infof("旧密码错误: userId=%d", userId)
		return nil, xerr.LoginPasswordWrong
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		l.Errorf("密码哈希失败: %v", err)
		return nil, xerr.RegisterHashFailed
	}

	user.PasswordHash = string(hashedPassword)
	if err := l.svcCtx.UsersModel.Update(l.ctx, user); err != nil {
		l.Errorf("更新密码失败: %v", err)
		return nil, xerr.RegisterInsertFailed
	}

	// 撤销该用户所有 Refresh Token，强制其他设备重新登录
	if err := l.svcCtx.RefreshTokensModel.RevokeByUserId(l.ctx, userId); err != nil {
		l.Errorf("撤销 Refresh Token 失败: %v", err)
	}

	l.Infof("密码修改成功: userId=%d", userId)

	return &types.ChangePasswordResp{}, nil
}

func parseUserId(uid string) (int64, error) {
	var id int64
	_, err := fmt.Sscanf(uid, "%d", &id)
	return id, err
}
