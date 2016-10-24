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
	reg = flag.String("reg", "172.16.7.119:8500", "register address")
	serviceName = flag.String("ser",utils.ServiceDefaultName,"get service 's config from consul")
	maxRedisConn = 0
)



func init() {
}

func main() {
	mlog.Start(mlog.LevelError,"")
	flag.Parse()
	serBytes, err := consul.Query(*reg, *serviceName)
	if err != nil{
		panic(err)
	}
	sers := make([]utils.ServerInfo,0)
	json.Unmarshal(serBytes,&sers)
	mlog.Info("services: %#v",sers)
	listenIP := getLocalIP()
	for _,ser :=  range sers {
		mlog.Info("start :%v",ser)
		switch ser.ServiceName{
		case utils.RegisterSer:
			registerService(ser,listenIP)
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
	if err != nil {
		panic(err)
	}
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", listenPort))
	rdc,err := redisapi.InitRedisClient(redisHost,int(redisDB),maxRedisConn,6,true)
	grpcSer := grpc.NewServer()
	captchaSer := &service.Service{Redisc:rdc}
	pb.RegisterCaptchaManagerServer(grpcSer,captchaSer)
	mlog.Info("start captcha service ok.")
	go grpcSer.Serve(lis)
}

func registerService(s utils.ServerInfo,listenIP string){
	redisHost := s.RedisConfig.Host
	redisDB := s.RedisConfig.DB
	listenPort := s.Port
	serviceName := s.ServiceName
	err := consul.Register(serviceName, listenIP, listenPort, *reg, time.Second * 30,  40)
	if err != nil {
		panic(err)
	}
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", listenPort))
	rdc,err := redisapi.InitRedisClient(redisHost,int(redisDB),maxRedisConn,6,true)
	grpcSer := grpc.NewServer()
	authSer := &service.Service{Redisc:rdc}
	pb.RegisterRegisterServer(grpcSer,authSer)
	pb.RegisterLoginServer(grpcSer,authSer)
	pb.RegisterUserInfoManagerServer(grpcSer,authSer)
	mlog.Info("start auth service ok.")
	go grpcSer.Serve(lis)

}

func relationService(s utils.ServerInfo,listenIP string){
	redisHost := s.RedisConfig.Host
	redisDB := s.RedisConfig.DB
	listenPort := s.Port
	serviceName := s.ServiceName
	err := consul.Register(serviceName, listenIP, listenPort, *reg, time.Second * 30,  40)
	if err != nil {
		panic(err)
	}
	rpclient.Init(*reg)
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", listenPort))
	rdc,err := redisapi.InitRedisClient(redisHost,int(redisDB),maxRedisConn,6,true)
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
	if err != nil {
		panic(err)
	}
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", listenPort))
	rdc,err := redisapi.InitRedisClient(redisHost,int(redisDB),maxRedisConn,6,true)
	grpcSer := grpc.NewServer()
	authSessionSer := &service.Service{Redisc:rdc}
	pb.RegisterSessionManagerServer(grpcSer,authSessionSer)

	mlog.Info("start relation service ok.")
	go grpcSer.Serve(lis)
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}