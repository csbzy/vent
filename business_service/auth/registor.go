package main

import (
	"golang.org/x/net/context"
	pb "github.com/chenshaobo/vent/business_service/proto"
	"github.com/uber-go/zap"
	"github.com/chenshaobo/vent/business_service/utils"
)
type Registor struct{

}

//func NewRegistor() Registor{
//	return Registor{}
//}
func (Registor) Register(ctx context.Context ,req *pb.RegisterReq)(*pb.RegisterRes ,error){
	logger.Info("register ",zap.Uint64("phoneNumber",req.PhoneNumber),zap.String("password",req.Password))
	//do register logic
	res := &pb.RegisterRes{}
	if rdc.Exists(utils.AccountPhonePrefix) {
		res.ErrCode = 10001
		return res
	}

	userID,err := rdc.Incr(utils.AccountCount,1)
	utils.PrintErr(err)

	res.ErrCode =0
	return res,nil
}
