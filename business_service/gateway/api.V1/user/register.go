package user

import (
	"github.com/kataras/iris"
	"github.com/golang/protobuf/proto"
	pb "github.com/chenshaobo/vent/business_service/proto"
	"github.com/chenshaobo/vent/business_service/rpclient"
	"github.com/chenshaobo/vent/business_service/utils"
	"golang.org/x/net/context"
	"github.com/chenshaobo/vent/business_service/gateway/api.V1/apiUtils"
)


//REGISTER FOR IS API FOR  /api/v1/user/register
func Register(c *iris.Context) {
	body := c.PostBody()
	c2s := &pb.RegisterC2S{}
	s2c := &pb.RegisterS2C{}

	err := proto.Unmarshal(body, c2s)
	if err != nil {
		s2c.ErrCode = utils.ErrParams
		utils.SetBody(c,s2c)
		return
	}
	conn := rpclient.Get(utils.RegisterSer)
	if conn == nil {
		s2c.ErrCode = utils.ErrServer
		utils.SetBody(c,s2c)
		return
	}
	rc := pb.NewRegisterClient(conn)
	s2c, err = rc.Register(context.Background(), c2s)
	utils.PrintErr(err)
	if err == nil {
		sessionKey := apiUtils.GetUserSessionKey(s2c.UserId)
		c.Session().Set(sessionKey,s2c.Session)
	}
	utils.SetBody(c,s2c)
}
