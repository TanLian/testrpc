package testrpc

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
)

type EdCode interface {
	encode(v interface{}) ([]byte, error)
	decode(data []byte, v interface{}) error
}

type GobEdCode int
type JsonEdCode int

func (edcode GobEdCode) encode(v interface{}) ([]byte, error) {
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	err := enc.Encode(v)
	if err != nil {
		return nil, err
	}
	return network.Bytes(), nil
}

func (edcode GobEdCode) decode(data []byte, v interface{}) error {
	buf := bytes.NewBuffer(data)
	return gob.NewDecoder(buf).Decode(v)
}

func (edcode JsonEdCode) encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (edcode JsonEdCode) decode(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func GetEdcode() (EdCode, error) {
	edcodeStr := new(Config).GetEdcodeConf()
	switch edcodeStr {
	case "gob":
		return *new(GobEdCode), nil
	case "json":
		return *new(JsonEdCode), nil
	default:
		return nil, errors.New("Unknown edcode protocol: " + edcodeStr)
	}
}
