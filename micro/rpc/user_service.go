package rpc

import (
	"context"
	"log"
)

type UserService struct {
	// 用反射来赋值
	// 类型是函数的字段，它不是方法（它不是定义在 UserService 上的方法）
	GetByID func(ctx context.Context, req *GetByIDReq) (*GetByIDResp, error)
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
}

func (u *UserServiceServer) GetByID(ctx context.Context, req *GetByIDReq) (*GetByIDResp, error) {
	log.Println("req:", req)
	return &GetByIDResp{
		Msg: "hello world",
	}, nil
}

func (u *UserServiceServer) Name() string {
	return "user-service"
}
