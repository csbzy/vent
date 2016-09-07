package user

import (
	"github.com/kataras/iris"
	"github.com/chenshaobo/vent/business_service/rpclient"
	pb "github.com/chenshaobo/vent/business_service/proto"
	"github.com/chenshaobo/vent/business_service/utils"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
)

func Login(c *iris.Context){
		body := c.PostBody()
		c2s := &pb.LoginC2S{}
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
		rc := pb.NewLoginClient(conn)
		res, err := rc.Login(context.Background(), c2s)
		utils.PrintErr(err)
		buf, _ := proto.Marshal(res)
		c.Gzip(buf, 1)
	}

