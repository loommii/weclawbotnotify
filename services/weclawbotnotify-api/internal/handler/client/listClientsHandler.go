// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package client

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"weclawbotnotify/services/weclawbotnotify-api/internal/logic/client"
	"weclawbotnotify/services/weclawbotnotify-api/internal/svc"
	"weclawbotnotify/services/weclawbotnotify-api/internal/types"
)

func ListClientsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListClientsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := client.NewListClientsLogic(r.Context(), svcCtx)
		resp, err := l.ListClients(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
