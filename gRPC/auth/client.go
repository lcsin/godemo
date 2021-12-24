package main

import (
	"context"
	"fmt"
	"godemo/gRPC/auth/pb"
	"google.golang.org/grpc"
)

type ClientCredential struct {
}

func (c *ClientCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"appID":  "lcx",
		"appKey": "root",
	}, nil
}

func (c *ClientCredential) RequireTransportSecurity() bool {
	return false
}

func main() {
	// grpc提供的auth拦截器
	grpc.WithPerRPCCredentials(&ClientCredential{})

	// grpc 客户端配置
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithPerRPCCredentials(&ClientCredential{}))

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
