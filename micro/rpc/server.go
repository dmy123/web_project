package rpc

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"reflect"
)

type Server struct {
	services map[string]Service
}

func (s *Server) RegisterService(service Service) {
	s.services[service.Name()] = service
}

func NewServer() *Server {
	return &Server{
		services: make(map[string]Service, 16),
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
		// lenBs 是长度字段的字节表示
		lenBs := make([]byte, numOfLengthBytes)
		_, err := conn.Read(lenBs)
		if err != nil {
			return err
		}

		// 我消息有多长？
		length := binary.BigEndian.Uint64(lenBs)

		reqBs := make([]byte, length)
		_, err = conn.Read(reqBs)
		if err != nil {
			return err
		}

		respData, err := s.handleMsg(reqBs)
		if err != nil {
			// 可能是业务error
			return err
		}
		respLen := len(respData)

		// 我要在这，构建响应数据
		// data = respLen 的 64 位表示 + respData
		res := make([]byte, respLen+numOfLengthBytes)
		// 第一步：
		// 把长度写进去前八个字节
		binary.BigEndian.PutUint64(res[:numOfLengthBytes], uint64(respLen))
		// 第二步：
		// 写入数据
		copy(res[numOfLengthBytes:], respData)

		_, err = conn.Write(res)
		if err != nil {
			return err
		}
	}
}

func (s *Server) handleMsg(req []byte) ([]byte, error) {
	// 还原调用信息
	r := &Request{}
	err := json.Unmarshal(req, r)
	if err != nil {
		return nil, err
	}
	// 还原了调用信息
	// 发起业务调用
	service, ok := s.services[r.ServiceName]
	if !ok {
		return nil, fmt.Errorf("unknown service: %s", r.ServiceName)
	}
	// 反射找到方法，并执行调用
	val := reflect.ValueOf(service)
	method := val.MethodByName(r.MethodName)
	in := make([]reflect.Value, 2)
	// context 如何传？
	in[0] = reflect.ValueOf(context.Background())
	inReq := reflect.New(method.Type().In(1).Elem())
	err = json.Unmarshal(r.Args, inReq.Interface())
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
	return res, nil
}
