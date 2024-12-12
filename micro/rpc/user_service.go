package rpc

import (
	"awesomeProject1/micro/proto/gen"
	"context"
	"log"
	"testing"
	"time"
)

type UserService struct {
	// 用反射来赋值
	// 类型是函数的字段，它不是方法（它不是定义在 UserService 上的方法）
	GetByID      func(ctx context.Context, req *GetByIDReq) (*GetByIDResp, error)
	GetByIdProto func(ctx context.Context, req *gen.GetByIdReq) (*gen.GetByIdResp, error)
}

func (u UserService) Name() string {
	return "user-service"
}

type GetByIDReq struct {
	Id int
}

type GetByIDResp struct {
	Msg string
}

type UserServiceServer struct {
	Err error
	Msg string
}

func (u *UserServiceServer) GetByID(ctx context.Context, req *GetByIDReq) (*GetByIDResp, error) {
	log.Println("req:", req)
	return &GetByIDResp{
		Msg: u.Msg,
	}, u.Err
}

func (u *UserServiceServer) GetByIdProto(ctx context.Context, req *gen.GetByIdReq) (*gen.GetByIdResp, error) {
	log.Println("req:", req)
	return &gen.GetByIdResp{
		User: &gen.User{
			Id:   0,
			Name: u.Msg,
		},
	}, u.Err
}

func (u *UserServiceServer) Name() string {
	return "user-service"
}

type UserServiceTimeout struct {
	// 用反射来赋值
	// 类型是函数的字段，它不是方法（它不是定义在 UserService 上的方法）
	GetByID      func(ctx context.Context, req *GetByIDReq) (*GetByIDResp, error)
	GetByIdProto func(ctx context.Context, req *gen.GetByIdReq) (*gen.GetByIdResp, error)
}

func (u UserServiceTimeout) Name() string {
	return "user-service-timeout"
}

type UserServiceServerTimeout struct {
	t     *testing.T
	sleep time.Duration
	Err   error
	Msg   string
}

func (u *UserServiceServerTimeout) Name() string {
	return "user-service-timeout"
}

func (u *UserServiceServerTimeout) GetByID(ctx context.Context, req *GetByIDReq) (*GetByIDResp, error) {
	if _, ok := ctx.Deadline(); !ok {
		u.t.Fatal("context deadline required")
	}
	time.Sleep(u.sleep)
	return &GetByIDResp{
		Msg: u.Msg,
	}, u.Err
}
