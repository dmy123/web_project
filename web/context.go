package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

type Context struct {
	Req        *http.Request
	Resp       http.ResponseWriter
	pathParams map[string]string
	queryValue url.Values
}

func (c *Context) RespJSON(code int, val any) error {
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}
	c.Resp.WriteHeader(code)
	n, err := c.Resp.Write(data)
	if n != len(data) {
		return errors.New("web： 未写入全部数据")
	}
	return err
}

func (c *Context) FromValue(key string) (val string, err error) {
	err = c.Req.ParseForm()
	if err != nil {
		return "", err
	}
	return c.Req.FormValue(key), nil
}

func (c *Context) BindJson(val any) error {
	if c.Req.Body == nil {
		return errors.New("web: body 为 nil")
	}
	decoder := json.NewDecoder(c.Req.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(val)
}

func (c *Context) QueryValue(key string) (string, error) {
	if c.queryValue == nil {
		c.queryValue = c.Req.URL.Query()
	}
	vals, exist := c.queryValue[key]
	if !exist || len(vals) == 0 {
		return "", errors.New("web: 找不到这个 key")
	}
	return vals[0], nil
}

func (c *Context) PathValue(key string) StringValue {
	val, ok := c.pathParams[key]
	if !ok {
		return StringValue{val: "", err: errors.New("web: key不存在")}
	}
	return StringValue{val: val, err: nil}
}

type StringValue struct {
	val string
	err error
}

func (s StringValue) String() (string, error) {
	return s.val, s.err
}

func (s StringValue) ToInt64() (int64, error) {
	if s.err != nil {
		return 0, s.err
	}
	return strconv.ParseInt(s.val, 10, 64)
}
