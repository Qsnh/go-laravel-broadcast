package main

import (
	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

type ChannelMessage struct {
	Channel string
	Data    []byte
}

var (
	redisHost         = os.Getenv("REDIS_HOST")
	redisPort         = os.Getenv("REDIS_PORT")
	redisPassword     = os.Getenv("REDIS_PASSWORD")
	subscribeChannels = os.Getenv("SUBSCRIBE_CHANNELS")
	SubscribeMessages = make(chan ChannelMessage, 10)
)

func SubscribeChannel() {
	// 连接redis
	dailOption := redis.DialPassword(redisPassword)
	conn, err := redis.Dial("tcp", redisHost+":"+redisPort, dailOption)
	if err != nil {
		log.WithField("host", redisHost).WithField("port", redisPort).Fatal(err)
	}
	defer conn.Close()
	// 订阅
	psc := redis.PubSubConn{Conn: conn}
	if subscribeChannels == "" {
		log.Fatal("无订阅频道")
	}
	channels := strings.Split(subscribeChannels, ",")
	for _, channel := range channels {
		log.WithField("channel", channel).Info("redis subscribe")
		if err := psc.PSubscribe(channel); err != nil {
			log.WithField("redis subscribe channel", channel).Error(err)
		}
	}
	// 监听消息
	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			// 收到订阅消息之后推送到订阅消息chan
			log.Info(v)
			SubscribeMessages <- ChannelMessage{Channel: v.Channel, Data: v.Data}
			//case redis.Subscription:
		case error:
			return
		}
	}
}
