package user

import (
	"github.com/kataras/iris"
	"github.com/golang/protobuf/proto"
	pb "github.com/chenshaobo/vent/business_service/proto"
	"github.com/chenshaobo/vent/business_service/rpclient"
	"github.com/chenshaobo/vent/business_service/utils"
	"golang.org/x/net/context"
	"github.com/chenshaobo/vent/business_service/gateway/api.V1/apiUtils"
	"strconv"
	"github.com/jbrodriguez/mlog"
)

func GetChangePasswordCaptcha(c *iris.Context){
	phoneNumber := c.Param("phoneNumber")
	c2s := &pb.PasswordCaptchaC2S{}
	s2c := &pb.PasswordCaptchaS2C{}
	phoneNumberUint,err := strconv.ParseUint(phoneNumber,10,64)
	if err  != nil{
		s2c.ErrCode = utils.ErrParams
		apiUtils.SetBody(c,s2c)
		return
	}

	c2s.PhoneNumber = phoneNumberUint
	conn := rpclient.Get(utils.UserSer)
	if conn == nil {
		s2c.ErrCode = utils.ErrServer
		apiUtils.SetBody(c,s2c)
		return
	}
	rc := pb.NewRegisterClient(conn)
	s2cTmp ,err  := rc.GetRegCaptcha(context.Background(),c2s)
	if err != nil {
		s2c.ErrCode = utils.ErrServer
		mlog.Error(err)
		apiUtils.SetBody(c,s2c)
		return
	}
	apiUtils.SetBody(c,s2cTmp)
}

func ChangePassword(c *iris.Context){
	body := c.PostBody()
	c2s := &pb.ChangePasswordC2S{}
	s2c := &pb.ChangePasswordS2C{}

	err := proto.Unmarshal(body, c2s)
	if err != nil {
		s2c.ErrCode = utils.ErrParams
		apiUtils.SetBody(c,s2c)
		mlog.Error(err)
		return
	}

	conn := rpclient.Get(utils.UserSer)
	if conn == nil {
		s2c.ErrCode = utils.ErrServer
		apiUtils.SetBody(c,s2c)
		mlog.Error(err)
		return
	}
	rc := pb.NewUserInfoManagerClient(conn)
	s2cTmp, err := rc.UserInfoModify(context.Background(), c2s)
	if err != nil{
		s2c.ErrCode = utils.ErrServer
		mlog.Error(err)
		apiUtils.SetBody(c,s2c)
	}
	s2c = s2cTmp
	apiUtils.SetBody(c,s2c)
}