package accesslog

import (
	"awesomeProject1/web"
	"fmt"
	"net/http"
	"testing"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	server := web.NewHTTPServer(web.ServerWithMiddleware(NewMiddlewareBuilder().Build()))
	server.Post("/a/b/*", func(ctx *web.Context) {
		fmt.Println("hello, it's me")
	})
	req, err := http.NewRequest(http.MethodPost, "/a/b/c", nil)
	if err != nil {
		t.Fatal(err)
	}
	server.ServeHTTP(nil, req)
}
