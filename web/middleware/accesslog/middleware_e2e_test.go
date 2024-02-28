//go:build e2e

package accesslog

import (
	"awesomeProject1/web"
	"testing"
)

func TestMiddlewareBuilder_BuildE2E(t *testing.T) {
	server := web.NewHTTPServer(web.ServerWithMiddleware(NewMiddlewareBuilder().Build()))
	server.Get("/a/b/*", func(ctx *web.Context) {
		ctx.Resp.Write([]byte("hello, it's me"))
	})
	server.Start(":8081")
}
