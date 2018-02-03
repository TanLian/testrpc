package testrpc

import (
	"errors"
	"net"
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

func Dial(network, address string) (*Client, error) {
	if network != "tcp" {
		return nil, errors.New("Unsupported protocol")
	}

	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}

	return NewClient(conn), nil
}
