package user

import (
	"github.com/kataras/iris"
	"github.com/chenshaobo/vent/business_service/rpclient"
	pb "github.com/chenshaobo/vent/business_service/proto"
	"github.com/chenshaobo/vent/business_service/utils"
	"github.com/chenshaobo/vent/business_service/gateway/api.V1/apiUtils"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	"github.com/jbrodriguez/mlog"
)

func Login(c *iris.Context) {
	body := c.PostBody()
	c2s := &pb.LoginC2S{}
	s2c := &pb.LoginS2C{}

	err := proto.Unmarshal(body, c2s)
	if err != nil {
		s2c.ErrCode = utils.ErrParams
		apiUtils.SetBody(c,s2c)
		return
	}

	conn := rpclient.Get(utils.UserSer)
	if conn == nil {
		s2c.ErrCode = utils.ErrServer
		apiUtils.SetBody(c,s2c)
		return
	}

	rc := pb.NewLoginClient(conn)
	s2cTmp, err := rc.Login(context.Background(), c2s)

	if err != nil {
		s2c.ErrCode = utils.ErrServer
		apiUtils.SetBody(c,s2c)
		return
	}
	s2c =s2cTmp
	if s2c.ErrCode > 0 || s2c.UserID <= 0 {
		apiUtils.SetBody(c,s2c)
		return
	}


	session, errCode := getSession(s2c.UserID)
	if errCode > 0 {
		s2c.ErrCode = errCode
		apiUtils.SetBody(c,s2c)
		return
	}

	s2c.Session = session
	utils.PrintErr(err)
	apiUtils.SetBody(c,s2c)
}

func getSession(userID uint64) (string ,uint32){
	authConn := rpclient.Get(utils.AuthSer)
	getSession := pb.NewSessionManagerClient(authConn)
	getSessionC2S := &pb.GetSessionC2S{UserID:userID}
	getSessionS2C := &pb.GetSessionS2C{}
	getSessionS2C,err := getSession.GetSession(context.Background(),getSessionC2S)
	if err != nil {
		utils.PrintErr(err)
		return "",utils.ErrServer
	}
	if getSessionS2C.ErrCode > 0 {
		mlog.Info("get session error:%v",getSessionS2C.ErrCode)
		return "",getSessionS2C.ErrCode
	}
	return getSessionS2C.Session,0
}