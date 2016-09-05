package user

import (
	"github.com/kataras/iris"
	"github.com/golang/protobuf/proto"
	pb "github.com/chenshaobo/vent/business_service/proto"
	"github.com/chenshaobo/vent/business_service/rclient"
	"github.com/chenshaobo/vent/business_service/utils"
	"golang.org/x/net/context"
	"github.com/uber-go/zap"
	"io/ioutil"
	"os"
	"fmt"
)


//REGISTER FOR IS API FOR  /api/v1/user/register
func Register(c *iris.Context) {
	body := c.PostBody()
	r := &pb.RegisterC2S{}
	fmt.Printf("body:%v",body)
	conn := rpclient.Get("registerService")
	if conn == nil {
		utils.Logger.Error("get register service error")
		c.EmitError(iris.StatusInternalServerError)
		return
	}
	err := proto.Unmarshal(body, r)
	ioutil.WriteFile("1.log",body,os.ModeAppend)

	if err != nil {
		utils.Logger.Error("Unmarshal",zap.String("unmarshal error",err.Error()),zap.Int("len",len(body)))
		c.EmitError(iris.StatusPreconditionFailed)
		return
	}
	rc := pb.NewRegisterClient(conn)
	res, err := rc.Register(context.Background(), &pb.RegisterReq{PhoneNumber:r.PhoneNumber, Password:r.Password})
	utils.Logger.Info("res", zap.Int("int", int(res.ErrCode)))
	buf, _ := proto.Marshal(res)
	c.Gzip(buf, 1)
}
