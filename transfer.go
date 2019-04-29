package testrpc

import (
	"net"
)

const (
	EachReadBytes = 500
)

type Transfer struct {
	conn net.Conn
}

func NewTransfer(conn net.Conn) *Transfer {
	return &Transfer{conn: conn}
}

func (trans *Transfer) ReadData() ([]byte, error) {
	finalData := make([]byte, 0)
	for {
		data := make([]byte, EachReadBytes)
		i, err := trans.conn.Read(data)
		if err != nil {
			return nil, err
		}
		finalData = append(finalData, data[:i]...)
		if i < EachReadBytes {
			break
		}
	}
	return finalData, nil
}

func (trans *Transfer) WriteData(data []byte) (int, error) {
	num, err := trans.conn.Write(data)
	return num, err
}
