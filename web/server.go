package web

import (
	"net"
	"net/http"
)

type HandleFunc func(ctx Context)

type Server interface {
	http.Handler
	Start(addr string) error

	// 添加路由注册功能
	AddRoute(method string, path string, handleFunc HandleFunc)
}

var _ Server = &HTTPServer{}

type HTTPServer struct {
}

func (h *HTTPServer) AddRoute(method string, path string, handleFunc HandleFunc) {
	//TODO implement me
	panic("implement me")
}

// ServeHTTP 处理请求的入口
func (h *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (h *HTTPServer) Start(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	// 此处，用户可注册after start回调；往admin注册实例；执行业务前置条件

	return http.Serve(l, h)
}
