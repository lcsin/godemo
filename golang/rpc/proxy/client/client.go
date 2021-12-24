package main

import (
	"fmt"
)

func main() {
	// 1. 建立连接
	client := NewHelloServiceProxy("tcp", "localhost:1234")
	// 2. 远程调用
	reply := new(string)
	err := client.Hello("lcx", reply)
	if err != nil {
		panic(err)
	}
	// 3. 调用结果
	fmt.Println("reply:", *reply)
}
