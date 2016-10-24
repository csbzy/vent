package service

import (
	"golang.org/x/net/context"
	pb "github.com/chenshaobo/vent/business_service/proto"
	"github.com/chenshaobo/vent/business_service/utils"
	"crypto/md5"
	"strconv"
	"encoding/hex"
	"github.com/jbrodriguez/mlog"
)

func (s *Service) Login(ctx context.Context ,req *pb.LoginC2S)(*pb.LoginS2C ,error){
	//do register logic
	mlog.Info("request login:%v, %v",req.PhoneNumber,req.Password)
	res := &pb.LoginS2C{}
	phoneStr := strconv.FormatUint(req.PhoneNumber,10)
	account := utils.AccountPhonePrefix + phoneStr
	if req.PhoneNumber < 9999999999 || req.PhoneNumber > 99999999999{
		res.ErrCode = utils.ErrParams
		return res,nil
	}
	if ! s.Redisc.Exists(account) {
		res.ErrCode = utils.ErrAccountNotExits
		return res,nil
	}

	userID,_ := s.Redisc.Get(account)
	if userID ==nil{
		res.ErrCode = utils.ErrAccountNotExits
		return res,nil
	}
	userIDStr := string(userID[:])
	passwordKey := utils.AccountPasswordPrefix+ userIDStr

	pwdBytes,_ := s.Redisc.Get(passwordKey)

	dbPasswordMD5 := string(pwdBytes[:])
	passwordMD5 := GetMD5Hash(req.Password)
	mlog.Info("db password :%v,login password :%v",dbPasswordMD5,passwordMD5)
	if  dbPasswordMD5  != passwordMD5 {
		res.ErrCode = utils.ErrPasswordWrong
		return res,nil
	}
	userIDInt,err:= strconv.ParseUint(userIDStr,10,64)
	if err != nil{
		res.ErrCode = utils.ErrServer
		return res,nil
	}
	res.UserID = userIDInt
	return res,nil
}

func GetMD5Hash(text string) string{
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}



