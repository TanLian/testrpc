package main

import (
	"log"
	"net"
	"os"
	"testrpc"
)

type Args struct {
	A, B int
}

func main() {

	// 连接本机的1234端口，返回一个net.Conn对象
	conn, err := net.Dial("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Println(err.Error())
		os.Exit(-1)
	}

	// main函数退出时关闭该网络连接
	defer conn.Close()

	// 创建一个rpc client对象
	client := testrpc.NewClient(conn)
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
