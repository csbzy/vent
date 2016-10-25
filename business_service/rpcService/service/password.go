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

	getCaptchaS2C,err := rc.GetCaptcha(context.Background(),&pb.GetCaptchaC2S{Type:utils.CaptchaTypePassword,Key:phoneStr})
	if err != nil{
		getCaptchaS2C = &pb.GetCaptchaS2C{}
		s2c.ErrCode = utils.ErrServer
		return s2c,nil
	}
	s2c.Captcha = getCaptchaS2C.Captcha
	return s2c,nil

}

func(s *Service) ChangePassword(c context.Context, c2s *pb.ChangePasswordC2S) (*pb.ChangePasswordS2C, error)  {

	mlog.Info("Request change password:%v",c2s)
	s2c := &pb.ChangePasswordS2C{}
	phoneStr := strconv.FormatUint(c2s.PhoneNumber,10)
	v,err := s.Redisc.Get(utils.AccountPhonePrefix +phoneStr)
	if err !=nil {
		s2c.ErrCode = utils.ErrServer
		mlog.Error(err)
		return s2c,nil
	}
	userIDStr := string(v)

	//check captcha
	conn := rpclient.Get(utils.CaptchaSer)
	if conn == nil {
		s2c.ErrCode = utils.ErrServer
		return s2c,nil
	}
	rc := pb.NewCaptchaManagerClient(conn)
	checkCaptcha := &pb.CheckCaptchaC2S{Type:utils.CaptchaTypePassword,Key:strconv.FormatUint(c2s.PhoneNumber,10),Captcha:c2s.Captcha}
	checkCaptchac2s ,err := rc.CheckCaptcha(context.Background(), checkCaptcha)
	if err !=nil{
		s2c.ErrCode = utils.ErrServer
		mlog.Error(err)
		return s2c,nil
	}
	if checkCaptchac2s.ErrCode >0 {
		s2c.ErrCode = checkCaptchac2s.ErrCode
		mlog.Error(err)
		return s2c,nil
	}


	err = s.Redisc.Set(utils.AccountPasswordPrefix+userIDStr,[]byte(GetMD5Hash(c2s.Password)))
	if err !=nil {
		s2c.ErrCode = utils.ErrServer
		mlog.Error(err)
		return s2c,nil
	}

	return s2c,nil
}
