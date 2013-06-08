package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
)

type Connection struct {
	ws        *websocket.Conn
	sendQueue chan string
	chat      ChatStorage
}

const (
	HIST_SIZE int = 5
)

func (conn *Connection) reader() {
	var message string

	for {
		err := websocket.Message.Receive(conn.ws, &message)
		if err != nil {
			fmt.Errorf("Receive error: %s\n", err)
			break
		}

		chatMessage := &Message{Text: message}
		chat.AddMessage(chatMessage)

		j, err := chatMessage.toJSON()
		if err != nil {
			fmt.Errorf("JSON marshaling error: %s\n", err)
			continue
		}
		conn.sendQueue <- j + "\r\n"
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
	fmt.Printf("History %d\n", from)
	for _, message := range conn.chat.Messages[from:] {
		data, err := message.toJSON()
		if err == nil {
			conn.sendQueue <- data + "\r\n"
		}
	}
}
