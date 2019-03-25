package main

import (
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	channelName := r.FormValue("channel")
	auth := Authorization(channelName, r.Cookies())
	if auth == false {
		// 断开websocket
		conn.Close()
		return
	}
	// 注册channel到连接的映射
	_ = append(ChannelsRegister.r[channelName], conn)
}
