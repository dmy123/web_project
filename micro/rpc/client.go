package rpc

import (
	"awesomeProject1/micro/rpc/message"
	"context"
	"encoding/json"
	"errors"
	"github.com/silenceper/pool"
	"net"
	"reflect"
	"time"
)

// InitClientProxy 要为GetById 之类的函数类型的字段赋值
func InitClientProxy(addr string, service Service) error {
	client, err := NewClient(addr)
	if err != nil {
		return err
	}
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
				req := &message.Request{
					ServiceName: service.Name(),
					MethodName:  fieldTyp.Name,
					Data:        reqData,
				}
				req.CalculateHeadLength()
				req.CalculateBodyLength()
				resp, err := p.Invoke(ctx, req)
				if err != nil {
					return []reflect.Value{retVal, reflect.ValueOf(err)}
				}

				var retErr error
				if len(resp.Error) > 0 {
					retErr = errors.New(string(resp.Error))
				}

				if len(resp.Data) > 0 {
					err = json.Unmarshal(resp.Data, retVal.Interface())
					if err != nil {
						return []reflect.Value{retVal, reflect.ValueOf(err)}
					}
				}

				var retErrVal reflect.Value
				if retErr == nil {
					retErrVal = reflect.Zero(reflect.TypeOf(new(error)).Elem())
				} else {
					retErrVal = reflect.ValueOf(retErr)
				}

				return []reflect.Value{retVal, retErrVal}
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
	//addr string
	pool pool.Pool
}

func NewClient(addr string) (*Client, error) {
	p, err := pool.NewChannelPool(&pool.Config{
		InitialCap:  1,
		MaxCap:      30,
		MaxIdle:     10,
		IdleTimeout: time.Minute,
		Factory: func() (interface{}, error) {
			return net.DialTimeout("tcp", addr, time.Second*3)
		},
		Close: func(obj interface{}) error {
			return obj.(net.Conn).Close()
		},
	})
	if err != nil {
		return nil, err
	}
	return &Client{
		//addr: addr,
		pool: p,
	}, nil
}

func (c *Client) Invoke(ctx context.Context, req *message.Request) (*message.Response, error) {
	//data, err := json.Marshal(req)
	//if err != nil {
	//	return nil, err
	//}
	data := message.EncodeReq(req)
	// 发请求
	//conn, err := net.DialTimeout("tcp", c.addr, time.Second*3)
	resp, _ := c.Send(data)
	//return &message.Response{
	//	Data: resp,
	//}, nil
	return message.DecodeResp(resp), nil
}

func (c *Client) Send(data []byte) ([]byte, error) {
	//conn, err := net.DialTimeout("tcp", c.addr, time.Second*3)
	co, err := c.pool.Get()
	if err != nil {
		return nil, err
	}
	conn := co.(net.Conn)
	defer func() {
		_ = conn.Close()
	}()
	//req := EncodeMsg(data)
	//_, err = conn.Write(req)
	_, err = conn.Write(data)
	if err != nil {
		return nil, err
	}

	return ReadMsg(conn)
}