package utils

const (
	RegisterSer = "registerService"
	ReleationSer = "releationService"


	SessionConfig = "SessionConfig"
)


type  ServerInfo struct {
	ServiceName    string `json:"serviceName"`
	Port          int `json:"port"`
	RedisConfig   RedisConfig `json:"redisConfig"`
}

type RedisConfig struct {
	Host string `json:"host"`
	DB   uint64  `json:"dbNum"`
}

