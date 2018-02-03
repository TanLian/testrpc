### 说明：
使用Golang实现的RPC，目前支持的序列化协议有GOB和JSON，默认采用GOB，具体可在conf.go里面配置，未来考虑支持更多的序列化，如xml、protobuf。
### 用法：
**Server.go**

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

**client.go**

client.go有两种使用方法，您可以先使用net.Dial返回net.Conn对象，再通过net.Conn对象生成rpc client对象，通过rpc client即可调用服务。也可以直接通过testrpc.Dial生成rpc client对象，再通过rpc client对象调用服务。

**用法一**：

* 连接本机的1234端口，返回一个net.Conn对象

```
conn, err := net.Dial("tcp", "127.0.0.1:1234")
if err != nil {
	log.Println(err.Error())
	return
}
```

* 创建一个rpc client对象

```
client := testrpc.NewClient(conn)
```

* 通过client调用服务

```
err = client.Call("Arith.Multiply", args, &reply)
if err != nil {
	log.Fatal("arith error:", err)
}
```
**用法二**

* 连接本机的1234端口，返回一个Client对象

```
client, err := testrpc.Dial("tcp", "127.0.0.1:1234")
if err != nil {
	log.Println(err.Error())
	return
}
```
* 通过client调用服务

```
err = client.Call("Arith.Multiply", args, &reply)
if err != nil {
	log.Fatal("arith error:", err)
}
```

文档地址：https://juejin.im/post/5a69e308518825733b0f151a