package web

import (
	"encoding/json"
	"errors"
	"net/http"
)

type Context struct {
	Req        *http.Request
	Resp       http.ResponseWriter
	pathParams map[string]string
}

func (c *Context) BindJson(val any) error {
	if c.Req.Body == nil {
		return errors.New("web: body 为 nil")
	}
	decoder := json.NewDecoder(c.Req.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(val)
}
