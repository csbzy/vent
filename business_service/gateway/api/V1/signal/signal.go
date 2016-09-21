package signal

import (
	"github.com/kataras/iris"
	"fmt"
	"time"
	"gopkg.in/square/go-jose.v1/json"
	"github.com/chenshaobo/vent/business_service/utils"
)

type signal struct{
	SignalType string `json:"type"`
	Room       string `json:"room"`
	Sdp        string `json:"sdp"`
	SdpMLineIndex int `json:"sdpMLineIndex"`
	SdpMid     string  `json:"sdpMid"`
	Candidate  string  `json:"candidate"`
}

func SetupSignalApi(){
	iris.StaticWeb("/api/v1/signal/", "./api/V1/signal/www/",3)


	iris.Config.Websocket.Endpoint = "/ws"
	iris.Config.Websocket.MaxMessageSize =  1024000
	iris.Config.Websocket.ReadBufferSize =  1024000
	iris.Config.Websocket.WriteBufferSize = 1024000
	iris.Config.Websocket.WriteTimeout = time.Second * 40
	iris.Config.Websocket.PongTimeout = time.Second *10
	iris.Config.Websocket.PingPeriod = (iris.Config.Websocket.PongTimeout)* 9 /10
	ws := iris.Websocket
	s := &signal{}
	ws.OnConnection(func(c iris.WebsocketConnection){
		c.OnMessage(func(data []byte) {
			fmt.Printf("recive:%v\n",string(data))
			err := json.Unmarshal(data,s)
			if  err !=nil {
				utils.PrintErr(err)
			}else{
				switch s.SignalType {
				case "enter":
					c.Join(s.Room)
				case "leave":
					c.Leave(s.Room)
				default:

				}
				c.EmitMessage(data)
			}
			//c.To(iris.Broadcast).EmitMessage(data)
			//c.EmitMessage(data)
		})
		c.OnDisconnect(func() {
			fmt.Printf("\nConnection with ID: %s has been disconnected!", c.ID())
		})
		c.OnError(func(err string){
			fmt.Printf("receive err:%s\n",err)
		})
	})
}
