package net

import (
	"awesomeProject1/micro/net/mocks"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func Test_handleConn(t *testing.T) {
	type args struct {
		mock func(ctrl *gomock.Controller) net.Conn
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "read error",
			args: args{
				mock: func(ctrl *gomock.Controller) net.Conn {
					conn := mocks.NewMockConn(ctrl)
					conn.EXPECT().Read(gomock.Any()).Return(0, errors.New("read error"))
					return conn
				},
			},
			wantErr: errors.New("read error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			s := &Server{}
			err := s.handleConn(tt.args.mock(ctrl))
			if !assert.Equal(t, tt.wantErr, err) {
				t.Errorf("handleConn() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
