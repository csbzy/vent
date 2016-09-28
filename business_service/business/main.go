package  main

import (
	"flag"
	"net"
	"fmt"
	"time"
	"google.golang.org/grpc"
	"github.com/chenshaobo/redisapi"
)

import(
	"github.com/chenshaobo/vent/business_service/consul"
	pb "github.com/chenshaobo/vent/business_service/proto"
	"github.com/chenshaobo/vent/business_service/rpcService/service"
	"github.com/jbrodriguez/mlog"
)

var (
	serv = flag.String("service", "releationService", "service name")
	port = flag.Int("port", 8200, "listening port")
	reg  = flag.String("reg", "127.0.0.1:8500", "register address")
)

func init(){
}

func main(){
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *port))
	if err != nil {
		panic(err)
	}

	rdc,err := redisapi.InitRedisClient("127.0.0.1:6379",0,6,6,true)

	err = consul.Register(*serv, "127.0.0.1", *port, *reg, time.Second * 30,  40)
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	geoSer := &service.Service{Redisc:rdc}
	pb.RegisterGeoManagerServer(s,geoSer)

	mlog.Info("start auth service ok.")
	s.Serve(lis)
}