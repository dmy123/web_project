package recover

import (
	"awesomeProject1/web"
	"fmt"
	"testing"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	builder := NewMiddlewareBuilder()
	builder.StatusCode = 500
	builder.Data = []byte("panic happen")
	builder.LogFunc = func(ctx *web.Context) {
		fmt.Printf("panic 路径：%s", ctx.Req.URL.String())
	}
	server := web.NewHTTPServer(web.ServerWithMiddleware(builder.Build()))
	server.Get("/user", func(ctx *web.Context) {
		ctx.RespJSON(500, "\"发生panic了\"")
		panic("发生panic了")
	})
	server.Start(":8081")
}
