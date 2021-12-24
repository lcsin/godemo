package main

import (
	"context"
	"fmt"
	"godemo/gRPC/timeout/pb"
	"google.golang.org/grpc"
	"time"
)

func main() {
	// 1. 建立连接
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := pb.NewGreeterClient(conn)

	// 设置超时时间
	ctx, _ := context.WithTimeout(context.Background(), time.Second*3)

	// 2. 远程调用
	reply, err := client.SayHello(ctx, &pb.HelloRequest{Name: "lcx"})
	if err != nil {
		panic(err)
	}

	// 3. 打印调用结果
	fmt.Println("reply:", reply)
}
