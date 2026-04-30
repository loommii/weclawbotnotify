package user

import (
	"context"
	"strconv"

	"weclawbotnotify/pkg/jwtx"
	pkgmw "weclawbotnotify/pkg/middleware"
	"weclawbotnotify/pkg/password"
	"weclawbotnotify/pkg/xerr"
	"weclawbotnotify/services/weclawbotnotify-api/internal/model"
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

// ChangePassword 修改密码：验证新密码强度 → 验证旧密码 → 更新密码 → 撤销所有 Refresh Token
func (l *ChangePasswordLogic) ChangePassword(req *types.ChangePasswordReq) (resp *types.ChangePasswordResp, err error) {
	claims, ok := l.ctx.Value(pkgmw.ClaimsKey).(*jwtx.JWTClaims)
	if !ok || claims == nil {
		l.Errorf("获取 JWT claims 失败")
		return nil, xerr.JwtError
	}

	userId, err := strconv.ParseInt(claims.UID, 10, 64)
	if err != nil {
		l.Errorf("解析 userId 失败: %v", err)
		return nil, xerr.JwtError
	}

	l.Infof("修改密码请求: userId=%d", userId)

	// 验证新密码强度（长度>=8，包含大小写字母和数字）
	if err := password.ValidateStrength(req.NewPassword); err != nil {
		l.Infof("新密码强度不足: userId=%d", userId)
		return nil, xerr.PasswordTooWeak
	}

	// 查询用户
	user, err := l.svcCtx.UsersModel.FindOne(l.ctx, userId)
	if err != nil {
		if err == model.ErrNotFound {
			l.Errorf("用户不存在: userId=%d", userId)
			return nil, xerr.JwtError
		}
		l.Errorf("查询用户失败: userId=%d, err=%v", userId, err)
		return nil, xerr.JwtError
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)); err != nil {
		l.Infof("旧密码错误: userId=%d", userId)
		return nil, xerr.LoginPasswordWrong
	}

	// 生成新密码哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		l.Errorf("密码哈希失败: userId=%d, err=%v", userId, err)
		return nil, xerr.RegisterHashFailed
	}

	// 更新密码
	user.PasswordHash = string(hashedPassword)
	if err := l.svcCtx.UsersModel.Update(l.ctx, user); err != nil {
		l.Errorf("更新密码失败: userId=%d, err=%v", userId, err)
		return nil, xerr.RegisterInsertFailed
	}

	// 撤销该用户所有 Refresh Token，强制所有设备重新登录
	if err := l.svcCtx.RefreshTokensModel.RevokeByUserId(l.ctx, userId); err != nil {
		l.Errorf("撤销 Refresh Token 失败: userId=%d, err=%v", userId, err)
	}

	l.Infof("密码修改成功: userId=%d", userId)

	return &types.ChangePasswordResp{}, nil
}
