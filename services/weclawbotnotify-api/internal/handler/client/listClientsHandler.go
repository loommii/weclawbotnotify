// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package client

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"weclawbotnotify/services/weclawbotnotify-api/internal/logic/client"
	"weclawbotnotify/services/weclawbotnotify-api/internal/svc"
)

func ListClientsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := client.NewListClientsLogic(r.Context(), svcCtx)
		resp, err := l.ListClients()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
