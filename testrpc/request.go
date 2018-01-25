package testrpc

type Request struct {
	MethodName string
	Args       interface{}
}

func NewRequest(methodName string, args interface{}) *Request {
	return &Request{
		MethodName: methodName,
		Args:       args}
}
