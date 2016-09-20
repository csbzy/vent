package user_test
import (
	pb "github.com/chenshaobo/vent/business_service/proto"
	"testing"

	"github.com/golang/protobuf/proto"
	"net/http"
	"bytes"
	"github.com/jbrodriguez/mlog"
)

var buf *bytes.Buffer

func TestRegister(t *testing.T) {
	mlog.Start(mlog.LevelInfo,"")
	buf = bytes.NewBuffer(make([]byte,0))
	mlog.Info("Start test")
	d :=&pb.RegisterC2S{PhoneNumber:13502700001,Password:"13502700001",Sex:1}
	request("POST","http://127.0.0.1:8080/api/v1/user/register",d,&pb.RegisterS2C{})

	l := &pb.LoginC2S{PhoneNumber:13502700001,Password:"13502700001"}
	request("PUT","http://127.0.0.1:8080/api/v1/user/login",l,&pb.LoginS2C{})

	userInfoMod := &pb.UserInfoModifyC2S{UserId:1,Nickname:"one",City:"shenzhen",Signature:"nice"}
	request("PUT","http://127.0.0.1:8080/api/v1/user/info",userInfoMod,&pb.UserInfoModifyS2C{})

	request("GET","http://127.0.0.1:8080/api/v1/user/info/1" ,nil,&pb.UserInfoGetS2C{})
	mlog.Info("test is ok!")
}


func request(method string,url string ,m proto.Message,s2c proto.Message){
	var res *http.Response
	var err error
	if method =="GET" {
		client := &http.Client{}
		res, err = client.Get(url)
	}else {
		data, _ := proto.Marshal(m)
		client := &http.Client{}
		req, _ := http.NewRequest(method, url, bytes.NewReader(data))
		res, err = client.Do(req)
	}
	if err != nil {
		mlog.Info("res:%v\nerr:%v", res, err)
		return
	}
	_, err = buf.ReadFrom(res.Body)
	by := buf.Bytes()
	err = proto.Unmarshal(by, s2c)
	buf.Reset()
	mlog.Info("request %v:%v\nparam:%v\nreturnBody:%v\nresponse:%v\nreturn:%v\n%v\n\n\n", method, url, m,by,res, s2c.String(), err)
}