package service

import (
	"golang.org/x/net/context"
	pb "github.com/chenshaobo/vent/business_service/proto"
	"github.com/chenshaobo/vent/business_service/utils"
	"github.com/chenshaobo/redisapi"
	"strconv"
)


func (s *Service) UserGeoUpload(ctx context.Context, c2s *pb.GeoUploadC2S) (*pb.CommonS2C, error){
	utils.Info("do upload user geo:%v",c2s)

	err := s.Redisc.GeoAdd("user.geo",redisapi.Coordinate{Longitude:c2s.Longitude,Latitude:c2s.Latitude},strconv.FormatUint(c2s.UserId,10))
	if err !=nil{
		return &pb.CommonS2C{ErrCode:utils.ErrServer},err
	}
	return &pb.CommonS2C{ErrCode:0},nil
}