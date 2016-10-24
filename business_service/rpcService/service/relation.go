package service

import (
	"golang.org/x/net/context"
	pb "github.com/chenshaobo/vent/business_service/proto"
	"github.com/chenshaobo/vent/business_service/utils"

	"github.com/chenshaobo/vent/business_service/rpclient"
	"strconv"
	"github.com/jbrodriguez/mlog"
)


func (s *Service) RecentContactGet(ctx context.Context, c2s *pb.RecentContactGetC2S) (*pb.RecentContactGetS2C, error){
	mlog.Info("do upload user geo:%v",c2s)
	s2c := &pb.RecentContactGetS2C{}
	userIDStr := strconv.FormatUint(c2s.UserID,10)
	recentContacts,err :=s.Redisc.Lrange(utils.RecentContactPrefix+userIDStr,0,-1)
	if err!=nil{
		mlog.Error(err)
		return &pb.RecentContactGetS2C{ErrCode:utils.ErrServer},nil
	}
	contacts := []*pb.Friend {}
	userInfoGet := &pb.UserInfoGetC2S{}

	conn := rpclient.Get(utils.RegisterSer)
	if conn == nil {
		s2c.ErrCode = utils.ErrServer
		return s2c,nil
	}
	rc := pb.NewUserInfoManagerClient(conn)


	for _,contactID := range recentContacts{
		//@对每一个最近联系人进行一次grpc请求,可以优化为合并请求
		userInfoGet.TargetUserID,err = strconv.ParseUint(contactID.(string),10,64)
		if err != nil {
			mlog.Error(err)
			continue
		}
		userInfoGets2c, err := rc.UserInfoGet(context.Background(), userInfoGet)
		if err !=nil{
			continue
		}
		contacts = append(contacts,&pb.Friend{UserID:userInfoGets2c.UserID,Signature:userInfoGets2c.Signature})

	}

	return &pb.RecentContactGetS2C{ErrCode:0,Friends:contacts},nil
}