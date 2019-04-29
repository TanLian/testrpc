package testrpc

import (
	"encoding/gob"
	"errors"
	"reflect"
)

type Service struct {
	Method    reflect.Method
	ArgType   reflect.Type
	ReplyType reflect.Type
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
