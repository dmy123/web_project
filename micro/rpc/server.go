package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"reflect"
)

type Server struct {
	services map[string]reflectionStub
}

func (s *Server) RegisterService(service Service) {
	s.services[service.Name()] = reflectionStub{
		s:     service,
		value: reflect.ValueOf(service),
	}
}

func NewServer() *Server {
	return &Server{
		services: make(map[string]reflectionStub, 16),
	}
}

func (s *Server) Start(network, addr string) error {
	listener, err := net.Listen(network, addr)
	if err != nil {
		// 比较常见的就是端口被占用
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go func() {
			if er := s.handleConn(conn); er != nil {
				_ = conn.Close()
			}
		}()
	}
}

// 我们可以认为，一个请求包含两部分
// 1. 长度字段：用八个字节表示
// 2. 请求数据：
// 响应也是这个规范
func (s *Server) handleConn(conn net.Conn) error {
	for {
		reqBs, err := ReadMsg(conn)
		if err != nil {
			return err
		}

		r := &Request{}
		err = json.Unmarshal(reqBs, r)
		if err != nil {
			return err
		}

		resp, err := s.Invoke(context.Background(), r)

		if err != nil {
			// 可能是业务error
			return err
		}

		res := EncodeMsg(resp.Data)

		_, err = conn.Write(res)
		if err != nil {
			return err
		}
	}
}

func (s *Server) Invoke(ctx context.Context, r *Request) (*Response, error) {
	rs, ok := s.services[r.ServiceName]
	if !ok {
		return nil, fmt.Errorf("unknown service: %s", r.ServiceName)
	}

	res, err := rs.invoke(ctx, r.MethodName, r.Args)
	if err != nil {
		return nil, err
	}

	return &Response{Data: res}, nil
}

type reflectionStub struct {
	s     Service
	value reflect.Value
}

func (s *reflectionStub) invoke(ctx context.Context, methodName string, data []byte) ([]byte, error) {
	// 反射找到方法，并执行调用
	//val := reflect.ValueOf(service)
	method := s.value.MethodByName(methodName)
	in := make([]reflect.Value, 2)
	// context 如何传？
	in[0] = reflect.ValueOf(context.Background())
	inReq := reflect.New(method.Type().In(1).Elem())
	err := json.Unmarshal(data, inReq.Interface())
	if err != nil {
		return nil, err
	}
	in[1] = inReq
	//results := val.MethodByName(r.MethodName).Call(in)
	results := method.Call(in)
	if results[1].Interface() != nil {
		return nil, results[1].Interface().(error)
	}
	res, err := json.Marshal(results[0].Interface())
	if err != nil {
		return nil, err
	}
	return res, nil
}
