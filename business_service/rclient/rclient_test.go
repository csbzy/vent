package rclient_test

import (
	. "github.com/chenshaobo/vent/business_service/rclient"
	"testing"
)

func TestGet(t *testing.T) {
	Init("172.16.7.119:800")
	Register("registerServices")
	if Get("registerServices") == nil{
		t.Log("GET SERVICE error")
	}else{
		t.Log("get service ok")
	}
}