package main

import (
	"context"
	"godemo/gRPC/stream/pb"
	"google.golang.org/grpc"
	"log"
	"sync"
	"time"
)

func main() {
	// 建立连接
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		return
	}
	defer conn.Close()
	// 创建连接对象
	client := pb.NewStreamServiceClient(conn)

	// 调用服务端推送流
	streamReq := &pb.StreamRequest{Request: "lcx"}
	resp, _ := client.GetStream(context.Background(), streamReq)
	for {
		data, err := resp.Recv()
		if err != nil {
			log.Println(err)
			break
		}
		log.Println(data.Reply)
	}

	// 调用客户端流
	putRes, _ := client.PutStream(context.Background())
	i := 1
	for {
		i++
		putRes.Send(&pb.StreamRequest{Request: "leicx"})
		time.Sleep(time.Second)
		if i > 10 {
			break
		}
	}

	// 服务端 客户端 双向流
	wg := sync.WaitGroup{}
	wg.Add(2)
	allStr, _ := client.AllStream(context.Background())
	go func() {
		defer wg.Done()
		for {
			data, err := allStr.Recv()
			if err != nil {
				log.Println("err:", err.Error())
				break
			} else {
				log.Println("收到服务端的响应:" + data.Reply)
			}
		}
	}()

	go func() {
		defer wg.Done()
		for {
			err := allStr.Send(&pb.StreamRequest{Request: "客户端请求:ray"})
			if err != nil {
				log.Println("err:", err.Error())
				break
			}
			time.Sleep(time.Second)
		}
	}()
	wg.Wait()
}
