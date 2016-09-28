package main

import (
	"flag"
	"net"
	"fmt"
	"time"
	"google.golang.org/grpc"
	"github.com/chenshaobo/redisapi"
)

import (
	"github.com/chenshaobo/vent/business_service/consul"
	 pb "github.com/chenshaobo/vent/business_service/proto"
	"github.com/chenshaobo/vent/business_service/rpcService/service"
	"github.com/jbrodriguez/mlog"
	"gopkg.in/square/go-jose.v1/json"
	"os"
	"os/signal"
	"syscall"
)

var (
	reg = flag.String("reg", "172.16.7.119:8500", "register address")
)

type  ServerInfo struct {
	ServiceName    string `json:"serviceName"`
	Port          int `json:"port"`
	RedisConfig   redisConfig `json:"redisConfig"`
}

type redisConfig struct {
	Host string `json:"host"`
	DB   int  `json:"dbNum"`
}

func init() {
}

func main() {
	flag.Parse()
	serBytes, err := consul.Query(*reg, "register01")
	if err != nil{
		panic(err)
	}
	sers := make([]ServerInfo,0)
	json.Unmarshal(serBytes,&sers)
	mlog.Info("services: %#v",sers)
	for _,ser :=  range sers {
		mlog.Info("start :%v",ser)
		switch ser.ServiceName{
		case "registerServer":
			authService(ser)
		default:

		}
	}


	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs,syscall.SIGKILL,syscall.SIGINT,syscall.SIGTERM)
	<-sigs
}


func authService(s ServerInfo){
	redisHost := s.RedisConfig.Host
	redisDB := s.RedisConfig.DB
	listenPort := s.Port
	serviceName := s.ServiceName

	mlog.Info("register consul")
	err := consul.Register(serviceName, "127.0.0.1", listenPort, *reg, time.Second * 30,  40)
	if err != nil {
		panic(err)
	}
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", listenPort))
	rdc,err := redisapi.InitRedisClient(redisHost,redisDB,6,6,true)
	grpcSer := grpc.NewServer()
	authSer := &service.Service{Redisc:rdc}
	pb.RegisterRegisterServer(grpcSer,authSer)
	pb.RegisterLoginServer(grpcSer,authSer)
	pb.RegisterUserInfoManagerServer(grpcSer,authSer)
	mlog.Info("start auth service ok.")
	go grpcSer.Serve(lis)

}