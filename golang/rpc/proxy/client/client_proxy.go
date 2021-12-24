package main

import (
	"godemo/golang/rpc/proxy/handler"
	"net/rpc"
)

type HelloServiceProxy struct {
	*rpc.Client
}

func NewHelloServiceProxy(proto, addr string) *HelloServiceProxy {
	client, err := rpc.Dial(proto, addr)
	if err != nil {
		panic(err)
	}
	return &HelloServiceProxy{
		client,
	}
}

func (s *HelloServiceProxy) Hello(request string, reply *string) error {
	return s.Call(handler.HelloServiceName+".Hello", request, reply)
}
