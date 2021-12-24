package main

import (
	"godemo/golang/rpc/proxy/handler"
	"net/rpc"
)

type HelloServicer interface {
	Hello(request string, reply *string) error
}

type HelloService struct {
}

func (s *HelloService) Hello(request string, reply *string) error {
	*reply = "hello, " + request
	return nil
}

func RegisterHelloService(srv HelloServicer) error {
	return rpc.RegisterName(handler.HelloServiceName, srv)
}
