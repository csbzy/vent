package user_test
import (
	pb "github.com/chenshaobo/vent/business_service/proto"
	"testing"

	"github.com/golang/protobuf/proto"
	"net/http"
	"bytes"
)

func TestRegister(t *testing.T) {
	d :=&pb.RegisterC2S{PhoneNumber:111111,Password:"111111"}
	data,err := proto.Marshal(d)
	if err != nil {
		t.Fatal("marshal error",err)
	}
	//
	//request := gorequest.New()
	//_,body,errs := request.Put("http://127.0.0.1:8080/api/v1/user/register").
	//	Set("Content-Type","application/x-protobuf").
	//	Send(data).EndBytes()

	client := &http.Client{}
	req,err  := http.NewRequest("PUT","http://127.0.0.1:8080/api/v1/user/register",bytes.NewReader(data))
	res,err := client.Do(req)
	if err !=nil {
		t.Errorf("request error ,return body:%v,err:%v",res.Body,err)
	}
	t.Logf("request return :%v,%v",res.Body,err)
}