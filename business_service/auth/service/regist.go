package service

import (
	"golang.org/x/net/context"
	pb "github.com/chenshaobo/vent/business_service/proto"
	"github.com/chenshaobo/vent/business_service/utils"
	"github.com/jbrodriguez/mlog"
	"strconv"
	"github.com/chenshaobo/redisapi"
)
type Service struct{
	Redisc *redisapi.RedisClient
}

func (s *Service) Register(ctx context.Context ,req *pb.RegisterC2S)(*pb.RegisterS2C ,error){
	//do register logic
	mlog.Info("request register:%v, %v",req.PhoneNumber,req.Password)
	res := &pb.RegisterS2C{}
	phoneStr := strconv.FormatUint(req.PhoneNumber,10)
	if req.PhoneNumber < 9999999999 || req.PhoneNumber > 99999999999{
		res.ErrCode = utils.ErrParams
		return res,nil
	}
	if s.Redisc.Exists(utils.AccountPhonePrefix + phoneStr) {
		res.ErrCode = utils.ErrAccountExits
		return res,nil
	}

	userID,err := s.Redisc.Incr(utils.AccountCount,1)
	if err !=nil {
		res.ErrCode = utils.ErrServer
		return res,nil
	}

	userIDStr := strconv.FormatInt(userID,10)
	err = s.Redisc.Set(utils.AccountPhonePrefix +phoneStr,[]byte(userIDStr))
	if err !=nil {
		res.ErrCode = utils.ErrServer
		return res,nil
	}

	err = s.Redisc.Set(utils.AccountPasswordPrefix+userIDStr,[]byte(GetMD5Hash(req.Password)))
	if err !=nil {
		res.ErrCode = utils.ErrServer
		return res,nil
	}

	err = s.Redisc.Sadd(utils.AccountUserList,userID)
	if err !=nil {
		res.ErrCode = utils.ErrServer
		return res,nil
	}

	sessionKey := utils.AccountSessionPrefix +userIDStr
	sessionStr := genSession()
	mlog.Info("session:%v",sessionStr)
	err = s.Redisc.Set(sessionKey,[]byte(sessionStr))
	if err != nil{
		res.ErrCode = utils.ErrServer
		return res,nil
	}
	s.Redisc.Expire(sessionKey,utils.DaySecond)


	//user info key
	userInfoKey := utils.UserInfoHashPrefix + userIDStr

	//init register user info
	err = s.Redisc.Hset(userInfoKey,"sex",req.Sex)
	err = s.Redisc.Hset(userInfoKey,"nickname",req.PhoneNumber)
	err = s.Redisc.Hset(userInfoKey,"city","深圳")
	err = s.Redisc.Hset(userInfoKey,"signature","今天天气真好")
	if err !=nil {
		res.ErrCode = utils.ErrServer
		return res,nil
	}
	res.ErrCode =0
	res.UserId = uint64(userID)
	res.Session= sessionStr
	return res,nil
}
