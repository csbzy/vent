package geography


import (
	"github.com/kataras/iris"
	"github.com/golang/protobuf/proto"
	pb "github.com/chenshaobo/vent/business_service/proto"
	"github.com/chenshaobo/vent/business_service/rpclient"
	"github.com/chenshaobo/vent/business_service/utils"
	"golang.org/x/net/context"
	"github.com/chenshaobo/vent/business_service/gateway/api.V1/apiUtils"
)

func SetupGeoApi(){
	userParty := iris.Party("/api/v1/coordinate")
	userParty.Put("",apiUtils.AuthSession,GeoUpload)
}



func GeoUpload(c *iris.Context){
	body := c.PostBody()

	c2s := &pb.GeoUploadC2S{}
	s2c := &pb.CommonS2C{}


	err := proto.Unmarshal(body, c2s)
	if err != nil {
		s2c.ErrCode = utils.ErrParams
		apiUtils.SetBody(c,s2c)
		return
	}

	conn := rpclient.Get(utils.ReleationSer)
	if conn == nil {
		s2c.ErrCode = utils.ErrServer
		apiUtils.SetBody(c,s2c)
		return
	}
	rc := pb.NewGeoManagerClient(conn)
	s2c, err = rc.UserGeoUpload(context.Background(), c2s)
	utils.PrintErr(err)
	apiUtils.SetBody(c,s2c)
}