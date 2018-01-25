使用Golang实现的RPC，用法如下：
Server.go

* 定义服务

```
type Args struct {
	A, B int
}

type Arith int

func (t *Arith) Multiply(args Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}
```

* 创建一个rpc server对象

```
newServer := testrpc.NewServer()
```

* 注册服务

```
newServer.Register(new(Arith))
```

* 监听端口以及处理连接

```
conn, e := l.Accept()
if e != nil {
	log.Fatalf("l.Accept: %v", e)
}
go newServer.ServeConn(conn)
```

client.go

```
type Args struct {
	A, B int
}

func main() {

	// 连接本机的1234端口，返回一个net.Conn对象
	conn, err := net.Dial("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Println(err.Error())
		return
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
```

