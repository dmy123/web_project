package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_setFuncField(t *testing.T) {
	type args struct {
		service Service
		mock    func(ctrl *gomock.Controller) Proxy
		//proxy   Proxy
	}
	arg, err := json.Marshal(&GetByIDReq{Id: 123})
	assert.NoError(t, err)
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "nil",
			args: args{
				service: nil,
				mock: func(ctrl *gomock.Controller) Proxy {
					return NewMockProxy(ctrl)
				},
			},
			wantErr: errors.New("rpc：不支持nil"),
		},
		//{
		//	name: "not pointer",
		//	args: args{
		//		service: new(*int),
		//	},
		//	wantErr: errors.New("rpc：只支持指向结构体的一级指针"),
		//},
		{
			name: "not pointer",
			args: args{
				service: UserService{},
				mock: func(ctrl *gomock.Controller) Proxy {
					return NewMockProxy(ctrl)
				},
			},
			wantErr: errors.New("rpc：只支持指向结构体的一级指针"),
		},
		{
			name: "pointer",
			args: args{
				service: &UserService{},
				mock: func(ctrl *gomock.Controller) Proxy {
					p := NewMockProxy(ctrl)
					p.EXPECT().Invoke(gomock.Any(), &Request{
						ServiceName: "user-service",
						MethodName:  "GetByID",
						Args:        arg,
					}).Return(&Response{}, nil)
					return p
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			err := setFuncField(tt.args.service, tt.args.mock(ctrl))
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			resp, err := tt.args.service.(*UserService).GetByID(context.Background(), &GetByIDReq{Id: 123})
			assert.Equal(t, tt.wantErr, err)
			t.Log(resp)
		})
	}
}

//type mockProxy struct {
//}
//
//func (m mockProxy) Invoke(ctx context.Context, req *Request) (*Response, error) {
//	//TODO implement me
//	panic("implement me")
//}
