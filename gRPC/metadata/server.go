package main

import (
	"context"
	"fmt"
	"godemo/gRPC/interceptor/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"net"
)

type HelloServer struct {
	pb.UnimplementedGreeterServer
}

func (s *HelloServer) SayHello(ctx context.Context, request *pb.HelloRequest) (*pb.HelloReply, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		// 遍历metadata
		for k, v := range md {
			fmt.Println(k, v)
		}
		// 获取指定key的数据
		fmt.Println("username:", md["name"])
		fmt.Println("password:", md["password"])
	}
	return &pb.HelloReply{Message: "hello, " + request.Name}, nil
}

func main() {
	// 1. 创建一个服务
	srv := grpc.NewServer()

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
