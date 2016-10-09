package api_test

import (
	pb "github.com/chenshaobo/vent/business_service/proto"
	"testing"
	"github.com/golang/protobuf/proto"
	"net/http"
	"bytes"
	"fmt"
	"strconv"
	"sync"
)

var wg sync.WaitGroup

func TestRegister(t *testing.T) {
	apiRun(13510162778)
	fmt.Println("test is ok!")
}

func runApi(workerSize int,b *testing.B) {
	var i int
	fmt.Printf("run times:%v\n",b.N)
	phoneNum := uint64(13510162778)
	totalJobSize :=  b.N
	jobs := make(chan uint64,totalJobSize)

	for i = 0;i <workerSize;i++{
		go worker(jobs)
	}
	for i = 0; i < (totalJobSize); i++ {
		jobs <- (phoneNum+uint64(i))
	}
	wg.Wait()
}



func BenchmarkApi1000(b  *testing.B) {
	runApi(1000,b)
}

func BenchmarkApi10000(b  *testing.B) {
	runApi(10000,b)
}






func worker(jobs chan uint64){
	for phoneNum := range jobs{
		apiRun(phoneNum)
	}
}
func apiRun(phonenum uint64) {
	wg.Add(1)
	var c2s, s2c proto.Message
	var userID uint64
	var session string
	c2s = &pb.RegisterC2S{PhoneNumber:phonenum, Password:"13502700001", Sex:1}
	s2c = &pb.RegisterS2C{}
	request("POST", "http://127.0.0.1:8080/api/v1/user", c2s, s2c, 0, "")

	c2s = &pb.LoginC2S{PhoneNumber:phonenum, Password:"13502700001"}
	s2c = &pb.LoginS2C{}
	request("PUT", "http://127.0.0.1:8080/api/v1/user/session", c2s, s2c, 0, "")
	userID = s2c.(*pb.LoginS2C).UserId
	session = s2c.(*pb.LoginS2C).Session
	if userID > 0 {
		c2s = &pb.UserInfoModifyC2S{UserId:userID, Nickname:"one", City:"shenzhen", Signature:"nice"}
		s2c = &pb.UserInfoModifyS2C{}
		request("PUT", "http://127.0.0.1:8080/api/v1/user/info", c2s, s2c, userID, session)

		s2c = &pb.UserInfoGetS2C{}
		request("GET", "http://127.0.0.1:8080/api/v1/user/info/" + strconv.FormatUint(userID, 10), nil, s2c, userID, session)

		c2s = &pb.GeoUploadC2S{UserId:userID, Latitude:22.5435866852, Longitude:113.9372047977}
		s2c = &pb.CommonS2C{}
		request("PUT", "http://127.0.0.1:8080/api/v1/coordinate", c2s, s2c, userID, session)
	}
	wg.Done()
}

func request(method string, url string, m proto.Message, s2c proto.Message, userID uint64, session string) {
	var res *http.Response
	var err error
	//fmt.Printf("cookie:%v,%v\n", userID, session)
	buf := bytes.NewBuffer(make([]byte, 0))
	if method == "GET" {
		client := &http.Client{}
		req, _ := http.NewRequest(method, url, nil)
		req.Header.Set("Cookie", "u=" + strconv.FormatUint(userID, 10) + "; s=" + session)
		res, err = client.Do(req)
	} else {
		data, _ := proto.Marshal(m)
		client := &http.Client{}
		req, _ := http.NewRequest(method, url, bytes.NewReader(data))
		req.Header.Set("Cookie", "u=" + strconv.FormatUint(userID, 10) + "; s=" + session)
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
	//fmt.Printf("request %v:%v\n" +
	//	"param:%v\n" +
	//	"returnBody:%v\n" +
	//	"response:%v\n" +
	//	"return:%v\n" +
	//	"err:%v\n\n\n", method, url, m, by, res, s2c.String(), err)
}