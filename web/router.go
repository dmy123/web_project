package web

import (
	"fmt"
	"regexp"
	"strings"
)

type router struct {
	trees map[string]*node
}

func newRouter() router {
	return router{trees: make(map[string]*node)}
}

type node struct {
	path     string
	children map[string]*node
	handler  HandleFunc

	route string
	// 正则匹配
	regexChild *node
	regexExpr  *regexp.Regexp
	// 通配符匹配
	starChild *node
	// 参数路径匹配
	paramChild *node

	// 参数路径or正则的参数名
	paramName string
}

type matchInfo struct {
	n          *node
	pathParams map[string]string
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

func (r *router) findRoute(method string, path string) (*matchInfo, bool) {
	root, exist := r.trees[method]
	if !exist {
		return nil, false
	}
	if path == "/" {
		return &matchInfo{n: root}, true
	}
	path = strings.Trim(path, "/")

	segs := strings.Split(path, "/")
	var pathParams map[string]string
	for _, seg := range segs {
		child, paramChild, exist := root.childOf(seg)
		if !exist {
			if root.path == "*" && root.children == nil && root.regexChild == nil && root.paramChild == nil {
				break
			}
			return nil, false
		}
		if paramChild {
			if pathParams == nil {
				pathParams = make(map[string]string)
			}
			pathParams[child.path[1:]] = seg
		}
		root = child
	}
	return &matchInfo{n: root, pathParams: pathParams}, true
}

// childOf 返回子节点，是否为路径参数，是否存在子节点
// 优先级：静态匹配、正则、参数路径、通配符匹配
func (n *node) childOf(path string) (*node, bool, bool) {

	if n.children == nil {
		if n.regexChild != nil && n.regexChild.regexExpr != nil && n.regexChild.regexExpr.Match([]byte(path)) {
			return n.regexChild, false, true
		}
		if n.paramChild != nil {
			return n.paramChild, true, true
		}
		return n.starChild, false, n.starChild != nil
	}
	child, exist := n.children[path]
	if !exist {
		if n.regexChild != nil && n.regexChild.regexExpr != nil && n.regexChild.regexExpr.Match([]byte(path)) {
			return n.regexChild, false, true
		}
		if n.paramChild != nil {
			return n.paramChild, true, true
		}
		return n.starChild, false, n.starChild != nil
	}
	return child, false, true
}

func (n *node) childOrCreate(path string) *node {
	if path == "*" {
		if n.paramChild != nil {
			panic(fmt.Sprintf("web: 非法路由，已有路径参数路由，不允许同时注册通配符路由和参数路由"))
		}
		if n.regexChild != nil {
			panic(fmt.Sprintf("web: 非法路由，已有正则路由，不允许同时注册通配符路由和参数路由"))
		}
	}

	if path[0] == ':' {
		// 判断是否为正则
		i := strings.IndexAny(path, "(")
		if i > 0 {
			// 取出正则，填充paramname
			n.regexChild = &node{path: path, paramName: path[1:i], regexExpr: regexp.MustCompile(path[i+1 : len(path)-1])}
			return n.regexChild
		}
		if n.starChild != nil {
			panic(fmt.Sprintf("web: 非法路由，已有通配符路由。不允许同时注册通配符路由和参数路由 [%s]", path))
		}
		n.paramChild = &node{path: path, paramName: path[1:]}
		return n.paramChild
	}
	if path == "*" {
		if n.paramChild != nil {
			panic(fmt.Sprintf("web: 非法路由，已有路径参数路由。不允许同时注册通配符路由和参数路由 [%s]", path))
		}
		if n.starChild != nil {
			return n.starChild
		}
		n.starChild = &node{path: "*"}
		return n.starChild
	}
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
