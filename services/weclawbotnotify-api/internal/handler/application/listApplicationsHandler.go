// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package application

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"weclawbotnotify/services/weclawbotnotify-api/internal/logic/application"
	"weclawbotnotify/services/weclawbotnotify-api/internal/svc"
)

func ListApplicationsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := application.NewListApplicationsLogic(r.Context(), svcCtx)
		resp, err := l.ListApplications()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
