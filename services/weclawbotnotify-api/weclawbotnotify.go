// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package main

import (
	"context"

	"flag"
	"fmt"
	"net/http"

	"weclawbotnotify/pkg/result"
	"weclawbotnotify/services/weclawbotnotify-api/internal/config"
	"weclawbotnotify/services/weclawbotnotify-api/internal/handler"
	"weclawbotnotify/services/weclawbotnotify-api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var configFile = flag.String("f", "etc/weclawbotnotify-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	httpHandler(server)

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

func httpHandler(server *rest.Server) {
	httpx.SetOkHandler(result.OkHandler)
	httpx.SetErrorHandlerCtx(func(_ context.Context, err error) (int, any) {
		return result.ErrorHandler(context.Background(), err)
	})

	server.Use(func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			next(w, r)
		}
	})
}
