package utils

import (
	"github.com/jbrodriguez/mlog"
	"github.com/kataras/iris"
	"github.com/golang/protobuf/proto"
)

func init(){
	mlog.Start(mlog.LevelInfo,"")
}

func PanicErr(err error){
	if err!=nil{
		panic(err)
	}
}


func PrintErr(err error){
	if err !=nil{
		mlog.Error(err)
	}
}




func SetBody(c *iris.Context,pm proto.Message){
	buf, _ := proto.Marshal(pm)
	c.Gzip(buf, 1)
}