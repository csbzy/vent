package api_test
import (
	pb "github.com/chenshaobo/vent/business_service/proto"
	"testing"
	"github.com/golang/protobuf/proto"
	"net/http"
	"bytes"
	"fmt"
	"strconv"
)

var buf *bytes.Buffer
var userID uint64
var session string
func TestRegister(t *testing.T) {
	userID = 0
	session ="1"
	var c2s,s2c proto.Message
	buf = bytes.NewBuffer(make([]byte,0))
	c2s =&pb.RegisterC2S{PhoneNumber:13502700001,Password:"13502700001",Sex:1}
	s2c = &pb.RegisterS2C{}
	request("POST","http://127.0.0.1:8080/api/v1/user",c2s,s2c)


	c2s = &pb.LoginC2S{PhoneNumber:13502700001,Password:"13502700001"}
	s2c = &pb.LoginS2C{}
	request("PUT","http://127.0.0.1:8080/api/v1/user/session",c2s,s2c)
	userID = s2c.(*pb.LoginS2C).UserId
	session = s2c.(*pb.LoginS2C).Session
	fmt.Printf("cookie:%v,%v\n",userID,session)

	c2s = &pb.UserInfoModifyC2S{UserId:1,Nickname:"one",City:"shenzhen",Signature:"nice"}
	s2c = &pb.UserInfoModifyS2C{}
	request("PUT","http://127.0.0.1:8080/api/v1/user/info",c2s,s2c)

	s2c = &pb.UserInfoGetS2C{}
	request("GET","http://127.0.0.1:8080/api/v1/user/info/1" ,nil,s2c)
	fmt.Println("test is ok!")
}


func request(method string,url string ,m proto.Message,s2c proto.Message){
	var res *http.Response
	var err error
	fmt.Printf("cookie:%v,%v\n",userID,session)
	if method =="GET" {
		client := &http.Client{}
		req, _ := http.NewRequest(method, url,nil)
		req.Header.Set("Cookie","u=" + strconv.FormatUint(userID,10) + "; s=" + session)
		res, err = client.Do(req)
	}else {
		data, _ := proto.Marshal(m)
		client := &http.Client{}
		req, _ := http.NewRequest(method, url, bytes.NewReader(data))
		req.Header.Set("Cookie","u=" + strconv.FormatUint(userID,10) + "; s=" + session)
		res, err = client.Do(req)
	}
	if err != nil {
		fmt.Printf("res:%v\nerr:%v\n", res, err)
		return
	}
	_, err = buf.ReadFrom(res.Body)
	by := buf.Bytes()
	err = proto.Unmarshal(by, s2c)
	buf.Reset()
	fmt.Printf("request %v:%v\nparam:%v\nreturnBody:%v\nresponse:%v\nreturn:%v\n%v\n\n\n", method, url, m,by,res, s2c.String(), err)
}