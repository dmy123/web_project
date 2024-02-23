package web

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

func TestRouter_AddRoute(t *testing.T) {
	// 1、构造路由树
	// 2、验证路由树
	testRoutes := []struct {
		method string
		path   string
	}{
		{
			method: http.MethodGet,
			path:   "/user/home",
		},
		{
			method: http.MethodGet,
			path:   "/",
		},
		{
			method: http.MethodGet,
			path:   "/user",
		},
		{
			method: http.MethodGet,
			path:   "/order/detail",
		},
		{
			method: http.MethodPost,
			path:   "/order/create",
		},
		{
			method: http.MethodPost,
			path:   "/login",
		},
		{
			method: http.MethodGet,
			path:   "/order/*",
		},
		{
			method: http.MethodGet,
			path:   "/*",
		},
		{
			method: http.MethodGet,
			path:   "/*/*",
		},
		{
			method: http.MethodGet,
			path:   "/param/:id",
		},
	}
	r := newRouter()
	var mockHandler HandleFunc = func(ctx *Context) {}
	for _, route := range testRoutes {
		r.addRoute(route.method, route.path, mockHandler)
	}
	wantRouter := &router{
		trees: map[string]*node{
			http.MethodGet: {
				path: "/",
				children: map[string]*node{
					"user": {
						path: "user",
						children: map[string]*node{
							"home": {
								path:    "home",
								handler: mockHandler,
							},
						},
						handler: mockHandler,
					},
					"order": {
						path: "order",
						children: map[string]*node{
							"detail": {
								path:    "detail",
								handler: mockHandler,
							},
						},
						starChild: &node{
							path:    "*",
							handler: mockHandler,
						},
					},
					"param": {
						path: "param",
						paramChild: &node{
							path:    ":id",
							handler: mockHandler,
						},
					},
				},
				handler: mockHandler,
				starChild: &node{
					path:    "*",
					handler: mockHandler,
					starChild: &node{
						path:    "*",
						handler: mockHandler,
					},
				},
			},
			http.MethodPost: {
				path: "/",
				children: map[string]*node{
					"login": {
						path:    "login",
						handler: mockHandler,
					},
					"order": {
						path: "order",
						children: map[string]*node{
							"create": {
								path:    "create",
								handler: mockHandler,
							},
						},
						//starChild: &node{
						//	path:    "*",
						//	handler: mockHandler,
						//},
					},
				},
			},
		},
	}
	msg, ok := wantRouter.equal(r)
	assert.True(t, ok, msg)

	r = newRouter()
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "", mockHandler)
	}, "web:路径不能为空字符串")
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/a/b/", mockHandler)
	}, "")
	assert.PanicsWithValue(t, "web: 非法路由，已有通配符路由。不允许同时注册通配符路由和参数路由 [:id]", func() {
		r.addRoute(http.MethodGet, "/a/*", mockHandler)
		r.addRoute(http.MethodGet, "/a/:id", mockHandler)
	})
	assert.PanicsWithValue(t, "web: 非法路由，已有路径参数路由。不允许同时注册通配符路由和参数路由 [*]", func() {
		r.addRoute(http.MethodGet, "/a/b/:id", mockHandler)
		r.addRoute(http.MethodGet, "/a/b/*", mockHandler)
	})
}

// 定义比较node的方法
func (r *router) equal(y router) (string, bool) {
	for k, v := range r.trees {
		yv, ok := y.trees[k]
		if !ok {
			return fmt.Sprintf("目标 router 里面没有方法 %s 的路由树", k), false
		}
		str, res := v.equal(yv)
		if !res {
			return k + "-" + str, false
		}
	}
	return "", true
}

func (m *matchInfo) equal(y *matchInfo) (string, bool) {
	nr, ok := m.n.equal(y.n)
	if !ok {
		return nr, ok
	}
	for k, v := range m.pathParams {
		val, exist := y.pathParams[k]
		if !exist {
			return "web: 参数不存在", false
		}
		if v != val {
			return "web: 参数值不匹配", false
		}
	}
	return "", true
}

