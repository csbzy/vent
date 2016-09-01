package utils

import (
	"github.com/uber-go/zap"
)

var Logger = zap.New(zap.NewJSONEncoder(
	zap.RFC3339Formatter("@timestamp"),
))
func PanicErr(err error){
	if err!=nil{
		panic(err)
	}
}


func PrintErr(err error){
	if err !=nil{
		Logger.Error(err.Error())
	}
}