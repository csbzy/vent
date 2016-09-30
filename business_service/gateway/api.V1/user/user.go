package user

import (
	"github.com/kataras/iris"
	"strconv"
	"github.com/chenshaobo/vent/business_service/utils"
)

func SetupUserApi(){
	userParty := iris.Party("/api/v1/user")
	userParty.Post("",Register)    //注册用户(创建)
	userParty.Put("/session",Login) //登录(更新session)
	userParty.Put("/info",InfoModify) //更新用户信息
	userParty.Get("/info/:userID",InfoGet) //获取用户信息
}


func GetUserSessionKey(userId uint64) string {
	userIDStr := strconv.FormatUint(userId,10)
	sessionKey := utils.AccountSessionPrefix + userIDStr
	return sessionKey
}