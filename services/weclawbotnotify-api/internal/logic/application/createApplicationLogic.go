package application

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"weclawbotnotify/pkg/jwtx"
	pkgmw "weclawbotnotify/pkg/middleware"
	"weclawbotnotify/pkg/xerr"
	"weclawbotnotify/services/weclawbotnotify-api/internal/model"
	"weclawbotnotify/services/weclawbotnotify-api/internal/svc"
	"weclawbotnotify/services/weclawbotnotify-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateApplicationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateApplicationLogic {
	return &CreateApplicationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CreateApplication 创建应用：参数校验 → 获取用户ID → 生成Token → 入库 → 返回应用信息
// Token 为随机字符串，仅在创建时返回一次明文，用户必须妥善保存
func (l *CreateApplicationLogic) CreateApplication(req *types.CreateApplicationReq) (resp *types.CreateApplicationResp, err error) {
	l.Infof("创建应用请求: name=%s", req.Name)

	// 1. 参数校验：应用名称不能为空
	if req.Name == "" {
		l.Errorf("参数校验失败: 应用名称为空")
		return nil, xerr.ApplicationParamNameEmpty
	}

	// 2. 从 JWT Context 中获取当前登录用户 ID
	userId, err := jwtx.GetUserIdFromContext(l.ctx, pkgmw.ClaimsKey)
	if err != nil {
		l.Errorf("获取用户ID失败: %v", err)
		return nil, xerr.RequestParamError
	}

	// 3. 生成应用 Token（格式：app_ + 32位随机hex字符串）
	token, err := generateAppToken()
	if err != nil {
		l.Errorf("生成应用Token失败: %v", err)
		return nil, xerr.ApplicationTokenFailed
	}

	// 4. 插入数据库，Token 明文存储
	result, err := l.svcCtx.ApplicationsModel.Insert(l.ctx, &model.Applications{
		UserId:      userId,
		Token:       token,
		Name:        req.Name,
		Description: req.Description,
		Status:      model.AppStatusActive,
		CreatedAt:   time.Now().Unix(),
	})
	if err != nil {
		l.Errorf("创建应用失败: %v", err)
		return nil, xerr.ApplicationInsertFailed
	}

	// 5. 获取自增主键 ID
	appId, err := result.LastInsertId()
	if err != nil {
		l.Errorf("获取应用ID失败: %v", err)
		return nil, xerr.ApplicationGetIdFailed
	}

	l.Infof("应用创建成功: appId=%d, name=%s, userId=%d", appId, req.Name, userId)

	// 6. 返回应用信息（Token 仅此一次返回给用户）
	return &types.CreateApplicationResp{
		Id:        appId,
		Token:     token,
		Name:      req.Name,
		CreatedAt: fmt.Sprintf("%d", time.Now().Unix()),
	}, nil
}

// generateAppToken 生成应用 Token，格式为 app_ + 32位随机hex字符串
// 使用 crypto/rand 保证随机性，示例：app_a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6
func generateAppToken() (string, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return "app_" + hex.EncodeToString(bytes), nil
}
