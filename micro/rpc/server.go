package rpc

import (
	"awesomeProject1/micro/rpc/message"
	"awesomeProject1/micro/rpc/serialize"
	json2 "awesomeProject1/micro/rpc/serialize/json"
	"context"
	"errors"
	"fmt"
	"net"
	"reflect"
	"strconv"
	"time"
)

type Server struct {
	services    map[string]reflectionStub
	serializers map[uint8]serialize.Serializer
}

func (s *Server) RegisterSerializer(sl serialize.Serializer) {
	s.serializers[sl.Code()] = sl
}

func (s *Server) RegisterService(service Service) {
	s.services[service.Name()] = reflectionStub{
		s:           service,
		value:       reflect.ValueOf(service),
		serializers: s.serializers,
	}
}

func NewServer() *Server {
	server := &Server{
		services:    make(map[string]reflectionStub, 16),
		serializers: make(map[uint8]serialize.Serializer, 4),
	}

	server.RegisterSerializer(&json2.Serializer{})
	return server
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

		//r := &message.Request{}
		//err = json.Unmarshal(reqBs, r)
		r := message.DecodeReq(reqBs)
		if err != nil {
			return err
		}

		ctx := context.Background()
		cancel := func() {}
		if deadlineStr, ok := r.Meta["Deadline"]; ok {
			if deadline, er := strconv.ParseInt(deadlineStr, 10, 64); er == nil {
				ctx, cancel = context.WithDeadline(ctx, time.UnixMilli(deadline))
			}
		}
		oneway, ok := r.Meta["one-way"]
		if ok && oneway == "true" {
			ctx = CtxWithOneway(ctx)
		}

		resp, err := s.Invoke(ctx, r)
		cancel()

		if err != nil {
			// 可能是业务error
			resp.Error = []byte(err.Error())
		}

		resp.CalculateHeadLength()
		resp.CalculateBodyLength()

		//res := EncodeMsg(resp.Data)
		//_, err = conn.Write(res)
		resp.CalculateHeadLength()
		resp.CalculateBodyLength()
		_, err = conn.Write(message.EncodeResp(resp))
		if err != nil {
			return err
		}
	}
}

func (s *Server) Invoke(ctx context.Context, r *message.Request) (*message.Response, error) {
	resp := &message.Response{
		MessageId:  r.MessageId,
		Version:    r.Version,
		Compresser: r.Compresser,
		Serializer: r.Serializer,
	}

	rs, ok := s.services[r.ServiceName]
	if !ok {
		return resp, fmt.Errorf("unknown service: %s", r.ServiceName)
	}

	if isOneway(ctx) {
		go func() {
			rs.invoke(ctx, r)
		}()
		return resp, errors.New("micro: ")
	}

	respData, err := rs.invoke(ctx, r)
	//if isOneway(ctx){
	//	return nil, errors.New("")
	//}
	resp.Data = respData
	if err != nil {
		return resp, err
	}
	return resp, nil
}

type reflectionStub struct {
	s           Service
	value       reflect.Value
	serializers map[uint8]serialize.Serializer
}

func (s *reflectionStub) invoke(ctx context.Context, req *message.Request) ([]byte, error) {
	// 反射找到方法，并执行调用
	//val := reflect.ValueOf(service)
	method := s.value.MethodByName(req.MethodName)
	in := make([]reflect.Value, 2)
	// context 如何传？
	in[0] = reflect.ValueOf(ctx)
	inReq := reflect.New(method.Type().In(1).Elem())
	serializer, ok := s.serializers[req.Serializer]
	if !ok {
		return nil, fmt.Errorf("unknown serializer: %q", req.Serializer)
	}
	err := serializer.Decode(req.Data, inReq.Interface())
	if err != nil {
		return nil, err
	}
	in[1] = inReq
	//results := val.MethodByName(r.MethodName).Call(in)
	results := method.Call(in)
	if results[1].Interface() != nil {
		err = results[1].Interface().(error)
	}

	var res []byte
	if results[0].IsNil() {
		return nil, err
	} else {
		var er error
		res, er = serializer.Encode(results[0].Interface())
		if er != nil {
			return nil, er
		}
	}
	return res, err
}
