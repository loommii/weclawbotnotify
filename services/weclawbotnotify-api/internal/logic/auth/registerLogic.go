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

// Register 用户注册：参数校验 → 单用户限制检查 → 用户名查重 → 创建用户 → 生成JWT
func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	l.Infof("用户注册请求: username=%s", req.Username)

	// 参数校验
	if req.Username == "" || req.Password == "" {
		l.Errorf("参数校验失败: 用户名或密码为空")
		return nil, xerr.RegisterParamEmpty
	}

	// V1 版本仅支持单用户注册：若已有用户存在则拒绝注册
	// TODO: 后续版本开放多用户注册时移除此限制
	count, err := l.svcCtx.UsersModel.Count(l.ctx)
	if err != nil {
		l.Errorf("查询用户数量失败: %v", err)
		return nil, xerr.RegisterQueryFailed
	}
	if count > 0 {
		l.Infof("拒绝注册: 系统已存在用户，V1版本仅支持单用户")
		return nil, xerr.RegisterClosed
	}

	// 检查用户名是否已存在
	_, err = l.svcCtx.UsersModel.FindOneByUsername(l.ctx, req.Username)
	if err == nil {
		l.Infof("用户名已存在: %s", req.Username)
		return nil, xerr.RegisterUsernameExist
	}
	if !isNotFound(err) {
		l.Errorf("查询用户失败: %v", err)
		return nil, xerr.RegisterQueryFailed
	}

	// 密码 bcrypt 哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		l.Errorf("密码哈希失败: %v", err)
		return nil, xerr.RegisterHashFailed
	}

	// 插入用户记录
	result, err := l.svcCtx.UsersModel.Insert(l.ctx, &model.Users{
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
	})
	if err != nil {
		l.Errorf("创建用户失败: %v", err)
		return nil, xerr.RegisterInsertFailed
	}

	// 获取自增 ID
	userId, err := result.LastInsertId()
	if err != nil {
		l.Errorf("获取用户ID失败: %v", err)
		return nil, xerr.RegisterGetIdFailed
	}

	// 生成 JWT 令牌
	token, err := l.svcCtx.JWTHelper.GenerateToken(jwtx.JWTClaims{
		UID:       strconv.FormatInt(userId, 10),
		TokenType: jwtx.Access,
	})
	if err != nil {
		l.Errorf("生成JWT令牌失败: %v", err)
		return nil, xerr.RegisterTokenFailed
	}

	l.Infof("用户注册成功: userId=%d, username=%s", userId, req.Username)

	return &types.RegisterResp{
		Token: token,
		User: types.UserInfo{
			Id:        userId,
			Username:  req.Username,
			CreatedAt: fmt.Sprintf("%d", time.Now().Unix()),
		},
	}, nil
}

// isNotFound 判断是否为数据库记录未找到错误
func isNotFound(err error) bool {
	return err == model.ErrNotFound
}
