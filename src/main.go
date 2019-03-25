package main

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

var (
	address = os.Getenv("WEBSOCKET_HOST")+":"+os.Getenv("WEBSOCKET_PORT")
	wsPath = os.Getenv("WEBSOCKET_WS_PATH")
)

func main() {
	// 加载env文件
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// 启动redis
	go func() {
		for message := range SubscribeMessages {
			// 收到消息进行推送
			ChannelsRegister.Broadcast(message)
		}
	}()
	// 启动http服务
	http.HandleFunc(wsPath, NewWebsocket)
	err = http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
	log.Info(address+wsPath)
}