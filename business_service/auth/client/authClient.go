package main

import (
	pb "github.com/chenshaobo/vent/business_service/proto"
	"flag"
	"github.com/chenshaobo/vent/business_service/consul"
	"google.golang.org/grpc"
	"github.com/chenshaobo/vent/business_service/utils"
	"time"
	"golang.org/x/net/context"
	"fmt"
)

var (
	serv = flag.String("service","registerService","service name")
	reg = flag.String("reg","172.16.7.119:8500","register address")
)
func main(){
	flag.Parse()
	r := consul.NewResolver(*serv)
	b := grpc.RoundRobin(r)
	conn,err := grpc.Dial(*reg,grpc.WithInsecure(),grpc.WithBalancer(b),grpc.WithTimeout(time.Second *10))
	utils.PanicErr(err)
	ticker:= time.NewTicker(20*time.Second)
	for t := range ticker.C{
		client :=pb.NewRegisterClient(conn)
		resp,err := client.Register(context.Background(),&pb.RegisterReq{PhoneNumber:1111,Password:"11111"})

		utils.PanicErr(err)

		fmt.Printf("interval %v  ;   %v\n",t,resp.ErrCode)
	}
}
