package main

import(
	"github.com/chenshaobo/vent/business_service/gateway/api/V1/user"
	"github.com/chenshaobo/vent/business_service/rclient"
	"github.com/kataras/iris"
	"flag"
	"github.com/chenshaobo/vent/business_service/utils"
	"github.com/uber-go/zap"
	pb "github.com/chenshaobo/vent/business_service/proto"
	"github.com/golang/protobuf/proto"
	"fmt"
)

var (
	reg = flag.String("reg","172.16.7.119:8500","service register ")
)

func main(){
	flag.Parse()
	TestUnmarshal()
	initApi()
}

func initApi(){
	utils.Logger.Info("init api")
	iris.UseFunc(log)
	userParty := iris.Party("/api/v1/user")
	userParty.Put("/register",user.Register)
	//iris.
	iris.Listen("0.0.0.0:8080")
}


func log(ctx *iris.Context){
	utils.Logger.Info("req",zap.String("url",string(ctx.Path())),)
	ctx.Next()
}

func fin(ctx *iris.Context){
	if err := recover();err !=nil{
		ctx.EmitError(iris.StatusInternalServerError)
	}
	ctx.Response.StatusCode()
	ctx.Next()
}


func TestUnmarshal(){
	d :=&pb.RegisterC2S{PhoneNumber:111111,Password:"111111"}
	data,err := proto.Marshal(d)
	fmt.Printf("DATA:%v,err:%v",data,err)
	r := &pb.RegisterC2S{}

	err = proto.Unmarshal(data,r)
	fmt.Printf("DATA:%v,err:%v",r,err)
}