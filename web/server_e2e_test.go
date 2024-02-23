//go:build e2e

package web

import (
	"fmt"
	"net/http"
	"testing"
)

func Test_Server(t *testing.T) {
	h := NewHTTPServer()

	//h.addRoute(http.MethodGet, "/user", func(ctx *Context) {
	//	fmt.Sprintf("1st thing")
	//	fmt.Sprintf("2rd thing")
	//})

	handler1 := func(ctx *Context) {
		fmt.Sprintf("1st thing")
	}
	handler2 := func(ctx *Context) {
		fmt.Sprintf("2rd thing")
	}

	h.addRoute(http.MethodGet, "/user", func(ctx *Context) {
		handler1(ctx)
		handler2(ctx)
	})

	h.addRoute(http.MethodGet, "/order/detail", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello, order detail"))
	})
	h.Start(":8081")
}
