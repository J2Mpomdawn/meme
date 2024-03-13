package control

import (
	"encoding/json"
	"meme/service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// websocket buffer
var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// pong wait time
var pong_wait = 60 * time.Second

// gvg-websocket streaming process
func GvgHandleConn(c *gin.Context) {
	//convert get to websocket
	conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		service.LogPrintln("red", "Upgrade", err)
	}

	service.FmtPrintln("blue", "start websocket")

	//end pocessing
	defer func() {
		//close communication
		conn.Close()
		service.FmtPrintln("blue", "websocket closed")

		//emppty channels
		c := make(chan bool)
		go func() {
			for len(service.WsGuild) > 0 {
				<-service.WsGuild
			}
			c <- false
		}()
		for len(service.WsGuild) > 0 {
			<-service.WsGuild
		}
		<-c
	}()

	//start
	service.WsFlg = true
	go ws_write(conn)
	ws_read(conn)
}

// gvg-websocket reading process
func ws_read(conn *websocket.Conn) {
	//read setting
	conn.SetReadDeadline(time.Now().Add(pong_wait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pong_wait))
		return nil
	})

	//reading
	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			service.LogPrintln("red", "ReadMessage", err)

			service.WsFlg = false
			break
		}

		switch string(msg) {
		case "\\q", "quit":
			//quit communication
			conn.WriteMessage(t, []byte("closed"))

			service.WsFlg = false
			return

		default:
			//parrot
			service.FmtPrint("blue", "msg: ")
			service.FmtPrintln("blue", string(msg))
			conn.WriteMessage(t, msg)
		}
	}
}

// gvg-websocket writing process
func ws_write(conn *websocket.Conn) {
	//ping-pong timer
	ticker := time.NewTicker(pong_wait * 9 / 10)
	defer ticker.Stop()

	for {
		var b []byte
		var err error

		select {
		case g := <-service.WsGuild:
			//received guild
			b, err = json.Marshal(g)

			service.FmtPrint("blue", "ws_guild: ")
			service.FmtPrintln("blue", g)

			if err != nil {
				service.LogPrintln("red", "Marchal", err)
				return
			}
		case r := <-service.WsRecord:
			//received record
			b, err = json.Marshal(r)

			service.FmtPrint("blue", "ws_record: ")
			service.FmtPrintln("blue", r)

			if err != nil {
				service.LogPrintln("red", "Marchal", err)
				return
			}
		case <-ticker.C:
			//ping
			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			err = conn.WriteMessage(websocket.PingMessage, nil)

			if err != nil {
				service.LogPrintln("red", "WriteMessage", err)
				return
			}
			continue
		}

		conn.WriteMessage(websocket.BinaryMessage, b)
	}
}
