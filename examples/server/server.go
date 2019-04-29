package main

import (
	"github.com/tanlian/testrpc"
)

type Args struct {
	A, B int
}

type Arith int

func (t *Arith) Multiply(args Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func main() {
	// 创建一个rpc server对象
	newServer := testrpc.NewServer()

	// 向rpc server对象注册一个Arith对象，注册后，client就可以调用Arith的Multiply方法
	newServer.Register(new(Arith))

	// 监听本机的1234端口
	newServer.Server("tcp", "127.0.0.1:1234")
}
