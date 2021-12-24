package main

import (
	"fmt"
	"godemo/gRPC/stream/pb"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
	"time"
)

type StreamServer struct {
	pb.UnimplementedStreamServiceServer
}

// GetStream 服务端流
func (s *StreamServer) GetStream(req *pb.StreamRequest, resp pb.StreamService_GetStreamServer) error {
	i := 0
	for {
		i++
		resp.Send(&pb.StreamReply{Reply: fmt.Sprintf("%v", time.Now().Unix())})
		time.Sleep(1 * time.Second)
		if i > 10 {
			break
		}
	}
	return nil
}

// PutStream 客户端流
func (s *StreamServer) PutStream(svc pb.StreamService_PutStreamServer) error {
	for {
		if tem, err := svc.Recv(); err == nil {
			log.Println(tem)
		} else {
			log.Println("break, err :", err)
			break
		}
	}

	return nil
}

// AllStream 双向流
func (s *StreamServer) AllStream(svc pb.StreamService_AllStreamServer) error {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for {
			if data, err := svc.Recv(); err == nil {
				log.Println("收到客户端请求:" + data.Request)
			} else {
				log.Println("request err:", err.Error())
				break
			}
		}
	}()

	go func() {
		defer wg.Done()
		for {
			err := svc.Send(&pb.StreamReply{Reply: "来自服务器:lcx"})
			if err != nil {
				log.Println("send err:", err.Error())
				break
			}
			time.Sleep(time.Second)
		}
	}()

	wg.Wait()
	return nil
}

func main() {
	srv := grpc.NewServer()

	pb.RegisterStreamServiceServer(srv, &StreamServer{})

	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	err = srv.Serve(listen)
	if err != nil {
		panic(err)
	}
}
