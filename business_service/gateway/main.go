package main

import(
	"github.com/kataras/iris"
	"flag"
	"github.com/jbrodriguez/mlog"
	"github.com/chenshaobo/vent/business_service/rpclient"
	"github.com/chenshaobo/vent/business_service/gateway/api.V1/user"
	"github.com/chenshaobo/vent/business_service/gateway/api.V1/signal"
	"github.com/chenshaobo/vent/business_service/gateway/api.V1/geography"
	"github.com/chenshaobo/vent/business_service/gateway/api.V1/relation"
	"github.com/chenshaobo/vent/business_service/utils"
)

var (
	DefaultIP = utils.GetLocalIP() + ":8500"
	reg = flag.String("reg",DefaultIP,"service register ")
)

func main(){
	flag.Parse()
	rpclient.Init(*reg)
	//the iris session manager not suit my busniess because of the session.sid
	//initSessionDB()
	initApi()
}

func initApi(){
	mlog.Start(mlog.LevelInfo,"")
	iris.UseFunc(log)
	user.SetupUserApi()
	signal.SetupSignalApi()
	geography.SetupGeoApi()
	relation.SetupRelationApi()

	iris.UseFunc(fin)
	///iris.UseFunc(fin)
	iris.AddServer(iris.ServerConfiguration{ListeningAddr: ":443", CertFile: "server.crt", KeyFile: "server.key"}) // you can close this server with .Close()
	iris.Listen("0.0.0.0:8080")

}

func log(ctx *iris.Context){
	mlog.Info("request:%v method:%v params: %v",
		string(ctx.Path()),
		string(ctx.Method()),
		string(ctx.PostBody()))
	ctx.Next()
}

func fin(ctx *iris.Context){
	defer func (){
		if err:= recover();err !=nil{
			ctx.SetStatusCode(500)
		}
	}()
}

