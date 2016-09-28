package service

import (
	"golang.org/x/net/context"
	pb "github.com/chenshaobo/vent/business_service/proto"
	"github.com/chenshaobo/vent/business_service/utils"
	"strconv"
)


func (s *Service) UserInfoModify(ctx context.Context ,req *pb.UserInfoModifyC2S,) (*pb.UserInfoModifyS2C,error){

	utils.Info("request info modify:%v",*req)

	userInfoKey := utils.UserInfoHashPrefix +  strconv.FormatUint(req.UserId,10)
	res := &pb.UserInfoModifyS2C{}
	if !s.Redisc.Exists(userInfoKey){
		res.ErrCode = utils.ErrAccountNotExits
		return res,nil
	}
	if req.Nickname !="" {
		s.Redisc.Hset(userInfoKey,"nickname",req.Nickname)
	}

	if req.City !=""{
		s.Redisc.Hset(userInfoKey,"city",req.City)
	}

	if req.Signature !=""{
		s.Redisc.Hset(userInfoKey,"signature",req.Signature)
	}


	res.ErrCode= 0
	return res,nil
}


func (s *Service) UserInfoGet(ctx context.Context,req *pb.UserInfoGetC2S)(*pb.UserInfoGetS2C,error){

	utils.Info("request get user info:%v",req)
	userInfoKey := utils.UserInfoHashPrefix +  strconv.FormatUint(req.UserId,10)
	userInfo,err := s.Redisc.HMget(userInfoKey,"nickname","sex","city","signature")
	utils.Info("userinfo:%v ,err:%v",userInfo,err)

	res := &pb.UserInfoGetS2C{}
	return res,nil
}