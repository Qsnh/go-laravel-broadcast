package main

import "github.com/gorilla/websocket"

type Channel struct {
	Name       string
	Broadcasts chan string
	Clients    []*websocket.Conn
}

type Channels struct {
	r map[string][]*websocket.Conn
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

var ChannelsRegister = &Channels{
	r: make(map[string][]*websocket.Conn),
}
