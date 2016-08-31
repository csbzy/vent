package  main

import (
	"github.com/uber-go/zap"
	"flag"
	"net"
	"fmt"
	"time"
	"google.golang.org/grpc"
)

import(
	"github.com/chenshaobo/vent/business_service/consul"
	pb "github.com/chenshaobo/vent/business_service/proto"
)

var (
	serv = flag.String("service", "registerService", "service name")
	port = flag.Int("port", 8100, "listening port")
	// reg  = flag.String("reg", "127.0.0.1:8500", "register address")
	reg  = flag.String("reg", "172.16.7.119:8500", "register address")
)
var logger zap.Logger
func init(){
	logger = zap.New(zap.NewJSONEncoder())
}
func main(){
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *port))
	if err != nil {
		panic(err)
	}

	err = consul.Register(*serv, "127.0.0.1", *port, *reg, time.Second * 30,  40)
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	pb.RegisterRegisterServer(s, &Registor{})
	s.Serve(lis)
}
