package main

import (
	"context"
	"fmt"
	"godemo/gRPC/interceptor/pb"
	"google.golang.org/grpc"
	"time"
)

// 定义grpc客户端拦截器
var clientInterceptor grpc.UnaryClientInterceptor

// 自定义服务调用
var clientInvoker grpc.UnaryInvoker

func main() {
	// 处理业务
	clientInvoker = func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		fmt.Println("my client invoker...")
		return nil
	}

	// 实现客户端拦截器业务
	clientInterceptor = func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		start := time.Now()
		clientInvoker(ctx, method, req, req, cc, opts...)
		err := invoker(ctx, method, req, reply, cc, opts...)
		fmt.Printf("method=%s req=%v rep=%v duration=%s error=%v\n", method, req, reply, time.Since(start), err)
		return err
	}

	// grpc 客户端配置
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithUnaryInterceptor(clientInterceptor))

	// 1. 建立连接
	conn, err := grpc.Dial("localhost:8080", opts...)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := pb.NewGreeterClient(conn)

	// 2. 远程调用
	reply, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "lcx"})
	if err != nil {
		panic(err)
	}

	// 3. 打印调用结果
	fmt.Println("reply:", reply)
}
