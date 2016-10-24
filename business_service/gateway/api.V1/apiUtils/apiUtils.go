package apiUtils

import (
	"strconv"
	"github.com/kataras/iris"

	"github.com/chenshaobo/vent/business_service/utils"
	"github.com/golang/protobuf/proto"
	"github.com/chenshaobo/vent/business_service/rpclient"
	pb "github.com/chenshaobo/vent/business_service/proto"
	"context"
	//"github.com/jbrodriguez/mlog"
)

func AuthSession(ctx *iris.Context){
	userIDStr := ctx.GetCookie(utils.CookieUserKey)
	if userIDStr ==""{
		ctx.Response.SetStatusCode(401)
		return
	}
	userIDInt , err := strconv.ParseUint(userIDStr,10,64)
	if err != nil {
		ctx.Response.SetStatusCode(401)
		return
	}
	reqSession := ctx.GetCookie(utils.CookieUserSession)
	if reqSession == ""{
		ctx.Response.SetStatusCode(401)
		return
	}
	authPass := Auth(userIDInt,reqSession)

	if !authPass{
		ctx.Response.SetStatusCode(401)
		return
	}
	ctx.Next()
}

func SetBody(c *iris.Context,pm proto.Message){
	buf, _ := proto.Marshal(pm)
	c.Gzip(buf, 1)
}

func Auth(userID uint64,session string) bool{
	authConn := rpclient.Get(utils.AuthSer)

	authSession := pb.NewSessionManagerClient(authConn)
	authReq := &pb.AuthC2S{UserID:userID,Session:session}
	authRes := &pb.CommonS2C{}
	authRes,err := authSession.AuthSession(context.Background(),authReq)
	if err !=nil || authRes.ErrCode > 0{

		return false
	}
	return true
}
