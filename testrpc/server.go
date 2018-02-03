package testrpc

import (
	"encoding/gob"
	"errors"
	//"log"
	"net"
	"reflect"
	"strings"
	"sync"
)

type Service struct {
	Method    reflect.Method
	ArgType   reflect.Type
	ReplyType reflect.Type
}

type Server struct {
	ServiceMap  map[string]map[string]*Service
	serviceLock sync.Mutex
	ServerType  reflect.Type
}

func (server *Server) Register(obj interface{}) error {
	server.serviceLock.Lock()
	defer server.serviceLock.Unlock()

	//通过obj得到其各个方法，存储在servicesMap中
	tp := reflect.TypeOf(obj)
	val := reflect.ValueOf(obj)
	serviceName := reflect.Indirect(val).Type().Name()
	if _, ok := server.ServiceMap[serviceName]; ok {
		return errors.New(serviceName + " already registed.")
	}

	s := make(map[string]*Service)
	numMethod := tp.NumMethod()
	for m := 0; m < numMethod; m++ {
		service := new(Service)
		method := tp.Method(m)
		mtype := method.Type
		mname := method.Name

		service.ArgType = mtype.In(1)
		service.ReplyType = mtype.In(2)
		service.Method = method
		s[mname] = service

		err := service.RegisterGobArgsType()
		if err != nil {
			return err
		}
	}
	server.ServiceMap[serviceName] = s
	server.ServerType = reflect.TypeOf(obj)
	return nil
}

func (server *Server) ServeConn(conn net.Conn) {
	trans := NewTransfer(conn)
	for {
		// 从conn读数据
		data, err := trans.ReadData()
		if err != nil {
			return
		}

		// decode
		var req Request
		edcode, err := GetEdcode()
		if err != nil {
			return
		}
		err = edcode.decode(data, &req)
		if err != nil {
			return
		}

		// 根据MethodName拿到service
		methodStr := strings.Split(req.MethodName, ".")
		if len(methodStr) != 2 {
			return
		}
		service := server.ServiceMap[methodStr[0]][methodStr[1]]

		// 构造argv
		argv, err := req.MakeArgs(edcode, *service)

		// 构造reply
		reply := reflect.New(service.ReplyType.Elem())

		// 调用对应的函数
		function := service.Method.Func
		out := function.Call([]reflect.Value{reflect.New(server.ServerType.Elem()), argv, reply})
		if out[0].Interface() != nil {
			return
		}

		// encode
		replyData, err := edcode.encode(reply.Elem().Interface())
		if err != nil {
			return
		}

		// 向conn写数据
		_, err = trans.WriteData(replyData)
		if err != nil {
			return
		}
	}
}

// 如果是GOB编解码，则注册Args的类型，防止gob编解码错误
func (service *Service) RegisterGobArgsType() error {
	edcodeStr := new(Config).GetEdcodeConf()
	switch edcodeStr {
	case "gob":
		args := reflect.New(service.ArgType)
		if args.Kind() == reflect.Ptr {
			args = args.Elem()
		}
		gob.Register(args.Interface())
		return nil
	case "json":
		return nil
	default:
		return errors.New("Unknown edcode protocol: " + edcodeStr)
	}
}
