/**
rpc server
*/
package main

import (
	"io"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type HelloService struct {
}

func (s *HelloService) Hello(request string, reply *string) error {
	*reply = "hello, " + request
	return nil
}

func main() {
	// 1. 注册RPC服务
	_ = rpc.RegisterName("HelloService", &HelloService{})
	// 2. 实例化Server
	http.HandleFunc("/jsonrpc", func(w http.ResponseWriter, r *http.Request) {
		var conn io.ReadWriteCloser = struct {
			io.Writer
			io.ReadCloser
		}{
			ReadCloser: r.Body,
			Writer:     w,
		}
		rpc.ServeRequest(jsonrpc.NewServerCodec(conn))
	})
	// 3. 启动服务
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
