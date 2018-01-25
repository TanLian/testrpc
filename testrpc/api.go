package testrpc

import (
	"net"
	//"reflect"
	"sync"
)

func NewServer() *Server {
	return &Server{
		ServiceMap:  make(map[string]map[string]*Service),
		serviceLock: sync.Mutex{}}
}

func NewClient(conn net.Conn) *Client {
	return &Client{
		conn: conn}
}
