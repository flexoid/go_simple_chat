package main

import (
	"encoding/json"
)

type Message struct {
	Author string
	Text   string
}

func (message *Message) toJSON() (string, error) {
	j, err := json.Marshal(message)
	return string(j), err
}

type ChatInstance struct {
	Messages    []*Message
	connections map[*Connection]bool

	broadcast  chan string
	register   chan *Connection
	unregister chan *Connection
}

func (instance *ChatInstance) AddMessage(message *Message) {
	instance.Messages = append(instance.Messages, message)
}

func (instance *ChatInstance) run() {
	for {
		select {
		case c := <-instance.register:
			instance.connections[c] = true

		case c := <-instance.unregister:
			delete(instance.connections, c)
			close(c.sendQueue)

		case m := <-instance.broadcast:
			for c := range instance.connections {
				c.sendQueue <- m
			}
		}
	}
}
