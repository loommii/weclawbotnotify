package auth

import (
	"context"
	"fmt"
	"strconv"

	"weclawbotnotify/pkg/jwtx"
	"weclawbotnotify/pkg/xerr"
	"weclawbotnotify/services/weclawbotnotify-api/internal/model"
	"weclawbotnotify/services/weclawbotnotify-api/internal/svc"
	"weclawbotnotify/services/weclawbotnotify-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Login 用户登录：参数校验 → 查找用户 → 验证密码 → 签发双令牌
func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	l.Infof("用户登录请求: username=%s", req.Username)

	if req.Username == "" || req.Password == "" {
		l.Errorf("参数校验失败: 用户名或密码为空")
		return nil, xerr.LoginParamEmpty
	}

	user, err := l.svcCtx.UsersModel.FindOneByUsername(l.ctx, req.Username)
	if err != nil {
		if isNotFound(err) {
			l.Infof("用户不存在: %s", req.Username)
			return nil, xerr.LoginUserNotFound
		}
		l.Errorf("查询用户失败: %v", err)
		return nil, xerr.LoginQueryFailed
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		l.Infof("密码错误: username=%s", req.Username)
		return nil, xerr.LoginPasswordWrong
	}

	// 签发双令牌
	accessToken, refreshToken, err := l.generateTokenPair(user.Id)
	if err != nil {
		l.Errorf("生成令牌失败: %v", err)
		return nil, xerr.LoginTokenFailed
	}

	l.Infof("用户登录成功: userId=%d, username=%s", user.Id, user.Username)

	return &types.LoginResp{
		Token:        accessToken,
		RefreshToken: refreshToken,
		User: types.UserInfo{
			Id:        user.Id,
			Username:  user.Username,
			CreatedAt: fmt.Sprintf("%d", user.CreatedAt),
		},
	}, nil
}

// generateTokenPair 签发双令牌并持久化 Refresh Token
func (l *LoginLogic) generateTokenPair(userId int64) (string, string, error) {
	uid := strconv.FormatInt(userId, 10)

	accessToken, err := l.svcCtx.AccessJWTHelper.GenerateToken(jwtx.JWTClaims{
		UID:       uid,
		TokenType: jwtx.Access,
	})
	if err != nil {
		return "", "", err
	}

	refreshToken, err := l.svcCtx.RefreshJWTHelper.GenerateToken(jwtx.JWTClaims{
		UID:       uid,
		TokenType: jwtx.Refresh,
	})
	if err != nil {
		return "", "", err
	}

	_, claims, err := l.svcCtx.RefreshJWTHelper.ValidateToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	_, err = l.svcCtx.RefreshTokensModel.Insert(l.ctx, &model.RefreshTokens{
		UserId:    userId,
		Jti:       claims.JTI,
		Revoked:   0,
		ExpiresAt: claims.ExpiresAt.Unix(),
	})
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
