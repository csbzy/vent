package utils

import (
	"github.com/jbrodriguez/mlog"
	"net"
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


func GetLocalIP() string {
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