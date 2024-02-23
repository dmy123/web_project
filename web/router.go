package web

import (
	"fmt"
	"strings"
)

type router struct {
	trees map[string]*node
}

func newRouter() *router {
	return &router{trees: make(map[string]*node)}
}

type node struct {
	path     string
	children map[string]*node
	handler  HandleFunc

	route string
}

func (r *router) addRoute(method string, path string, handler HandleFunc) {
	if path == "" {
		panic("web:路径不能为空字符串")
	}
	if path[0] != '/' {
		panic("web: 路由必须以 / 开头")
	}
	if path != "/" && path[len(path)-1] == '/' {
		panic("web: 路由不能以 / 结尾")
	}

	if r.trees == nil {
		r.trees = make(map[string]*node)
	}

	root, exist := r.trees[method]
	if !exist {
		root = &node{
			path: "/",
		}
		r.trees[method] = root
	}
	if path == "/" {
		root.handler = handler
		return
	}

	segs := strings.Split(path[1:], "/")

	for _, s := range segs {
		if s == "" {
			panic(fmt.Sprintf("web: 非法路由。不允许使用 //a/b, /a//b 之类的路由, [%s]", path))
		}
		root = root.childOrCreate(s)
	}
	if root.handler != nil {
		panic(fmt.Sprintf("web: 路由冲突[%s]", path))
	}
	root.handler = handler
	root.route = path
}

func (r *router) findRoute(method string, path string) (*node, bool) {
	root, exist := r.trees[method]
	if !exist {
		return nil, false
	}
	if path == "/" {
		return root, true
	}
	path = strings.Trim(path, "/")

	segs := strings.Split(path, "/")
	for _, seg := range segs {
		child, exist := root.childOf(seg)
		if !exist {
			return nil, false
		}
		root = child
	}
	return root, true
}

func (n *node) childOf(path string) (*node, bool) {
	if n.children == nil {
		return nil, false
	}
	child, exist := n.children[path]
	if !exist {
		return nil, false
	}
	return child, true
}

func (n *node) childOrCreate(path string) *node {
	if n.children == nil {
		n.children = make(map[string]*node)
	}
	child, exist := n.children[path]
	if !exist {
		child = &node{path: path}
		n.children[path] = child
	}
	return child
}
