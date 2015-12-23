package main

import (
	"github.com/gorilla/websocket"
	"net/http"
)

var connections = make(map[*websocket.Conn]bool)
var drops = make(chan Drop)

func websocketHandler(w http.ResponseWriter, r *http.Request) {

	conn, _ := websocket.Upgrade(w, r, nil, 1024, 1024)
	addConnection(conn)
	listen(conn)

}

func addConnection(conn *websocket.Conn) {
	connections[conn] = true
}

func removeConnection(conn *websocket.Conn) {
	delete(connections, conn)
}

func sendDrops(d Drop) {
	for c := range connections {
		c.WriteJSON(d)
	}
}

func listen(conn *websocket.Conn) {

	defer removeConnection(conn)
	for {
		//_, msg, err := conn.ReadMessage()

		var a Drop

		err := conn.ReadJSON(&a)

		if err != nil {
			break
		}
		drops <- a
	}

}

func updateListen(d DB) {

	for {
		select {
		case drop := <-drops:
			d.putDrops(drop)
			sendDrops(drop)
		}

	}

}
