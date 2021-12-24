/**
rpc server
*/
package main

import (
	"net"
	"net/rpc"
)

type HelloService struct {
}

func (s *HelloService) Hello(request string, reply *string) error {
	*reply = "hello, " + request
	return nil
}

func main() {
	// 1. 实例化一个server
	listen, _ := net.Listen("tcp", ":1234")
	// 2. 注册处理逻辑handler
	_ = rpc.RegisterName("HelloService", &HelloService{})
	// 3.启动服务
	conn, _ := listen.Accept() // 接受一个连接
	rpc.ServeConn(conn)
}
