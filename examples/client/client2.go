package main

import (
	"github.com/tanlian/testrpc"
	"log"
)

type Args struct {
	A, B int
}

func main() {

	// 连接本机的1234端口，返回一个Client对象
	client, err := testrpc.Dial("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Println(err.Error())
		return
	}

	// main函数退出时关闭该client
	defer client.Close()

	// 调用远端Arith.Multiply函数
	args := Args{7, 8}
	var reply int
	err = client.Call("Arith.Multiply", args, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	log.Println(reply)

}
