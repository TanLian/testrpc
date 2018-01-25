package testrpc

import (
	"encoding/json"
)

type EdCode int

func (edcode EdCode) encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (edcode EdCode) decode(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
