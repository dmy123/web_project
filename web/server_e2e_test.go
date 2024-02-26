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

	h.addRoute(http.MethodPost, "/form", func(ctx *Context) {
		ctx.Req.ParseForm()
		ctx.Resp.Write([]byte(fmt.Sprintf("hello, %s", ctx.Req.URL.Path)))
	})

	h.addRoute(http.MethodPost, "/values/:id", func(ctx *Context) {
		ctx.Req.ParseForm()
		id, err := ctx.PathValue("id").ToInt64()
		if err != nil {
			ctx.Resp.WriteHeader(400)
			ctx.Resp.Write([]byte("id 输入不对"))
			return
		}
		ctx.Resp.Write([]byte(fmt.Sprintf("hello, %s, id = %v", ctx.Req.URL.Path, id)))
	})

	type User struct {
		Name string `json:"name"`
	}

	h.addRoute(http.MethodGet, "/user/123", func(ctx *Context) {
		ctx.RespJSON(202, User{
			Name: "tom",
		})
	})

	h.Start(":8081")
}
