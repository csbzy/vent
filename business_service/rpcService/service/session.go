package service

import (
	"golang.org/x/net/context"
	pb "github.com/chenshaobo/vent/business_service/proto"
	"github.com/chenshaobo/vent/business_service/utils"
	"github.com/satori/go.uuid"
	"strconv"
	"github.com/jbrodriguez/mlog"
)

func (s *Service) AuthSession(ctx context.Context, c2s *pb.AuthReq) (*pb.CommonS2C, error){
	res := &pb.CommonS2C{}
	userIDStr := strconv.FormatUint(c2s.UserId,10)
	sessionKey := utils.AccountSessionPrefix +userIDStr
	sessionByte ,err := s.Redisc.Get(sessionKey)
	if err !=nil || len(sessionByte) == 0 || c2s.Session != string(sessionByte[:]) {

		res.ErrCode = utils.ErrSessionNotMatch
	}
	mlog.Info("session:%v ,redis session:%v",c2s.Session,string(sessionByte[:]))
	return res,nil
}


func (s *Service) GetSession(ctx context.Context, c2s *pb.GetSessionReq) (*pb.GetSessionRes, error){
	res := &pb.GetSessionRes{}
	userIDStr := strconv.FormatUint(c2s.UserId,10)
	sessionKey := utils.AccountSessionPrefix +userIDStr
	sessionByte ,err := s.Redisc.Get(sessionKey)
	if err !=nil || len(sessionByte) == 0  {
		sessionByte = []byte(genSession())
		err = s.Redisc.Set(sessionKey,sessionByte)
		if err !=nil{
			mlog.Error(err)
		}
	}
	s.Redisc.Expire(sessionKey,utils.DaySecond)

	res.UserId = c2s.UserId
	res.Session = string(sessionByte[:])

	return res,nil
}



func genSession() string{
	return uuid.NewV4().String()
}