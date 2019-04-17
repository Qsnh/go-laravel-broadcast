package main

import (
	_ "github.com/joho/godotenv/autoload"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

var (
	address     = os.Getenv("WEBSOCKET_HOST") + ":" + os.Getenv("WEBSOCKET_PORT")
	wsPath      = os.Getenv("WEBSOCKET_WS_PATH")
	tlsEnabled  = os.Getenv("TLS_ENABLED")
	tlsKeyFile  = os.Getenv("TLS_KEY_FILE")
	tlsCertFile = os.Getenv("TLS_CERT_FILE")
)

func main() {
	// 启动redis
	go SubscribeChannel()
	go func() {
		for message := range SubscribeMessages {
			// 收到消息进行推送
			go ChannelsRegister.Broadcast(message)
		}
	}()
	// 数据定时输出
	go metrics.Report()
	// 心跳
	go HeartbeatTimer()
	// 启动http服务
	log.Info(address + wsPath)
	http.HandleFunc(wsPath, func(w http.ResponseWriter, r *http.Request) {
		NewWebsocket(w, r)
	})
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(metrics.GetJson()))
	})
	var err error
	if tlsEnabled == "true" {
		// HTTPS
		err = http.ListenAndServeTLS(address, tlsCertFile, tlsKeyFile, nil)
	} else {
		err = http.ListenAndServe(address, nil)
	}
	if err != nil {
		log.WithField("address", address).Fatal("ListenAndServe:", err)
	}
}
