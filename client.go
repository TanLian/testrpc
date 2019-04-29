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

	// 构造一个Request
	request := NewRequest(methodName, req)

	// 如果是GOB编码，则要注册相应类型，防止gob编解码错误
	err := request.RegisterGobArgsType()
	if err != nil {
		log.Println(err.Error())
		return err
	}

	// encode
	edcode, err := GetEdcode()
	if err != nil {
		return err
	}
	data, err := edcode.encode(request)
	if err != nil {
		log.Println(err.Error())
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
