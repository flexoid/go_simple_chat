package main

import (
	"code.google.com/p/go.net/websocket"
	"net/http"
	"text/template"
)

var (
	chat ChatInstance = ChatInstance{
		broadcast:   make(chan string, 128),
		register:    make(chan *Connection),
		unregister:  make(chan *Connection),
		connections: make(map[*Connection]bool),
	}
	homeTempl = template.Must(template.ParseFiles("home.html"))
)

func handler(w http.ResponseWriter, r *http.Request) {
	homeTempl.Execute(w, r.Host)
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
	http.HandleFunc("/", handler)
	http.Handle("/ws", websocket.Handler(WsServer))
	go chat.run()
	http.ListenAndServe(":8080", nil)
}
