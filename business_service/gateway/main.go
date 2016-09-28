package main

import(
	"github.com/chenshaobo/vent/business_service/gateway/api/V1/user"
	"github.com/kataras/iris"
	"flag"
	"github.com/chenshaobo/vent/business_service/rpclient"
	"github.com/jbrodriguez/mlog"
	"github.com/chenshaobo/vent/business_service/gateway/api/V1/signal"
	"github.com/chenshaobo/vent/business_service/gateway/api/V1/geography"
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
	//iris.
	//user api

	user.SetupUserApi()
	signal.SetupSignalApi()
	geography.SetupGeoApi()
	///iris.UseFunc(fin)
	iris.AddServer(iris.ServerConfiguration{ListeningAddr: ":443", CertFile: "server.crt", KeyFile: "server.key"}) // you can close this server with .Close()
	iris.Listen("0.0.0.0:8080")
}



func log(ctx *iris.Context){
	mlog.Info("request:%v,params: %v",string(ctx.Path()),string(ctx.PostBody()))
	ctx.Next()
}

func fin(ctx *iris.Context){
	ctx.SetConnectionClose()
}

