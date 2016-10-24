package rpclient

import(
	"google.golang.org/grpc"
	"time"
	consulb "github.com/chenshaobo/vent/business_service/consul"
	"github.com/chenshaobo/vent/business_service/utils"
	"sync"
	"github.com/jbrodriguez/mlog"
)

var (
	_defaultRClients RClients //var will Automatic initialization to zero value
	once sync.Once
)

//export API
func Init(registerHost string){
	once.Do(func(){ _defaultRClients.init(registerHost)})
}


func Register(serviceName string){
	if _defaultRClients.registerHost !=""{
		_defaultRClients.register(serviceName,_defaultRClients.registerHost)
	}
}

func Get(serviceName string) *grpc.ClientConn{
	return _defaultRClients.get(serviceName)
}

type RClients struct{
	cc map[string] *grpc.ClientConn
	registerHost string
}


func (rc *RClients) init(registerHost string){
	rc.registerHost = registerHost
	rc.cc = map[string] *grpc.ClientConn{}
}
//register use to get a connect from register for serviceName
func (rc *RClients) register(serviceName string,registerHost string) {
	if _,ok := rc.cc[serviceName]; !ok {
		r := consulb.NewResolver(serviceName)
		b := grpc.RoundRobin(r)
		mlog.Info("register:  service name:%v,register host:%v",serviceName,registerHost)
		conn, err := grpc.Dial(registerHost, grpc.WithInsecure(), grpc.WithBalancer(b), grpc.WithTimeout(time.Second * 10))
		mlog.Info("get conn ok")
		utils.PrintErr(err)
		rc.cc[serviceName] = conn
	}
}

func (rc *RClients) get(serviceName string) *grpc.ClientConn{
	if conn,ok := rc.cc[serviceName]; ok{
		return conn
	}else {
		rc.register(serviceName, rc.registerHost)
		conn := rc.cc[serviceName]
		return conn
	}
}
