package main

import (
	"golang.org/x/net/context"
	pb "github.com/chenshaobo/vent/business_service/proto"
	"github.com/uber-go/zap"
)
type Registor struct{

}

//func NewRegistor() Registor{
//	return Registor{}
//}
func (Registor) Register(ctx context.Context ,req *pb.RegisterReq)(*pb.RegisterRes ,error){
	logger.Info("register ",zap.Uint64("phoneNumber",req.PhoneNumber),zap.String("password",req.Password))
	return &pb.RegisterRes{ErrCode:0},nil
}
