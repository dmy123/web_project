package web

import "net/http"

type Context struct {
	Req        *http.Request
	Resp       http.ResponseWriter
	pathParams map[string]string
}
