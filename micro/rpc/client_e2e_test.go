package rpc

import (
	"awesomeProject1/micro/proto/gen"
	"awesomeProject1/micro/rpc/serialize/proto"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestInitServiceProto(t *testing.T) {
	server := NewServer()
	service := &UserServiceServer{}
	server.RegisterService(service)
	server.RegisterSerializer(&proto.Serializer{})
	go func() {
		err := server.Start("tcp", ":8081")
		t.Log(err)
	}()
	time.Sleep(time.Second * 3)
	usClient := &UserService{}
	c, err := NewClient(":8081", ClientWithSerializer(&proto.Serializer{}))
	require.NoError(t, err)
	err = c.InitService(usClient)
	require.NoError(t, err)

	testCases := []struct {
		name     string
		mock     func()
		wantErr  error
		wantResp *GetByIDResp
	}{
		{
			name: "no error",
			mock: func() {
				service.Msg = "hello world"
				service.Err = nil
			},
			wantResp: &GetByIDResp{
				Msg: "hello world",
			},
		},
		{
			name: "error",
			mock: func() {
				service.Msg = ""
				service.Err = errors.New("mock error")
			},
			wantResp: &GetByIDResp{},
			wantErr:  errors.New("mock error"),
		},
		{
			name: "both",
			mock: func() {
				service.Msg = "hello world"
				service.Err = errors.New("mock error")
			},
			wantResp: &GetByIDResp{
				Msg: "hello world",
			},
			wantErr: errors.New("mock error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()
			resp, er := usClient.GetByIdProto(context.Background(), &gen.GetByIdReq{
				Id: 123,
			})
			assert.Equal(t, tc.wantErr, er)
			if resp != nil && resp.User != nil {
				assert.Equal(t, tc.wantResp.Msg, resp.User.Name)
			}
		})
	}
}

func TestInitClientProxy(t *testing.T) {
	server := NewServer()
	service := &UserServiceServer{}
	server.RegisterService(service)
	go func() {
		err := server.Start("tcp", ":8081")
		t.Log(err)
	}()
	time.Sleep(time.Second * 3)
	usClient := &UserService{}
	c, err := NewClient(":8081")
	require.NoError(t, err)
	err = c.InitService(usClient)
	require.NoError(t, err)

	//ctx := context.WithValue(context.Background(), "oneway", true)

	testCases := []struct {
		name     string
		mock     func()
		wantErr  error
		wantResp *GetByIDResp
	}{
		{
			name: "no error",
			mock: func() {
				service.Msg = "hello world"
				service.Err = nil
			},
			wantResp: &GetByIDResp{
				Msg: "hello world",
			},
		},
		{
			name: "error",
			mock: func() {
				service.Msg = ""
				service.Err = errors.New("mock error")
			},
			wantResp: &GetByIDResp{},
			wantErr:  errors.New("mock error"),
		},
		{
			name: "both",
			mock: func() {
				service.Msg = "hello world"
				service.Err = errors.New("mock error")
			},
			wantResp: &GetByIDResp{
				Msg: "hello world",
			},
			wantErr: errors.New("mock error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()
			//resp, er  := usClient.GetByID(CtxWithOneway(context.Background()), &GetByIDReq{})
			// 不处理resp
			resp, er := usClient.GetByID(context.Background(), &GetByIDReq{
				Id: 123,
			})
			assert.Equal(t, tc.wantErr, er)
			assert.Equal(t, tc.wantResp, resp)
		})
	}
}

func TestOneWay(t *testing.T) {
	server := NewServer()
	service := &UserServiceServer{}
	server.RegisterService(service)
	go func() {
		err := server.Start("tcp", ":8081")
		t.Log(err)
	}()
	time.Sleep(time.Second * 3)
	usClient := &UserService{}
	c, err := NewClient(":8081")
	require.NoError(t, err)
	err = c.InitService(usClient)
	require.NoError(t, err)

	testCases := []struct {
		name     string
		mock     func()
		wantErr  error
		wantResp *GetByIDResp
	}{
		{
			name: "oneway",
			mock: func() {
				service.Msg = "hello world"
				service.Err = errors.New("mock error")
			},
			wantResp: &GetByIDResp{},
			wantErr:  errors.New("micro: oneway, no need to deal resp"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()
			ctx := CtxWithOneway(context.Background())
			resp, er := usClient.GetByID(ctx, &GetByIDReq{
				Id: 123,
			})
			assert.Equal(t, tc.wantErr, er)
			assert.Equal(t, tc.wantResp, resp)
			time.Sleep(time.Second * 3)
			assert.Equal(t, "hello world", service.Msg)
		})
	}
}
