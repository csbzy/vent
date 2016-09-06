package user

import (
	"github.com/kataras/iris"
	"github.com/golang/protobuf/proto"
	pb "github.com/chenshaobo/vent/business_service/proto"
	"github.com/chenshaobo/vent/business_service/rpclient"
	"github.com/chenshaobo/vent/business_service/utils"
	"golang.org/x/net/context"
)


//REGISTER FOR IS API FOR  /api/v1/user/register
func Register(c *iris.Context) {
	body := c.PostBody()
	c2s := &pb.RegisterC2S{}
	conn := rpclient.Get("registerService")
	if conn == nil {
		c.EmitError(iris.StatusInternalServerError)
		return
	}
	err := proto.Unmarshal(body, c2s)

	if err != nil {
		c.EmitError(iris.StatusPreconditionFailed)
		return
	}
	rc := pb.NewRegisterClient(conn)
	res, err := rc.Register(context.Background(), c2s)
	utils.PrintErr(err)
	buf, _ := proto.Marshal(res)
	c.Gzip(buf, 1)
}
