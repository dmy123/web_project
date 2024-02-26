package web

import (
	"fmt"
	"net/http"
	"testing"
)

func TestServer(t *testing.T) {
	var server Server = &HTTPServer{}
	http.ListenAndServe(":8081", server)

	server.Start(":8081")
}

func TestHTTPServer_ServeHTTP(t *testing.T) {
	server := NewHTTPServer()
	server.mdls = []Middleware{
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("第一个before")
				next(ctx)
				fmt.Println("第一个after")
			}
		},
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("第二个before")
				next(ctx)
				fmt.Println("第二个after")
			}
		},
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("第三个中断")

			}
		},
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("看不到")

			}
		},
	}
	server.ServeHTTP(nil, &http.Request{})
}
