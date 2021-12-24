package main

import (
	"context"
	"fmt"
	"godemo/gRPC/interceptor/pb"
	"google.golang.org/grpc"
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

// 自定义handler，可以在这里面处理拦截的请求
var serverHandler grpc.UnaryHandler

func main() {
	// 处理业务
	serverHandler = func(ctx context.Context, req interface{}) (interface{}, error) {
		fmt.Println("my server handler...")
		return nil, nil
	}

	// 实现grpc拦截器业务
	serverInterceptor = func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		// 继续处理请求
		fmt.Println("接收到新请求")
		serverHandler(ctx, req)
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
