package consul

import(
	consul "github.com/hashicorp/consul/api"
	"github.com/chenshaobo/vent/business_service/utils"
	"github.com/jbrodriguez/mlog"
)


func Query(consulHost string, key string) ([]byte,error) {
	conf := &consul.Config{Scheme: "http", Address: consulHost}
	client, err := consul.NewClient(conf)
	utils.PrintErr(err)
	kv,_,err := client.KV().Get(key,nil)
	if (err!=nil || kv ==nil){

		mlog.Info("get key error:%v,data :%v",err,kv)
		return nil,err
	}

	return kv.Value,err
}
