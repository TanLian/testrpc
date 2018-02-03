package testrpc

// 支持的协议：gob、json
const (
	EDCODE = "gob"
)

type Config int

func (conf *Config) GetEdcodeConf() string {
	return EDCODE
}
