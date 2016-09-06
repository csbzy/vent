package main

import (
	"golang.org/x/net/context"
	pb "github.com/chenshaobo/vent/business_service/proto"
	"github.com/chenshaobo/vent/business_service/utils"
	"github.com/jbrodriguez/mlog"
	"strconv"
)
type Registor struct{

}

//func NewRegistor() Registor{
//	return Registor{}
//}
func (Registor) Register(ctx context.Context ,req *pb.RegisterC2S)(*pb.RegisterS2C ,error){
	//do register logic
	mlog.Info("request register:%v, %v",req.PhoneNumber,req.Password)
	res := &pb.RegisterS2C{}
	phoneStr := strconv.FormatUint(req.PhoneNumber,10)
	if req.PhoneNumber < 9999999999 || req.PhoneNumber > 99999999999{
		res.ErrCode = utils.ErrParams
		return res,nil
	}
	if rdc.Exists(utils.AccountPhonePrefix + phoneStr) {
		res.ErrCode = utils.ErrAccountExits
		return res,nil
	}

	userID,err := rdc.Incr(utils.AccountCount,1)
	if err !=nil {
		res.ErrCode = utils.ErrServer
		return res,nil
	}

	userIDStr := strconv.FormatInt(userID,10)
	err = rdc.Set(utils.AccountPhonePrefix +phoneStr,[]byte(userIDStr))
	if err !=nil {
		res.ErrCode = utils.ErrServer
		return res,nil
	}

	err = rdc.Set(utils.AccountPasswordPrefix+userIDStr,[]byte(req.Password))
	if err !=nil {
		res.ErrCode = utils.ErrServer
		return res,nil
	}

	err = rdc.Sadd(utils.AccountUserList,userID)
	if err !=nil {
		res.ErrCode = utils.ErrServer
		return res,nil
	}
	res.ErrCode =0
	res.UserId = uint64(userID)

	return res,nil
}
