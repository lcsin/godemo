/**
rpc server
*/
package main

import (
	"net"
	"net/rpc"
)

func main() {
	// 1. 实例化一个server
	listen, _ := net.Listen("tcp", ":1234")
	// 2. 注册处理逻辑handler
	err := RegisterHelloService(&HelloService{})
	if err != nil {
		panic(err)
	}
	// 3.启动服务
	for {
		// 接受连接
		conn, _ := listen.Accept()
		// goroutine处理连接
		go rpc.ServeConn(conn)
	}
}
