package auth

import (
	"context"
	"strconv"
	"time"

	"weclawbotnotify/pkg/jwtx"
	"weclawbotnotify/pkg/xerr"
	"weclawbotnotify/services/weclawbotnotify-api/internal/model"
	"weclawbotnotify/services/weclawbotnotify-api/internal/svc"
	"weclawbotnotify/services/weclawbotnotify-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefreshLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshLogic {
	return &RefreshLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Refresh 刷新令牌：验证签名 → 校验TokenType → 查询JTI → 检查撤销状态 → 轮换签发新TokenPair
func (l *RefreshLogic) Refresh(req *types.RefreshReq) (resp *types.RefreshResp, err error) {
	l.Infof("刷新令牌请求")

	if req.RefreshToken == "" {
		return nil, xerr.RefreshTokenInvalid
	}

	// 验证 Refresh Token 签名和格式
	_, claims, err := l.svcCtx.RefreshJWTHelper.ValidateToken(req.RefreshToken)
	if err != nil {
		l.Infof("Refresh Token 验证失败: %v", err)
		return nil, xerr.RefreshTokenInvalid
	}

	// 校验 TokenType 必须为 refresh
	if !claims.IsRefreshToken() {
		l.Infof("TokenType 不是 refresh: %s", claims.TokenType)
		return nil, xerr.RefreshTokenInvalid
	}

	// 查询 JTI 是否存在（防重放）
	storedToken, err := l.svcCtx.RefreshTokensModel.FindOneByJti(l.ctx, claims.JTI)
	if err != nil {
		if err == model.ErrNotFound {
			l.Infof("JTI 不存在，令牌已撤销: %s", claims.JTI)
			return nil, xerr.RefreshTokenRevoked
		}
		l.Errorf("查询 JTI 失败: %v", err)
		return nil, xerr.RefreshTokenInvalid
	}

	// 校验是否已撤销
	if storedToken.Revoked == 1 {
		l.Infof("Refresh Token 已撤销: jti=%s", claims.JTI)
		return nil, xerr.RefreshTokenRevoked
	}

	// 校验过期时间
	if storedToken.ExpiresAt < time.Now().Unix() {
		l.Infof("Refresh Token 已过期: jti=%s", claims.JTI)
		return nil, xerr.RefreshTokenExpired
	}

	userId, err := strconv.ParseInt(claims.UID, 10, 64)
	if err != nil {
		l.Errorf("解析 UID 失败: %v", err)
		return nil, xerr.RefreshTokenInvalid
	}

	// 轮换：撤销旧 Refresh Token（一次性使用）
	if err := l.svcCtx.RefreshTokensModel.RevokeByJti(l.ctx, claims.JTI); err != nil {
		l.Errorf("撤销旧 Refresh Token 失败: %v", err)
		return nil, xerr.RefreshTokenInvalid
	}

	// 签发新 TokenPair
	newAccessToken, newRefreshToken, err := l.generateTokenPair(userId)
	if err != nil {
		l.Errorf("生成新令牌失败: %v", err)
		return nil, xerr.RefreshTokenInvalid
	}

	l.Infof("令牌刷新成功: userId=%d", userId)

	return &types.RefreshResp{
		Token:        newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

// generateTokenPair 签发双令牌并持久化 Refresh Token
func (l *RefreshLogic) generateTokenPair(userId int64) (string, string, error) {
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
