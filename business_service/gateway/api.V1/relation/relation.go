package relation
import (
	"github.com/kataras/iris"
	"github.com/golang/protobuf/proto"
	pb "github.com/chenshaobo/vent/business_service/proto"
	"github.com/chenshaobo/vent/business_service/rpclient"
	"github.com/chenshaobo/vent/business_service/utils"
	"golang.org/x/net/context"
	"github.com/chenshaobo/vent/business_service/gateway/api.V1/apiUtils"
	//"github.com/jbrodriguez/mlog"
)

func SetupRelationApi(){
	userParty := iris.Party("/api/v1/relation")
	userParty.Get("/recentContact",apiUtils.AuthSession,GetFriends)
}



func GetFriends(c *iris.Context){
	body := c.PostBody()

	c2s := &pb.RecentContactGetC2S{}
	s2c := &pb.RecentContactGetS2C{}


	err := proto.Unmarshal(body, c2s)
	if err != nil {
		s2c.ErrCode = utils.ErrParams
		apiUtils.SetBody(c,s2c)
		return
	}

	conn := rpclient.Get(utils.RelationSer)
	if conn == nil {
		s2c.ErrCode = utils.ErrServer
		apiUtils.SetBody(c,s2c)
		return
	}
	rc := pb.NewRelationClient(conn)
	s2c, err = rc.RecentContactGet(context.Background(), c2s)

	apiUtils.SetBody(c,s2c)
}