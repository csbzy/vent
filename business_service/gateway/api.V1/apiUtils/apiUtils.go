package apiUtils

import (
	"github.com/kataras/iris"
	"strconv"
	"github.com/jbrodriguez/mlog"
	"github.com/chenshaobo/vent/business_service/utils"
)

func AuthSession(ctx *iris.Context){
	userIDStr := ctx.GetCookie(utils.CookieUserKey)
	if userIDStr ==""{
		return
	}
	userIDInt , err := strconv.ParseUint(userIDStr,10,64)
	if err != nil {
		return
	}
	reqSession := ctx.GetCookie(utils.CookieUserSession)
	if reqSession == ""{
		return
	}
	dbSession := ctx.Session().Get(GetUserSessionKey(userIDInt))
	mlog.Info("%v 's session %v  session from redis:%v", userIDInt,reqSession,dbSession)
	if reqSession != dbSession {
		return
	}
	ctx.Next()
}


func GetUserSessionKey(userId uint64) string {
	userIDStr := strconv.FormatUint(userId,10)
	sessionKey := utils.AccountSessionPrefix + userIDStr
	return sessionKey
}