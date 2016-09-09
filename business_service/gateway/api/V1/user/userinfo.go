package user

import (
	"github.com/kataras/iris"
	pb "github.com/chenshaobo/vent/business_service/proto"
	"github.com/chenshaobo/vent/business_service/rpclient"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	"github.com/chenshaobo/vent/business_service/utils"
	"strconv"
	"github.com/jbrodriguez/mlog"
)

func InfoModify(c *iris.Context){
	body := c.PostBody()
	c2s := &pb.UserInfoModifyC2S{}
	s2c := &pb.UserInfoModifyS2C{}

	err := proto.Unmarshal(body, c2s)
	if err != nil {
		s2c.ErrCode = utils.ErrParams
		utils.SetBody(c,s2c)
		mlog.Info("err params")
		return
	}

	conn := rpclient.Get("registerService")
	if conn == nil {
		s2c.ErrCode = utils.ErrServer
		utils.SetBody(c,s2c)
		mlog.Info("err conn")
		return
	}
	rc := pb.NewUserInfoManagerClient(conn)
	s2c, err = rc.UserInfoModify(context.Background(), c2s)
	utils.PrintErr(err)
	utils.SetBody(c,s2c)
}

func InfoGet(c *iris.Context){
	userID ,err:= strconv.ParseUint(c.Param("userID"),10,64)
	s2c := &pb.UserInfoGetS2C{}
	if err !=nil {
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

	c2s := &pb.UserInfoGetC2S{UserId:userID}
	rc := pb.NewUserInfoManagerClient(conn)
	s2c, err = rc.UserInfoGet(context.Background(), c2s)
	utils.PrintErr(err)
	utils.SetBody(c,s2c)
}

