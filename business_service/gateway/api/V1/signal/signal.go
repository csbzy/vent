package signal

import (
	"fmt"
	"github.com/chenshaobo/vent/business_service/utils"
	"github.com/kataras/iris"
	"gopkg.in/square/go-jose.v1/json"
	"time"
	"github.com/kataras/go-websocket"
)

type signal struct {
	SignalType    string `json:"type"` // enter | leave | offer |answer |
	Room          string `json:"room"`
	Sdp           string `json:"sdp"`
	SdpMLineIndex int    `json:"sdpMLineIndex"`
	SdpMid        string `json:"sdpMid"`
	Candidate     string `json:"candidate"`
}

func SetupSignalApi() {
	iris.StaticWeb("/api/v1/signal/", "./api/V1/signal/www/", 3)

	iris.Config.Websocket.Endpoint = "/ws"
	iris.Config.Websocket.MaxMessageSize = 1024000
	iris.Config.Websocket.ReadBufferSize = 1024000
	iris.Config.Websocket.WriteBufferSize = 1024000
	iris.Config.Websocket.WriteTimeout = time.Second * 40
	iris.Config.Websocket.PongTimeout = time.Second * 10
	iris.Config.Websocket.PingPeriod = (iris.Config.Websocket.PongTimeout) * 9 / 10
	ws := iris.Websocket
	ws.OnConnection(func(c iris.WebsocketConnection) {
		s := &signal{}
		var curRoom *string
		c.OnMessage(func(data []byte) {
			fmt.Printf("recive:%v\n", string(data))
			err := json.Unmarshal(data, s)
			if err != nil {
				utils.PrintErr(err)
			} else {
				//@todo auth account before handle busniess
				switch s.SignalType {
				case "enter":
					curRoom = &s.Room
					c.Join(s.Room)
				case "leave":
					c.Leave(*curRoom)
					curRoom = nil
				default:
					c.To(websocket.Broadcast).EmitMessage(data)
				}
			}
		})
		c.OnDisconnect(func() {
			fmt.Printf("\nConnection with ID: %s has been disconnected!", c.ID())
		})
		c.OnError(func(err string) {
			fmt.Printf("receive err:%s\n", err)
		})
	})
}
