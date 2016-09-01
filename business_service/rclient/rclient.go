package rclient

import(
	"google.golang.org/grpc"
	"time"
	consulb "github.com/chenshaobo/vent/business_service/consul"
	"github.com/chenshaobo/vent/business_service/utils"
	"github.com/uber-go/zap"
	"sync"
)

var (
	_defaultRClients RClients
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
		utils.Logger.Info("register ",zap.String("service name",serviceName))
		conn, err := grpc.Dial(registerHost, grpc.WithInsecure(), grpc.WithBalancer(b), grpc.WithTimeout(time.Second * 10))
		utils.Logger.Info("dial pass")
		utils.PrintErr(err)
		rc.cc[serviceName] = conn
	}
}

func (rc *RClients) get(serviceName string) *grpc.ClientConn{
	if conn,ok := rc.cc[serviceName]; ok{
		return conn
	}else {
		rc.register(serviceName, rc.registerHost)
		utils.Logger.Error("GET service conn fault", zap.String("serviceName", serviceName))
		conn := rc.cc[serviceName]
		return conn
	}
}
