package main

import (
	"github.com/gorilla/websocket"
	"sync"
)

type Channels struct {
	mu sync.Mutex
	r  map[string][]*websocket.Conn
}

func (c *Channels) Broadcast(message ChannelMessage) {
	clients := c.r[message.Channel]
	if len(clients) == 0 {
		return
	}
	for _, client := range clients {
		_ = client.WriteMessage(websocket.TextMessage, message.Data)
	}
}

func (c *Channels) RemoveConn(channel string, index int) {
	c.mu.Lock()
	c.r[channel] = append(c.r[channel][:index], c.r[channel][index+1:]...)
	c.mu.Unlock()
}

var ChannelsRegister = &Channels{
	r: make(map[string][]*websocket.Conn),
}
