package rpc

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"reflect"
	"time"
)

// InitClientProxy 要为GetById 之类的函数类型的字段赋值
func InitClientProxy(addr string, service Service) error {
	client := NewClient(addr)
	return setFuncField(service, client)
}

func setFuncField(service Service, p Proxy) error {
	if service == nil {
		return errors.New("rpc：不支持nil")
	}
	val := reflect.ValueOf(service)
	typ := val.Type()
	if typ.Kind() != reflect.Pointer || typ.Elem().Kind() != reflect.Struct {
		return errors.New("rpc：只支持指向结构体的一级指针")
	}

	val = val.Elem()
	typ = typ.Elem()

	numField := typ.NumField()
	for i := 0; i < numField; i++ {
		fieldTyp := typ.Field(i)
		fieldVal := val.Field(i)
		if fieldVal.CanSet() {
			// 捕捉本地调用，而后调用set方法篡改了它，改成发起rpc调用
			fn := func(args []reflect.Value) (results []reflect.Value) {

				// args[0]是context
				ctx := args[0].Interface().(context.Context)
				// args[1]是req

				retVal := reflect.New(fieldTyp.Type.Out(0).Elem())

				reqData, err := json.Marshal(args[1].Interface())
				if err != nil {
					return []reflect.Value{retVal, reflect.ValueOf(err)}
				}
				req := &Request{
					ServiceName: service.Name(),
					MethodName:  fieldTyp.Name,
					Args:        reqData,
					//Args:        []any{args[1].Interface()},
					//Args: slice.Map[reflect.Value, any](args, func(idx int, src reflect.Value) any {
					//	return src.Interface()
					//}),
				}
				//var p Proxy
				resp, err := p.Invoke(ctx, req)
				if err != nil {
					return []reflect.Value{retVal, reflect.ValueOf(err)}
				}
				err = json.Unmarshal(resp.Data, retVal.Interface())
				if err != nil {
					return []reflect.Value{retVal, reflect.ValueOf(err)}
				}

				fmt.Println("resp:", resp)
				return []reflect.Value{retVal, reflect.Zero(reflect.TypeOf(new(error)).Elem())}
			}
			fnVal := reflect.MakeFunc(fieldTyp.Type, fn)
			fieldVal.Set(fnVal)
		}
	}

	return nil
}

// 长度字段使用的字节数量
const numOfLengthBytes = 8

type Client struct {
	addr string
}

func NewClient(addr string) *Client {
	return &Client{addr: addr}
}

func (c *Client) Invoke(ctx context.Context, req *Request) (*Response, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	// 发请求
	//conn, err := net.DialTimeout("tcp", c.addr, time.Second*3)
	resp, err := c.Send(data)
	return &Response{
		Data: resp,
	}, nil
}

func (c *Client) Send(data []byte) ([]byte, error) {
	conn, err := net.DialTimeout("tcp", c.addr, time.Second*3)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = conn.Close()
	}()

	reqLen := len(data)

	// 我要在这，构建请求数据
	// data = reqLen 的 64 位表示 + respData
	req := make([]byte, reqLen+numOfLengthBytes)
	// 第一步：
	// 把长度写进去前八个字节
	binary.BigEndian.PutUint64(req[:numOfLengthBytes], uint64(reqLen))
	// 第二步：
	// 写入数据
	copy(req[numOfLengthBytes:], data)

	_, err = conn.Write(req)
	if err != nil {
		return nil, err
	}

	lenBs := make([]byte, numOfLengthBytes)
	_, err = conn.Read(lenBs)
	if err != nil {
		return nil, err
	}

	// 我响应有多长？
	length := binary.BigEndian.Uint64(lenBs)

	respBs := make([]byte, length)
	_, err = conn.Read(respBs)
	if err != nil {
		return nil, err
	}

	return respBs, nil
}
