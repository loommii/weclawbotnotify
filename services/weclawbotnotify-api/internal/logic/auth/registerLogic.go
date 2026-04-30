package auth

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"weclawbotnotify/pkg/jwtx"
	"weclawbotnotify/pkg/xerr"
	"weclawbotnotify/services/weclawbotnotify-api/internal/model"
	"weclawbotnotify/services/weclawbotnotify-api/internal/svc"
	"weclawbotnotify/services/weclawbotnotify-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Register 用户注册：参数校验 → 单用户限制 → 用户名查重 → 创建用户 → 签发双令牌
func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	l.Infof("用户注册请求: username=%s", req.Username)

	if req.Username == "" || req.Password == "" {
		l.Errorf("参数校验失败: 用户名或密码为空")
		return nil, xerr.RegisterParamEmpty
	}

	count, err := l.svcCtx.UsersModel.Count(l.ctx)
	if err != nil {
		l.Errorf("查询用户数量失败: %v", err)
		return nil, xerr.RegisterQueryFailed
	}
	if count > 0 {
		l.Infof("拒绝注册: 系统已存在用户，V1版本仅支持单用户")
		return nil, xerr.RegisterClosed
	}

	_, err = l.svcCtx.UsersModel.FindOneByUsername(l.ctx, req.Username)
	if err == nil {
		l.Infof("用户名已存在: %s", req.Username)
		return nil, xerr.RegisterUsernameExist
	}
	if !isNotFound(err) {
		l.Errorf("查询用户失败: %v", err)
		return nil, xerr.RegisterQueryFailed
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		l.Errorf("密码哈希失败: %v", err)
		return nil, xerr.RegisterHashFailed
	}

	result, err := l.svcCtx.UsersModel.Insert(l.ctx, &model.Users{
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
	})
	if err != nil {
		l.Errorf("创建用户失败: %v", err)
		return nil, xerr.RegisterInsertFailed
	}

	userId, err := result.LastInsertId()
	if err != nil {
		l.Errorf("获取用户ID失败: %v", err)
		return nil, xerr.RegisterGetIdFailed
	}

	// 签发双令牌
	accessToken, refreshToken, err := l.generateTokenPair(userId)
	if err != nil {
		l.Errorf("生成令牌失败: %v", err)
		return nil, xerr.RegisterTokenFailed
	}

	l.Infof("用户注册成功: userId=%d, username=%s", userId, req.Username)

	return &types.RegisterResp{
		Token:        accessToken,
		RefreshToken: refreshToken,
		User: types.UserInfo{
			Id:        userId,
			Username:  req.Username,
			CreatedAt: fmt.Sprintf("%d", time.Now().Unix()),
		},
	}, nil
}

// generateTokenPair 签发双令牌并持久化 Refresh Token
func (l *RegisterLogic) generateTokenPair(userId int64) (string, string, error) {
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

	// 从 refreshToken 中提取 JTI 和过期时间
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

func isNotFound(err error) bool {
	return err == model.ErrNotFound
}
