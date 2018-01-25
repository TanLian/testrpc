package main

import (
	"log"
	"net"

	"testrpc"
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
	arith := new(Arith)
	newServer.Register(arith)

	// 监听本机的1234端口
	l, e := net.Listen("tcp", "127.0.0.1:1234")
	if e != nil {
		log.Fatalf("net.Listen tcp :0: %v", e)
	}

	for {
		// 阻塞直到从1234端口收到一个网络连接
		conn, e := l.Accept()
		if e != nil {
			log.Fatalf("l.Accept: %v", e)
		}

		//开始工作
		go newServer.ServeConn(conn)
	}
}
