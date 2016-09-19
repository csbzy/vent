package signal

import (
	"github.com/kataras/iris"
	"fmt"
	"time"
)

func SetupSignalApi(){
	iris.StaticWeb("/api/v1/signal/", "./api/V1/signal/www/",3)

	ws := iris.Websocket
	iris.Config.Websocket.Endpoint = "/ws"
	iris.Config.Websocket.MaxMessageSize = 10240
	iris.Config.Websocket.ReadBufferSize = 10240
	iris.Config.Websocket.WriteBufferSize = 10240
	iris.Config.Websocket.WriteTimeout = time.Hour
	iris.Config.Websocket.PongTimeout = time.Hour
	ws.OnConnection(func(c iris.WebsocketConnection){
		c.OnMessage(func(data []byte) {
			fmt.Printf("recive:%v\n",string(data))
			//c.To(iris.Broadcast).EmitMessage(data)
			c.EmitMessage(data)
		})
		c.OnDisconnect(func() {
			fmt.Printf("\nConnection with ID: %s has been disconnected!", c.ID())
		})
		c.OnError(func(err string){
			fmt.Printf("receive err:%s\n",err)
		})
	})
}