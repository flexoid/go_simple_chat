package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"net/http"
)

var (
	chat ChatStorage
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Access to %s", r.URL)
}

func WsServer(ws *websocket.Conn) {
	connection := &Connection{ws: ws, sendQueue: make(chan string, 128), chat: chat}
	go connection.writer()
	connection.sendLastMessages()
	connection.reader()
}

func main() {
	// http.HandleFunc("/", handler)
	http.Handle("/ws", websocket.Handler(WsServer))
	http.ListenAndServe(":8080", nil)
}
