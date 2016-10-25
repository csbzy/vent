package service

import (
	"golang.org/x/net/context"
	pb "github.com/chenshaobo/vent/business_service/proto"
	"github.com/chenshaobo/vent/business_service/utils"
	"strconv"
	"github.com/jbrodriguez/mlog"
	"github.com/chenshaobo/vent/business_service/rpclient"
)
func(s *Service) GetPasswordCaptcha(c context.Context,c2s *pb.PasswordCaptchaC2S) (*pb.PasswordCaptchaS2C, error){
	mlog.Info("request get password captcha:%v",c2s.PhoneNumber)

	s2c := &pb.PasswordCaptchaS2C{}
	phoneStr := strconv.FormatUint(c2s.PhoneNumber,10)
	if !s.Redisc.Exists(utils.AccountPhonePrefix + phoneStr) {
		s2c.ErrCode = utils.ErrAccountNotExits
		return s2c,nil
	}

	conn := rpclient.Get(utils.CaptchaSer)
	if conn == nil {
		s2c.ErrCode = utils.ErrServer
		return s2c,nil
	}
	rc := pb.NewCaptchaManagerClient(conn)

	getCaptchaS2C,err := rc.GetCaptcha(context.Background(),&pb.GetCaptchaC2S{Type:utils.CaptchaTypeReg,Key:phoneStr})
	if err != nil{
		getCaptchaS2C = &pb.GetCaptchaS2C{}
		s2c.ErrCode = utils.ErrServer
		return s2c,nil
	}
	s2c.Captcha = getCaptchaS2C.Captcha
	return s2c,nil

}
