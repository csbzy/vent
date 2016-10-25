package utils

const (
	UserSer = "userService"
	RelationSer = "relationService"
	AuthSer = "authService"
	CaptchaSer = "captchaService"
	SessionConfig = "SessionConfig"
	ServiceDefaultName = "service01"
)

const(
	CaptchaTypeReg = "captchaReg:"
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

