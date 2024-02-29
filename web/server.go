package web

import (
	"fmt"
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

// option模式
type HTTPServerOption func(server *HTTPServer)

type HTTPServer struct {
	router

	mdls []Middleware

	log func(msg string, args ...any)
}

func NewHTTPServer(opts ...HTTPServerOption) *HTTPServer {
	res := &HTTPServer{router: newRouter(), log: func(msg string, args ...any) {
		fmt.Printf(msg, args...)
	}}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

// NewHTTPServerV1缺乏拓展性
//func NewHTTPServerV1(mdls ...Middleware) *HTTPServer {
//	res := &HTTPServer{router: newRouter(), mdls: mdls}
//	return res
//}

func ServerWithMiddleware(mdls ...Middleware) HTTPServerOption {
	return func(server *HTTPServer) {
		server.mdls = mdls
	}
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

	root := h.serve
	for i := len(h.mdls) - 1; i >= 0; i-- {
		root = h.mdls[i](root)
	}

	var m Middleware = func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			next(ctx)
			h.flashResp(ctx)
		}
	}
	m(root)
	root(ctx)

	//h.serve(ctx)
	//ctx.Resp.Write(ctx.RespData)
	//if ctx.RespStatusCode != 0{
	//	ctx.Resp.WriteHeader(ctx.RespStatusCode)
	//}

}
func (h *HTTPServer) flashResp(ctx *Context) {
	if ctx.RespStatusCode != 0 {
		ctx.Resp.WriteHeader(ctx.RespStatusCode)
	}
	n, err := ctx.Resp.Write(ctx.RespData)
	if err != nil || n != len(ctx.RespData) {
		//log.Fatalln("写入响应失败")
		h.log("写入响应失败%v", err)
	}
}

func (h *HTTPServer) serve(ctx *Context) {

	// 查找路由，执行命中业务逻辑
	n, ok := h.findRoute(ctx.Req.Method, ctx.Req.URL.Path)
	if !ok || n.n.handler == nil {
		ctx.RespStatusCode = 404
		ctx.RespData = []byte("NOT FOUND ERR")
		//ctx.Resp.WriteHeader(404)
		//ctx.Resp.Write([]byte("NOT FOUND ERR"))
		return
	}
	ctx.pathParams = n.pathParams
	ctx.MatchedRoute = n.n.route
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
