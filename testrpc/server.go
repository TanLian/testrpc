package testrpc

import (
	"errors"
	"fmt"
	"log"
	"net"
	"reflect"
	"strconv"
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
		var edcode EdCode
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
		reqArgs := req.Args.(map[string]interface{})
		argv := reflect.New(service.ArgType)
		err = service.MakeArgType(reqArgs, argv)
		if err != nil {
			log.Println(err.Error())
			return
		}
		if argv.Kind() == reflect.Ptr {
			argv = argv.Elem()
		}

		// 构造reply
		reply := reflect.New(service.ReplyType.Elem())

		// 调用对应的函数
		function := service.Method.Func
		out := function.Call([]reflect.Value{reflect.New(server.ServerType.Elem()), argv, reply})
		if out[0].Interface() != nil {
			log.Println(out[0])
			return
		}

		// encode
		replyData, err := edcode.encode(reply.Elem().Interface())
		if err != nil {
			log.Println(err.Error())
			return
		}

		// 向conn写数据
		_, err = trans.WriteData(replyData)
		if err != nil {
			log.Println(err.Error())
			return
		}
	}
}

// 用data填充obj
func (service *Service) MakeArgType(data map[string]interface{}, obj reflect.Value) error {
	for k, v := range data {
		err := service.SetField(obj, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

//用map的值替换结构的值
func (service *Service) SetField(obj reflect.Value, name string, value interface{}) error {
	structValue := obj.Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)

	var err error
	if structFieldType != val.Type() {
		val, err = TypeConversion(fmt.Sprintf("%v", value), structFieldValue.Kind())
		if err != nil {
			return err
		}
	}

	structFieldValue.Set(val)
	return nil
}

// 将string类型的value值转换成reflect.Value类型
func TypeConversion(value string, ntype reflect.Kind) (reflect.Value, error) {
	switch ntype {
	case reflect.String:
		return reflect.ValueOf(value), nil
	case reflect.Int:
		i, err := strconv.Atoi(value)
		return reflect.ValueOf(i), err
	case reflect.Int8:
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int8(i)), err
	case reflect.Int16:
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int16(i)), err
	case reflect.Int32:
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int32(i)), err
	case reflect.Int64:
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(i), err
	case reflect.Float32:
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(float32(i)), err
	case reflect.Float64:
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(i), err
	default:
		return reflect.ValueOf(value), errors.New("unknown type：" + ntype.String())
	}
}
