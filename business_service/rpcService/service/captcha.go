package service

import (
	"golang.org/x/net/context"
	"github.com/jbrodriguez/mlog"
	"strconv"
	pb "github.com/chenshaobo/vent/business_service/proto"
	"github.com/chenshaobo/vent/business_service/utils"
	"encoding/binary"
	"math/rand"

)

func (s *Service ) GetCaptcha(ctx context.Context,req *pb.GetCaptchaC2S)(*pb.GetCaptchaS2C,error){
	mlog.Info("request get captcha with type:%v,key:%v" ,req.Type,req.Key)
	s2c := *pb.GetCaptchaS2C{}
	key := strconv.FormatInt(uint64(req.Type), 10) + ":" + req.Key
	v,err := s.Redisc.Get(key)

	if err ==nil && v !=nil{
		captcha := binary.BigEndian.Uint32(v)
		s2c.Captcha = captcha
		return s2c,nil
	}else if err ==nil && v == nil {
		captcha := getRandomNum(100000,999999)
		err = s.Redisc.Set(key,[]byte(captcha))
		mlog.Fatal(err)
		err = s.Redisc.Expire(key,utils.CaptchaExpireSecond)
		mlog.Fatal(err)
		if err !=nil{
			s2c.ErrCode = utils.ErrServer
		}else {
			s2c.Captcha = captcha
		}
		return s2c,nil
	}else{
		s2c.ErrCode= utils.ErrServer
		return s2c,nil
	}
}

func (s *Service ) CheckCaptcha(ctx context.Context, req *pb.CheckCaptchaC2S) (*pb.CheckCaptchaS2C, error){
	mlog.Info("request get captcha with type:%v,key:%v" ,req.Type,req.Key)
	s2c := *pb.CheckCaptchaS2C{}
	key := strconv.FormatInt(uint64(req.Type), 10) + ":" + req.Key
	v,_ := s.Redisc.Get(key)
	if v == nil {
		s2c.ErrCode = utils.ErrCaptchaNotMatch
	}else{
		captcha := binary.BigEndian.Uint32(v)
		if req.Captcha == captcha{
			s2c.ErrCode = 0
		}else{
			s2c.ErrCode = utils.ErrCaptchaNotMatch
		}
	}
	return s2c,nil
}

func getRandomNum(min int,max int ) int{
	return rand.Intn(max- min) + min
}