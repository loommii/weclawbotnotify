package application

import (
	"context"
	"fmt"
	"math"

	"weclawbotnotify/pkg/jwtx"
	pkgmw "weclawbotnotify/pkg/middleware"
	"weclawbotnotify/pkg/xerr"
	"weclawbotnotify/services/weclawbotnotify-api/internal/svc"
	"weclawbotnotify/services/weclawbotnotify-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListApplicationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListApplicationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListApplicationsLogic {
	return &ListApplicationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListApplicationsLogic) ListApplications(req *types.ListApplicationsReq) (resp *types.ListApplicationsResp, err error) {
	userId, err := jwtx.GetUserIdFromContext(l.ctx, pkgmw.ClaimsKey)
	if err != nil {
		l.Errorf("获取用户ID失败: %v", err)
		return nil, xerr.RequestParamError
	}

	page := req.Page
	pageSize := req.PageSize
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	total, err := l.svcCtx.ApplicationsModel.CountByUserId(l.ctx, userId)
	if err != nil {
		l.Errorf("查询应用总数失败: userId=%d, err=%v", userId, err)
		return nil, xerr.ApplicationQueryFailed
	}

	totalPages := 0
	if total > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(pageSize)))
	}

	offset := (page - 1) * pageSize
	applications, err := l.svcCtx.ApplicationsModel.FindByUserIdWithPagination(l.ctx, userId, offset, pageSize)
	if err != nil {
		l.Errorf("查询应用列表失败: userId=%d, err=%v", userId, err)
		return nil, xerr.ApplicationQueryFailed
	}

	list := make([]types.ApplicationInfo, 0, len(applications))
	for _, app := range applications {
		lastUsedAt := ""
		if app.LastUsedAt.Valid {
			lastUsedAt = formatTimestamp(app.LastUsedAt.Int64)
		}
		list = append(list, types.ApplicationInfo{
			Id:          app.Id,
			Name:        app.Name,
			Description: app.Description,
			CreatedAt:   formatTimestamp(app.CreatedAt),
			LastUsedAt:  lastUsedAt,
		})
	}

	l.Infof("查询应用列表成功: userId=%d, page=%d, pageSize=%d, total=%d", userId, page, pageSize, total)

	return &types.ListApplicationsResp{
		PageInfo: types.PageInfo{
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: totalPages,
		},
		List: list,
	}, nil
}

func formatTimestamp(ts int64) string {
	if ts == 0 {
		return ""
	}
	return fmt.Sprintf("%d", ts)
}
