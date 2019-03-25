package main

import (
	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
	"strings"
)

type ChannelMessage struct {
	Channel string
	Data    []byte
}

var (
	redisHost         = GetEnv("REDIS_HOST")
	redisPort         = GetEnv("REDIS_PORT")
	redisPassword     = GetEnv("REDIS_PASSWORD")
	subscribeChannels = GetEnv("SUBSCRIBE_CHANNELS")
	SubscribeMessages = make(chan ChannelMessage, 10)
)

func init() {
	// 连接redis
	dailOption := redis.DialPassword(redisPassword)
	conn, err := redis.Dial("tcp", redisHost+":"+redisPort, dailOption)
	if err != nil {
		log.WithField("host", redisHost).WithField("port", redisPort).Fatal(err)
	}
	defer conn.Close()
	// 订阅
	psc := redis.PubSubConn{Conn: conn}
	channels := strings.Split(subscribeChannels, ",")
	if len(channels) == 0 {
		log.Fatal("无订阅频道")
	}
	for _, channel := range channels {
		psc.Subscribe(channel)
	}
	// 监听消息
	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			// 收到订阅消息之后推送到订阅消息chan
			SubscribeMessages <- ChannelMessage{Channel: v.Channel, Data: v.Data}
			//case redis.Subscription:
			//case error:
		}
	}
}
