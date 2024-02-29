package errhdl

import (
	"awesomeProject1/web"
	"net/http"
	"testing"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	server := web.NewHTTPServer(web.ServerWithMiddleware(NewMiddlewareBuilder().AddCode(http.StatusNotFound, []byte(`<html>
	<h1>404 NOT FOUND 我的自定义错误页面</h1>
	</html>`)).Build()))
	server.Start(":8081")
}
