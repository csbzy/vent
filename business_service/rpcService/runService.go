package main

import (
	"flag"
	"net"
	"fmt"
	"time"
	"google.golang.org/grpc"
	"github.com/chenshaobo/redisapi"
	"github.com/chenshaobo/vent/business_service/consul"
	 pb "github.com/chenshaobo/vent/business_service/proto"
	"github.com/chenshaobo/vent/business_service/rpcService/service"
	"github.com/jbrodriguez/mlog"
	"gopkg.in/square/go-jose.v1/json"
	"os"
	"os/signal"
	"syscall"
	"github.com/chenshaobo/vent/business_service/utils"
	"github.com/chenshaobo/vent/business_service/rpclient"
)

var (
	DefaultIP = utils.GetLocalIP() + ":8500"
	reg = flag.String("reg", DefaultIP, "register address")
	serviceName = flag.String("ser",utils.ServiceDefaultName,"get service 's config from consul")
	maxRedisConn = 0
)



func init() {
}

func main() {
	mlog.Start(mlog.LevelInfo,"")
	flag.Parse()
	serBytes, err := consul.Query(*reg, *serviceName)
	if err != nil{
		panic(err)
	}
	sers := make([]utils.ServerInfo,0)
	json.Unmarshal(serBytes,&sers)
	mlog.Info("services: %#v",sers)
	listenIP := utils.GetLocalIP()
	for _,ser :=  range sers {
		mlog.Info("start :%v",ser)
		switch ser.ServiceName{
		case utils.UserSer:
			UserService(ser,listenIP)
		case utils.RelationSer:
			relationService(ser,listenIP)
		case utils.AuthSer:
			authSessionService(ser,listenIP)
		case utils.CaptchaSer:
			captchaService(ser,listenIP)
		default:
			mlog.Info("unknow service name :%v",ser.ServiceName)
		}
	}


	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs,syscall.SIGKILL,syscall.SIGINT,syscall.SIGTERM)
	<-sigs
}


func captchaService(s utils.ServerInfo,listenIP string){
	redisHost := s.RedisConfig.Host
	redisDB := s.RedisConfig.DB
	listenPort := s.Port
	serviceName := s.ServiceName
	err := consul.Register(serviceName, listenIP, listenPort, *reg, time.Second * 30,  40)
	utils.PanicErr(err)
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", listenPort))
	utils.PanicErr(err)
	rdc,err := redisapi.InitRedisClient(redisHost,int(redisDB),maxRedisConn,6,true)
	utils.PanicErr(err)
	grpcSer := grpc.NewServer()
	captchaSer := &service.Service{Redisc:rdc}
	pb.RegisterCaptchaManagerServer(grpcSer,captchaSer)
	mlog.Info("start captcha service ok.")
	go grpcSer.Serve(lis)
}

func UserService(s utils.ServerInfo,listenIP string){
	redisHost := s.RedisConfig.Host
	redisDB := s.RedisConfig.DB
	listenPort := s.Port
	serviceName := s.ServiceName
	err := consul.Register(serviceName, listenIP, listenPort, *reg, time.Second * 30,  40)
	utils.PanicErr(err)
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", listenPort))
	utils.PanicErr(err)
	rdc,err := redisapi.InitRedisClient(redisHost,int(redisDB),maxRedisConn,6,true)
	utils.PanicErr(err)
	grpcSer := grpc.NewServer()
	authSer := &service.Service{Redisc:rdc}
	pb.RegisterRegisterServer(grpcSer,authSer)
	pb.RegisterLoginServer(grpcSer,authSer)
	pb.RegisterUserInfoManagerServer(grpcSer,authSer)
	pb.RegisterPasswordManagerServer(grpcSer,authSer)
	mlog.Info("start auth service ok.")
	go grpcSer.Serve(lis)

}

func relationService(s utils.ServerInfo,listenIP string){
	redisHost := s.RedisConfig.Host
	redisDB := s.RedisConfig.DB
	listenPort := s.Port
	serviceName := s.ServiceName
	err := consul.Register(serviceName, listenIP, listenPort, *reg, time.Second * 30,  40)
	utils.PanicErr(err)
	rpclient.Init(*reg)
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", listenPort))
	utils.PanicErr(err)
	rdc,err := redisapi.InitRedisClient(redisHost,int(redisDB),maxRedisConn,6,true)
	utils.PanicErr(err)
	grpcSer := grpc.NewServer()
	relationSer := &service.Service{Redisc:rdc}

	pb.RegisterGeoManagerServer(grpcSer,relationSer)
	pb.RegisterRelationServer(grpcSer,relationSer)
	mlog.Info("start relation service ok.")
	go grpcSer.Serve(lis)
}

func authSessionService(s utils.ServerInfo,listenIP string){
	redisHost := s.RedisConfig.Host
	redisDB := s.RedisConfig.DB
	listenPort := s.Port
	serviceName := s.ServiceName
	err := consul.Register(serviceName, listenIP, listenPort, *reg, time.Second * 30,  40)
	utils.PanicErr(err)
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", listenPort))
	utils.PanicErr(err)
	rdc,err := redisapi.InitRedisClient(redisHost,int(redisDB),maxRedisConn,6,true)
	utils.PanicErr(err)
	grpcSer := grpc.NewServer()
	authSessionSer := &service.Service{Redisc:rdc}
	pb.RegisterSessionManagerServer(grpcSer,authSessionSer)

	mlog.Info("start relation service ok.")
	go grpcSer.Serve(lis)
}

