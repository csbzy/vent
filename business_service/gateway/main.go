package main

import(
	"github.com/chenshaobo/vent/business_service/gateway/api"
)
import "github.com/kataras/iris"

func main(){
	initApi()
}

func initApi(){
	iris.API("/api/users/register",api.RegisterApi{})
	//iris.
	iris.Listen("0.0.0.0:8080")
}
