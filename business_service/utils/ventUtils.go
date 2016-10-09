package utils

import (
	"github.com/jbrodriguez/mlog"
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


