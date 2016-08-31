package api

import (
	"github.com/kataras/iris"
	"github.com/golang/protobuf/proto"
	pb "github.com/chenshaobo/vent/business_service/proto"
	"fmt"
)

type RegisterApi struct {
	*iris.Context
}

func (r RegisterApi) Put(){
	body := r.PostBody()
	pbm := &pb.RegisterC2S{}
	proto.Unmarshal(body,pbm)
	fmt.Printf("data%v\n",pbm.PhoneNumber)
	r.Write("ok")
}

func (r RegisterApi) Get(){
	body := r.PostBody()
	pbm := &pb.RegisterC2S{}
	proto.Unmarshal(body,pbm)
	fmt.Printf("data%v\n",pbm.PhoneNumber)
	r.Write("ok")
}