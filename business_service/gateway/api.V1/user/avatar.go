package user

//import (
//	"github.com/kataras/iris"
//	pb "github.com/chenshaobo/vent/business_service/proto"
//	"github.com/chenshaobo/vent/business_service/rpclient"
//	"github.com/golang/protobuf/proto"
//	"golang.org/x/net/context"
//	"github.com/chenshaobo/vent/business_service/utils"
//	"strconv"
//	"github.com/jbrodriguez/mlog"
//	"github.com/chenshaobo/vent/business_service/gateway/api.V1/apiUtils"
//	"mime/multipart"
//	"github.com/valyala/fasthttp"
//)
//
//func UpdateAvatar(c *iris.Context){
//	useIDInt,err := c.ParamInt64("userID")
//
//	f := c.FormFile("avatar")
//	//f.Header
//	fasthttp.SaveMultipartFile(f,"")
//	c2s := &pb.UserInfoModifyC2S{}
//	s2c := &pb.UserInfoModifyS2C{}
//	if err !=nil{
//		s2c.ErrCode = utils.ErrServer
//		mlog.Error(err)
//		apiUtils.SetBody(c,s2c)
//	}
//
//
//
//}
