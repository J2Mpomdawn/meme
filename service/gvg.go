package service

import (
	"flag"
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"

	"meme/analysis"
	"meme/model"
)

var addr = flag.String("addr", "api.mentemori.icu", "http service address")
var Current_sub model.Value_StreamId
var Buffer = make(chan []byte)
var ReqFlg = make(chan bool)
var ExtFlg = make(chan bool)
var recvErrFlg = false

func SetCurrentSub() {
	//渡す項目はenvから取得する
	Current_sub = model.Value_StreamId{
		WorldId:  1099,
		GroupId:  0,
		Class:    0,
		Block:    0,
		CastleId: 0,
	}
}

func GetBuffer() []byte {
	return analysis.PackStreamId(Current_sub)
}

func Gvg() {
	flag.Parse()
	log.SetFlags(0)

	u := url.URL{Scheme: "wss", Host: *addr, Path: "/gvg"}
	log.Printf("connecting to %s", u.String())

	c, res, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err, "\n", res.StatusCode)
	}

	defer func() {
		log.Println("exit [Gvg]")
		c.Close()

		if recvErrFlg {
			recvErrFlg = false

			go Gvg()
			b := make([]byte, 4)
			b[0] = 0
			b[1] = 0
			b[2] = 88
			b[3] = 34
			Buffer <- b
			<-ReqFlg
		}
	}()

	done := make(chan struct{})

	go recv(c, done)
	wait(c, done)
}

func recv(c *websocket.Conn, done chan struct{}) {
	defer func() {
		log.Println("exit [recv]")
		close(done)
	}()

	go func() {
		defer log.Println("exit [recv-for]")
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)

				recvErrFlg = true
				return
			}

			analysis.GvgAnalysis(message, Current_sub)
			RegisterGuild(analysis.Guilds)
			RegisterRecord(analysis.Castles)
		}
	}()
	<-ExtFlg
}

func wait(c *websocket.Conn, done chan struct{}) {
	/*
	 * Time allowed to read the next pong message from the peer
	 * pongWait = 60
	 *
	 * Send pings to peer with this period. Must be less than pongWait.
	 * pingPeriod = (pongWait * 9) / 10
	 *            = (60 * 9) / 10
	 *            = 540 / 10
	 *            = 54
	 */
	ticker := time.NewTicker(50 * time.Second)
	defer ticker.Stop()
	defer log.Println("exit [wait]")

	for {
		select {
		case <-done:
			log.Println("done")
			return
		case b := <-Buffer:
			w, err := c.NextWriter(websocket.BinaryMessage)
			if err != nil {
				log.Println("writer b:", err)
			}
			_, err = w.Write(b)
			w.Close()
			if err != nil {
				log.Println("write b:", err)
				return
			}

			log.Println(b)
			ReqFlg <- false
		case <-ticker.C:
			err := c.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				log.Println("write t:", err)
				return
			}
			log.Println("tic")
		}
	}
}
