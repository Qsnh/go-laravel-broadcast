package main

import (
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	checkOrigin = os.Getenv("WEBSOCKET_CHECK_ORIGIN")
	upgrader    = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// 跨域控制
		CheckOrigin: func(r *http.Request) bool {
			if checkOrigin == "" {
				return false
			}
			if checkOrigin == "*" {
				return true
			}
			origins := strings.Split(checkOrigin, ",")
			requestOrigin := r.Header["Origin"][0]
			log.WithField("origin", requestOrigin).Info("check origin")
			for _, origin := range origins {
				if requestOrigin == origin {
					return true
				}
			}
			return false
		},
	}
)

func NewWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	channelName := r.FormValue("channel")
	if strings.HasPrefix(channelName, "private-") || strings.HasPrefix(channelName, "presence-") {
		if auth := Authorization(channelName, r.Cookies()); auth == false {
			// 断开websocket
			conn.Close()
			return
		}
	}
	// 注册channel到连接的映射
	ChannelsRegister.r[channelName] = append(ChannelsRegister.r[channelName], conn)
	// 统计
	metrics.ClientCount.Inc(1)
}

func HeartbeatTimer() {
	t := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-t.C:
			go HeartbeatHandler()
		}
	}
}

func HeartbeatHandler() {
	for channel, conns := range ChannelsRegister.r {
		for index, conn := range conns {
			if err := conn.WriteMessage(websocket.TextMessage, []byte("hb")); err != nil {
				// clear
				ChannelsRegister.RemoveConn(channel, index)
			}
		}
	}
}
