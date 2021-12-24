package main

import (
	"context"
	"fmt"
	"godemo/gRPC/auth/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func main() {
	// 1. 建立连接
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := pb.NewGreeterClient(conn)

	// auth metadata
	md := metadata.New(map[string]string{
		"appID":  "lcx",
		"appKey": "root",
	})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	// 2. 远程调用
	reply, err := client.SayHello(ctx, &pb.HelloRequest{Name: "lcx"})
	if err != nil {
		panic(err)
	}

	// 3. 打印调用结果
	fmt.Println("reply:", reply)
}
