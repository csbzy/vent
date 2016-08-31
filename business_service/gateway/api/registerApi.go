package api

import (
	"github.com/kataras/iris"
	"github.com/golang/protobuf/proto"
	pb "github.com/chenshaobo/vent/business_service/proto"
)

type RegisterApi struct {
	*iris.Context
}

func (r *RegisterApi) Put(){
	body := r.PostBody()
	pbm := &pb.RegisterS2C{}
	proto.Unmarshal(body,pbm)

}