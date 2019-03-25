package main

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

var (
	address = GetEnv("WEBSOCKET_HOST") + ":" + GetEnv("WEBSOCKET_PORT")
	wsPath  = GetEnv("WEBSOCKET_WS_PATH")
)

func main() {
	// 启动redis
	go func() {
		for message := range SubscribeMessages {
			// 收到消息进行推送
			ChannelsRegister.Broadcast(message)
		}
	}()
	// 启动http服务
	http.HandleFunc(wsPath, NewWebsocket)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
	log.Info(address + wsPath)
}
