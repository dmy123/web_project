package web

import (
	"net"
	"net/http"
)

type HandleFunc func(ctx *Context)

type Server interface {
	http.Handler
	Start(addr string) error

	// 添加路由注册功能
	addRoute(method string, path string, handleFunc HandleFunc)
}

var _ Server = &HTTPServer{}

type HTTPServer struct {
	router
}

func NewHTTPServer() *HTTPServer {
	return &HTTPServer{router: newRouter()}
}

//func (h *HTTPServer) addRoute(method string, path string, handleFunc HandleFunc) {
//	//TODO implement me
//	panic("implement me")
//}

func (h *HTTPServer) Get(path string, handleFunc HandleFunc) {
	h.addRoute(http.MethodGet, path, handleFunc)
}

func (h *HTTPServer) Post(path string, handleFunc HandleFunc) {
	h.addRoute(http.MethodPost, path, handleFunc)
}

// ServeHTTP 处理请求的入口
func (h *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := &Context{
		Req:        request,
		Resp:       writer,
		pathParams: map[string]string{},
	}
	h.serve(ctx)
}

func (h *HTTPServer) serve(ctx *Context) {

	// 查找路由，执行命中业务逻辑
	n, ok := h.findRoute(ctx.Req.Method, ctx.Req.URL.Path)
	if !ok || n.n.handler == nil {
		ctx.Resp.WriteHeader(404)
		ctx.Resp.Write([]byte("NOT FOUND ERR"))
		return
	}
	ctx.pathParams = n.pathParams
	n.n.handler(ctx)
}

func (h *HTTPServer) Start(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	// 此处，用户可注册after start回调；往admin注册实例；执行业务前置条件

	return http.Serve(l, h)
}
