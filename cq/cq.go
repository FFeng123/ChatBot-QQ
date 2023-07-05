package cq

import (
	"time"

	"github.com/apex/log"
	"golang.org/x/net/websocket"
)

var (
	ThisQQ int64

	addr string

	ws *websocket.Conn
)

func Connect(Addr string) {
	addr = Addr
	ReConnect()
}

func DisConnect() {
	if ws != nil {
		ws.Close()
	}
}

func ReConnect() {
	DisConnect()
	for {
		var err error
		ws, err = websocket.Dial(addr, "", addr)
		if err == nil {
			go listen(ws)
			break
		}
		log.Warn("无法连接到CQHttp正在重试···")
		time.Sleep(3 * time.Second)
	}
}
