package user

import (
	"github.com/kataras/iris"
	"github.com/golang/protobuf/proto"
	pb "github.com/chenshaobo/vent/business_service/proto"
	"github.com/chenshaobo/vent/business_service/rclient"
	"github.com/chenshaobo/vent/business_service/utils"
	"golang.org/x/net/context"
	"github.com/uber-go/zap"
)


func Register(c *iris.Context) {
	body := c.PostBody()
	pbm := &pb.RegisterC2S{}
	conn := rclient.Get("registerService")
	if conn != nil {
		c.EmitError(iris.StatusInternalServerError)
		return
	}
	err := proto.Unmarshal(body, pbm)
	if err != nil {
		utils.Logger.Info("err")
		c.EmitError(iris.StatusInternalServerError)
		return
	}
	rc := pb.NewRegisterClient(conn)
	res, err := rc.Register(context.Background(), &pb.RegisterReq{PhoneNumber:pbm.PhoneNumber, Password:pbm.Password})
	utils.Logger.Info("res", zap.Int("int", int(res.ErrCode)))
	buf, _ := proto.Marshal(res)
	c.Gzip(buf, 1)
}
