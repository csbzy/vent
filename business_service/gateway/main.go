package main

import(
	"github.com/kataras/iris"
	"flag"
	"github.com/jbrodriguez/mlog"
	"github.com/chenshaobo/vent/business_service/rpclient"
	"github.com/chenshaobo/vent/business_service/gateway/api.V1/user"
	"github.com/chenshaobo/vent/business_service/gateway/api.V1/signal"
	"github.com/chenshaobo/vent/business_service/gateway/api.V1/geography"
	"github.com/chenshaobo/vent/business_service/utils"
	"github.com/chenshaobo/vent/business_service/consul"
	"gopkg.in/square/go-jose.v1/json"
	"github.com/kataras/go-sessions/sessiondb/redis"
	"github.com/kataras/go-sessions/sessiondb/redis/service"
	"strconv"
)

var (
	reg = flag.String("reg","172.16.7.119:8500","service register ")
)

func main(){
	flag.Parse()
	rpclient.Init(*reg)
	initSessionDB()
	initApi()
}

func initApi(){
	iris.UseFunc(log)
	user.SetupUserApi()
	signal.SetupSignalApi()
	geography.SetupGeoApi()
	///iris.UseFunc(fin)
	iris.AddServer(iris.ServerConfiguration{ListeningAddr: ":443", CertFile: "server.crt", KeyFile: "server.key"}) // you can close this server with .Close()
	iris.Listen("0.0.0.0:8080")
}

func initSessionDB(){
	serBytes, err := consul.Query(*reg, utils.SessionConfig)
	if err != nil{
		panic(err)
	}
	sessionConfig := &utils.RedisConfig{}
	json.Unmarshal(serBytes,&sessionConfig)
	mlog.Info("services: %v,serbyte:%v",sessionConfig,serBytes)
	db := redis.New(service.Config{Network: service.DefaultRedisNetwork,
		Addr:          sessionConfig.Host,
		Password:      "",
		Database:      strconv.FormatUint(sessionConfig.DB,10),
		MaxIdle:       0,
		MaxActive:     0,
		IdleTimeout:   service.DefaultRedisIdleTimeout,
		Prefix:        "",
		MaxAgeSeconds: service.DefaultRedisMaxAgeSeconds}) // optionally configure the bridge between your redis server

	iris.UseSessionDB(db)

}



func log(ctx *iris.Context){
	userIDInt , err := strconv.ParseUint(ctx.GetCookie("u"),10,64)
	if err != nil{
		return
	}
   //@todo session you wen ti
	mlog.Info("request:%v,   session:%v,%v  session from redis:%v   params: %v",string(ctx.Path()),ctx.GetCookie("u"),ctx.GetCookie("s"),ctx.Session().Get(user.GetUserSessionKey(userIDInt)),string(ctx.PostBody()))
	ctx.Next()
}

func fin(ctx *iris.Context){
	ctx.SetConnectionClose()
}

