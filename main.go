package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"net/http"
)

var (
	chat ChatInstance = ChatInstance{
		broadcast:   make(chan string, 128),
		register:    make(chan *Connection),
		unregister:  make(chan *Connection),
		connections: make(map[*Connection]bool),
	}
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Access to %s", r.URL)
}

func WsServer(ws *websocket.Conn) {
	connection := &Connection{ws: ws, sendQueue: make(chan string, 128), chat: &chat}
	chat.register <- connection
	defer func() { chat.unregister <- connection }()
	go connection.writer()
	connection.sendLastMessages()
	connection.reader()
}

func main() {
	// http.HandleFunc("/", handler)
	http.Handle("/ws", websocket.Handler(WsServer))
	go chat.run()
	http.ListenAndServe(":8080", nil)
}
