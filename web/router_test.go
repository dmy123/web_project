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
					},
				},
				handler: mockHandler,
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
					},
				},
			},
		},
	}
	msg, ok := wantRouter.equal(*r)
	assert.True(t, ok, msg)

	r = newRouter()
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "", mockHandler)
	}, "web:路径不能为空字符串")
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/a/b/", mockHandler)
	}, "")
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

// 定义比较node的方法
func (n *node) equal(y *node) (string, bool) {
	if y == nil {
		return "", false
	}

	// 比较path
	if n.path != y.path {
		return fmt.Sprintf("%s 节点 path 不相等 x %s, y %s", n.path, n.path, y.path), false
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
	}
	r := newRouter()
	var mockHandler HandleFunc = func(ctx *Context) {}
	for _, route := range testRoutes {
		r.addRoute(route.method, route.path, mockHandler)
	}

	testCases := []struct {
		name      string
		method    string
		path      string
		wantNode  *node
		wantFound bool
	}{
		{
			name:      "root",
			method:    http.MethodDelete,
			path:      "/",
			wantFound: true,
			wantNode: &node{
				path:    "/",
				handler: mockHandler,
			},
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
			wantNode: &node{
				path:    "create",
				handler: mockHandler,
			},
		},
		{
			name:      "no handler",
			method:    http.MethodPost,
			path:      "/order",
			wantFound: true,
			wantNode: &node{
				path: "order",
				children: map[string]*node{
					"create": {
						path:    "create",
						handler: mockHandler,
					},
				},
			},
		},
		{
			name:   "path not found",
			method: http.MethodGet,
			path:   "/abc",
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			res, ok := r.findRoute(tt.method, tt.path)
			assert.Equal(t, tt.wantFound, ok)
			if !ok {
				return
			}
			_, ok = res.equal(tt.wantNode)
			assert.Equal(t, ok, true)
			//tt.wantNode.equal(res)
		})
	}
}
