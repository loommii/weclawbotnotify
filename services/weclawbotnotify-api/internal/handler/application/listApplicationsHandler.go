// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package application

import (
	"net/http"

	"weclawbotnotify/services/weclawbotnotify-api/internal/logic/application"
	"weclawbotnotify/services/weclawbotnotify-api/internal/svc"
	"weclawbotnotify/services/weclawbotnotify-api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ListApplicationsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListApplicationsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := application.NewListApplicationsLogic(r.Context(), svcCtx)
		resp, err := l.ListApplications(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
