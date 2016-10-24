package user

import (
	"github.com/kataras/iris"
	"github.com/chenshaobo/vent/business_service/gateway/api.V1/apiUtils"
)

func SetupUserApi(){
	userParty := iris.Party("/api/v1/user")
	userParty.Get("/check_code/register",GetRegisterCaptcha)
	userParty.Post("/check_code/register",Register)    //注册用户(创建)
	userParty.Put("/session",Login) //登录(更新session)
	userParty.Put("/info",apiUtils.AuthSession,InfoModify) //更新用户信息
	userParty.Get("/info/:userID",apiUtils.AuthSession,InfoGet) //获取用户信息
}

