package rpc

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

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
	err := InitClientProxy(":8081", usClient)
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
			resp, er := usClient.GetByID(context.Background(), &GetByIDReq{
				Id: 123,
			})
			assert.Equal(t, tc.wantErr, er)
			assert.Equal(t, tc.wantResp, resp)
		})
	}
}
