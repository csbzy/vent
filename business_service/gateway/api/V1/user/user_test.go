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
	buf = bytes.NewBuffer(make([]byte,256))
	d :=&pb.RegisterC2S{PhoneNumber:13502700001,Password:"13502700001"}

	request("PUT","http://127.0.0.1:8080/api/v1/user/register",d,&pb.RegisterS2C{})

	l := &pb.LoginC2S{PhoneNumber:13502700001,Password:"13502700001"}
	request("POST","http://127.0.0.1:8080/api/v1/user/login",l,&pb.LoginS2C{})
}


func request(method string,url string ,m proto.Message,s2c proto.Message){
	data,err := proto.Marshal(m)
	mlog.Error(err)
	client := &http.Client{}
	req,err  := http.NewRequest(method,url,bytes.NewReader(data))
	res,err := client.Do(req)
	_,err = buf.ReadFrom(res.Body)
	err =proto.Unmarshal(buf.Bytes(),s2c)
	buf.Reset()
	mlog.Info("\nrequest:%v,body:%v\nreturn:%v,%v",url,m,s2c.String(),err)
}