// 定义比较node的方法
func (n *node) equal(y *node) (string, bool) {
	if y == nil {
		return "", false
	}

	// 比较path
	if n.path != y.path {
		return fmt.Sprintf("%s 节点 path 不相等 x %s, y %s", n.path, n.path, y.path), false
	}

	if n.starChild != nil {
		msg, ok := n.starChild.equal(y.starChild)
		if !ok {
			return msg, ok
		}
	}

	// 比较handler
	nhv := reflect.ValueOf(n.handler)
	yhv := reflect.ValueOf(y.handler)
	if nhv != yhv {
		return fmt.Sprintf("%s 节点 handler 不相等 x %s, y %s", n.path, nhv.Type().String(), yhv.Type().String()), false
	}
	// 比较children
	if len(n.children) != len(y.children) {
		return fmt.Sprintf("%s 子节点长度不等", n.path), false
	}

	for k, v := range n.children {
		yv, ok := y.children[k]
		if !ok {
			return fmt.Sprintf("%s 目标节点缺少子节点 %s", n.path, k), false
		}
		str, res := v.equal(yv)
		if !res {
			return n.path + "-" + str, false
		}
	}
	return "", true
}

func Test_router_findRoute(t *testing.T) {
	testRoutes := []struct {
		method string
		path   string
	}{
		{
			method: http.MethodGet,
			path:   "/",
		},
		{
			method: http.MethodGet,
			path:   "/user",
		},
		{
			method: http.MethodPost,
			path:   "/order/create",
		},
		{
			method: http.MethodDelete,
			path:   "/",
		},
		{
			method: http.MethodPost,
			path:   "/order/*",
		},
		{
			method: http.MethodGet,
			path:   "/param/:id",
		},
	}
	r := newRouter()
	var mockHandler HandleFunc = func(ctx *Context) {}
	for _, route := range testRoutes {
		r.addRoute(route.method, route.path, mockHandler)
	}

	testCases := []struct {
		name          string
		method        string
		path          string
		wantMatchInfo *matchInfo
		wantFound     bool
	}{
		{
			name:      "root",
			method:    http.MethodDelete,
			path:      "/",
			wantFound: true,
			wantMatchInfo: &matchInfo{n: &node{
				path:    "/",
				handler: mockHandler,
			}},
		},
		{
			name:   "method not found",
			method: http.MethodHead,
		},
		{
			name:      "two layer",
			method:    http.MethodPost,
			path:      "/order/create",
			wantFound: true,
			wantMatchInfo: &matchInfo{n: &node{
				path:    "create",
				handler: mockHandler,
			}},
		},
		{
			name:      "order start",
			method:    http.MethodPost,
			path:      "/order/abc",
			wantFound: true,
			wantMatchInfo: &matchInfo{n: &node{
				path:    "*",
				handler: mockHandler,
			}},
		},
		{
			name:      "no handler",
			method:    http.MethodPost,
			path:      "/order",
			wantFound: true,
			wantMatchInfo: &matchInfo{n: &node{
				path: "order",
				children: map[string]*node{
					"create": {
						path:    "create",
						handler: mockHandler,
					},
				},
				starChild: &node{
					path:    "*",
					handler: mockHandler,
				},
			}},
		},
		{
			name:   "path not found",
			method: http.MethodGet,
			path:   "/abc",
		},
		{
			// 命中 /param/:id
			name:      ":id",
			method:    http.MethodGet,
			path:      "/param/123",
			wantFound: true,
			wantMatchInfo: &matchInfo{n: &node{
				path:    ":id",
				handler: mockHandler,
			},
				pathParams: map[string]string{"id": "123"},
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			res, ok := r.findRoute(tt.method, tt.path)
			assert.Equal(t, tt.wantFound, ok)
			if !ok {
				return
			}
			_, ok = res.equal(tt.wantMatchInfo)
			assert.Equal(t, ok, true)
			//tt.wantMatchInfo.equal(res)
		})
	}
}
