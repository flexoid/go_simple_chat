package main

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"fmt"
)

type Connection struct {
	ws        *websocket.Conn
	sendQueue chan string
	chat      *ChatInstance
}

const (
	HIST_SIZE int = 5
)

func (conn *Connection) reader() {
	var message string
	var err error

	for {
		err = websocket.Message.Receive(conn.ws, &message)
		if err != nil {
			fmt.Errorf("Receive error: %s\n", err)
			break
		}

		var chatMessage Message
		err = json.Unmarshal([]byte(message), &chatMessage)
		if err != nil {
			fmt.Errorf("JSON unmarchaling error: %s\n", err)
			continue
		}

		chat.AddMessage(&chatMessage)
		conn.chat.broadcast <- message
	}

	conn.ws.Close()
}

func (conn *Connection) writer() {
	for message := range conn.sendQueue {
		err := websocket.Message.Send(conn.ws, message)
		if err != nil {
			fmt.Errorf("Send error: %s\n", err)
			break
		}
	}
}

func (conn *Connection) sendLastMessages() {
	from := 0
	if len(conn.chat.Messages)-HIST_SIZE > 0 {
		from = len(conn.chat.Messages) - HIST_SIZE
	}
	if len(conn.chat.Messages[from:]) > 0 {
		data, err := json.Marshal(conn.chat.Messages[from:])
		if err == nil {
			conn.sendQueue <- string(data) + "\r\n"
		}
	}
}
