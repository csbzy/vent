package utils

const (
	ErrParams = 10000 + iota //请求参数错误
	ErrServer  // 服务器未知错误
	ErrAccountExits //账号已存在
	ErrAccountNotExits //账号不存在
	ErrPasswordWrong  //密码错误
	ErrSessionNotMatch //session错误
)
