package rpc

import (
	"awesomeProject1/micro/rpc/message"
	"awesomeProject1/micro/rpc/serialize"
	json2 "awesomeProject1/micro/rpc/serialize/json"
	"context"
	"errors"
	"github.com/silenceper/pool"
	"net"
	"reflect"
	"time"
)

// InitService 要为GetById 之类的函数类型的字段赋值
func (c *Client) InitService(service Service) error {
	return setFuncField(service, c, c.serializer)
}

func setFuncField(service Service, p Proxy, s serialize.Serializer) error {
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

				reqData, err := s.Encode(args[1].Interface())
				if err != nil {
					return []reflect.Value{retVal, reflect.ValueOf(err)}
				}
				meta := make(map[string]string)
				if isOneway(ctx) {
					meta["one-way"] = "true"
				}
				req := &message.Request{
					ServiceName: service.Name(),
					MethodName:  fieldTyp.Name,
					Data:        reqData,
					Serializer:  s.Code(),
					Meta:        meta,
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
					err = s.Decode(resp.Data, retVal.Interface())
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
	pool       pool.Pool
	serializer serialize.Serializer
}

type ClientOpt func(c *Client)

func ClientWithSerializer(sl serialize.Serializer) ClientOpt {
	return func(c *Client) {
		c.serializer = sl
	}
}

func NewClient(addr string, opt ...ClientOpt) (*Client, error) {
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
	res := &Client{
		//addr: addr,
		pool:       p,
		serializer: &json2.Serializer{},
	}
	for _, opt := range opt {
		opt(res)
	}
	return res, nil
}

func (c *Client) Invoke(ctx context.Context, req *message.Request) (*message.Response, error) {
	//data, err := json.Marshal(req)
	//if err != nil {
	//	return nil, err
	//}
	data := message.EncodeReq(req)
	// 发请求
	//conn, err := net.DialTimeout("tcp", c.addr, time.Second*3)
	resp, err := c.send(ctx, data)
	//return &message.Response{
	//	Data: resp,
	//}, nil
	return message.DecodeResp(resp), err
}

func (c *Client) send(ctx context.Context, data []byte) ([]byte, error) {
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

	if isOneway(ctx) {
		return nil, errors.New("micro: oneway, no need to deal resp")
	}

	return ReadMsg(conn)
}
