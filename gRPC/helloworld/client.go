package main

import (
	"context"
	"fmt"
	"godemo/gRPC/helloworld/pb"
	"google.golang.org/grpc"
)

func main() {
	// 1. 建立连接
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
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
