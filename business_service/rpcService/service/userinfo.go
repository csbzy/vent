package service

import (
	"golang.org/x/net/context"
	pb "github.com/chenshaobo/vent/business_service/proto"
	"github.com/chenshaobo/vent/business_service/utils"
	"strconv"
	"github.com/jbrodriguez/mlog"
	"github.com/chenshaobo/redisapi"
)


func (s *Service) UserInfoModify(ctx context.Context ,req *pb.UserInfoModifyC2S,) (*pb.UserInfoModifyS2C,error){

	mlog.Info("request info modify:%v",*req)

	userInfoKey := utils.UserInfoHashPrefix +  strconv.FormatUint(req.UserID,10)
	res := &pb.UserInfoModifyS2C{}
	if !s.Redisc.Exists(userInfoKey){
		res.ErrCode = utils.ErrAccountNotExits
		return res,nil
	}
	infos := make([]redisapi.ScoreStruct,0)
	if req.Nickname !="" {
		infos = append(infos,redisapi.ScoreStruct{Member:"nickname",Score:req.Nickname})
	}

	if req.City !=""{
		infos = append(infos,redisapi.ScoreStruct{Member:"city",Score:req.City})
	}

	if req.Signature !=""{
		infos = append(infos,redisapi.ScoreStruct{Member:"signature",Score:req.Signature})
	}

	if req.Avatar !=""{
		infos = append(infos,redisapi.ScoreStruct{Member:"avatar",Score:req.Avatar})
	}

	s.Redisc.HMset(userInfoKey,infos)
	res.ErrCode= 0
	return res,nil
}


func (s *Service) UserInfoGet(ctx context.Context,req *pb.UserInfoGetC2S)(*pb.UserInfoGetS2C,error){

	mlog.Info("request get user info:%v",req)
	res := &pb.UserInfoGetS2C{}
	userInfoKey := utils.UserInfoHashPrefix +  strconv.FormatUint(req.TargetUserID,10)
	userInfo,err := s.Redisc.HMget(userInfoKey,"nickname","sex","city","signature")
	if err != nil {
		return res,err
	}
	mlog.Info("user key:%v,userinfo:%v ,err:%v",userInfoKey,userInfo,err)
	res.Nickname = userInfo[0].Score.(string)
	sex,err := strconv.ParseUint(userInfo[1].Score.(string),10,32)
	res.Sex = uint32(sex)
	res.City = userInfo[2].Score.(string)
	res.Signature = userInfo[3].Score.(string)
	mlog.Info("user:%v",res)
	return res,nil
}