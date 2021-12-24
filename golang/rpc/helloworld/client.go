/**
rpc client
*/
package main

import (
	"fmt"
	"net/rpc"
)

func main() {
	// 1. 建立连接
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		panic(err)
	}
	// 2. 远程调用
	reply := new(string)
	err = client.Call("HelloService.Hello", "lcx", reply)
	if err != nil {
		panic(err)
	}
	// 3. 调用结果
	fmt.Println("reply:", *reply)

}
