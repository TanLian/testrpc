package testrpc

import (
	"log"
	"net"
)

type Client struct {
	conn net.Conn
}

func (client *Client) Close() {
	client.conn.Close()
}

func (client *Client) Call(methodName string, req interface{}, reply interface{}) error {

	// 构造一个Transfer
	request := NewRequest(methodName, req)

	// encode
	var edcode EdCode
	data, err := edcode.encode(request)
	if err != nil {
		return err
	}

	// write
	trans := NewTransfer(client.conn)
	_, err = trans.WriteData(data)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	// read
	data2, err := trans.ReadData()
	if err != nil {
		log.Println(err.Error())
		return err
	}

	// decode and assin to reply
	edcode.decode(data2, reply)

	// return
	return nil
}
