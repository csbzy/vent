package user

import (
	"github.com/kataras/iris"
	"github.com/chenshaobo/vent/business_service/rpclient"
	pb "github.com/chenshaobo/vent/business_service/proto"
	"github.com/chenshaobo/vent/business_service/utils"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
)

func Login(c *iris.Context) {
	body := c.PostBody()
	c2s := &pb.LoginC2S{}
	s2c := &pb.LoginS2C{}

	err := proto.Unmarshal(body, c2s)
	if err != nil {
		s2c.ErrCode = utils.ErrParams
		utils.SetBody(c,s2c)
		return
	}

	conn := rpclient.Get("registerService")
	if conn == nil {
		s2c.ErrCode = utils.ErrServer
		utils.SetBody(c,s2c)
		return
	}

	rc := pb.NewLoginClient(conn)
	s2c, err = rc.Login(context.Background(), c2s)
	utils.PrintErr(err)
	utils.SetBody(c,s2c)
}

