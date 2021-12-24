package main

import (
	"context"
	"fmt"
	"godemo/gRPC/auth/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net"
)

type HelloServer struct {
	pb.UnimplementedGreeterServer
}

func (s *HelloServer) SayHello(ctx context.Context, request *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "hello, " + request.Name}, nil
}

// 定义grpc拦截器
var serverInterceptor grpc.UnaryServerInterceptor

func authHandler(ctx context.Context, req interface{}) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "无Token认证信息")
	}
	var (
		appID  string
		appKey string
	)

	if val, ok := md["appid"]; ok {
		appID = val[0]
	}
	if val, ok := md["appkey"]; ok {
		appKey = val[0]
	}
	if appID != "lcx" && appKey != "root" {
		return nil, status.Errorf(codes.Unauthenticated, "无效的Token信息")
	}
	auth := map[string]string{
		"appID":  appID,
		"appKey": appKey,
	}
	return auth, nil
}

func main() {
	// 实现grpc拦截器业务
	serverInterceptor = func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		fmt.Println("接收到新请求")
		// 处理token认证
		token, err := authHandler(ctx, req)
		if err != nil {
			return nil, err
		}
		fmt.Println("token:", token)
		// 继续处理请求
		res, err := handler(ctx, req)
		fmt.Println("请求处理完成")
		return res, err
	}
	// 添加grpc配置项
	var opts []grpc.ServerOption
	opts = append(opts, grpc.UnaryInterceptor(serverInterceptor))

	// 1. 创建一个服务
	srv := grpc.NewServer(opts...)

	// 2. 注册服务
	pb.RegisterGreeterServer(srv, &HelloServer{})

	// 3. 启动服务
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	err = srv.Serve(listen)
	if err != nil {
		panic(err)
	}
}
