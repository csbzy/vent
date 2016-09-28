package service

import (
	"golang.org/x/net/context"
	pb "github.com/chenshaobo/vent/business_service/proto"
	"github.com/chenshaobo/vent/business_service/utils"
)


func (s *Service) UserGeoUpload(ctx context.Context, c2s *pb.GeoUploadC2S) (*pb.CommonS2C, error){
	utils.Info("do upload user geo:%v",c2s)

	return &pb.CommonS2C{ErrCode:0},nil
}