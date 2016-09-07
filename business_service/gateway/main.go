package main

import(
	"github.com/chenshaobo/vent/business_service/gateway/api/V1/user"
	"github.com/kataras/iris"
	"flag"
	"github.com/chenshaobo/vent/business_service/rpclient"
	"github.com/jbrodriguez/mlog"
)

var (
	reg = flag.String("reg","172.16.7.119:8500","service register ")
)

func main(){
	flag.Parse()
	rpclient.Init(*reg)
	initApi()
}

func initApi(){
	iris.UseFunc(log)
	userParty := iris.Party("/api/v1/user")
	userParty.Put("/register",user.Register)
	userParty.Post("/login",user.Login)
	//iris.
	//iris.UseFunc(fin)
	iris.Listen("0.0.0.0:8080")
}


func log(ctx *iris.Context){
	mlog.Info("request:%v,params:%v",string(ctx.Path()),string(ctx.PostBody()))
	ctx.Next()
}

func fin(ctx *iris.Context){
	ctx.SetConnectionClose()
}

