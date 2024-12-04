package rpc

import "context"

type Service interface {
	Name() string
}

type Proxy interface {
	Invoke(ctx context.Context, req *Request) (*Response, error)
}

type Request struct {
	ServiceName string
	MethodName  string
	Args        []byte
	//Args        any
	//Args        []any
}

type Response struct {
	Data []byte
}
