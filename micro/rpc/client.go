package rpc

import (
	"awesomeProject1/micro/rpc/compresser"
	"awesomeProject1/micro/rpc/message"
	"awesomeProject1/micro/rpc/serialize"
	json2 "awesomeProject1/micro/rpc/serialize/json"
	"context"
	"errors"
	"github.com/silenceper/pool"
	"net"
	"reflect"
	"strconv"
	"time"
)

// InitService 要为GetById 之类的函数类型的字段赋值
func (c *Client) InitService(service Service) error {
	return setFuncField(service, c, c.serializer, c.compresser)
}

func setFuncField(service Service, p Proxy, s serialize.Serializer, c compresser.Compresser) error {
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

				rd, err := c.Compress(reqData)
				if err != nil {
					return []reflect.Value{retVal, reflect.ValueOf(err)}
				}
				reqData = rd

				meta := make(map[string]string, 2)
				if deadline, ok := ctx.Deadline(); ok {
					meta["Deadline"] = strconv.FormatInt(deadline.UnixMilli(), 10)
				}
				if isOneway(ctx) {
					meta["one-way"] = "true"
				}
				req := &message.Request{
					ServiceName: service.Name(),
					MethodName:  fieldTyp.Name,
					Data:        reqData,
					Serializer:  s.Code(),
					Compresser:  c.Code(),
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
					data, err := c.Uncompress(resp.Data)
					if err != nil {
						return []reflect.Value{retVal, reflect.ValueOf(err)}
					}
					err = s.Decode(data, retVal.Interface())
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
	compresser compresser.Compresser
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
		compresser: &compresser.DoNothingCompresser{},
	}
	for _, opt := range opt {
		opt(res)
	}
	return res, nil
}

func (c *Client) Invoke(ctx context.Context, req *message.Request) (*message.Response, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	ch := make(chan struct{}, 1)
	defer func() { close(ch) }()
	var (
		resp *message.Response
		err  error
	)
	go func() {
		resp, err = c.doInvoke(ctx, req)
		ch <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-ch:
		return resp, err
	}

}

func (c *Client) doInvoke(ctx context.Context, req *message.Request) (*message.Response, error) {
	data := message.EncodeReq(req)

	resp, err := c.send(ctx, data)
	if err != nil {
		return nil, err
	}
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
