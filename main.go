package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Geridev/socket/websocket/server"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var v = &server.Valami{}

func AutoId() string {
	var err error
	return uuid.Must(uuid.New(), err).String()
}

func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	client := &server.Client{
		Id:   AutoId(),
		Conn: conn,
	}

	v.GetUserName(client)
	v.AddClient(client)
	fmt.Println(len(v.Client))

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			v.RemoveClient(client)
			log.Println(len(v.Client))
			return
		}
		v.SendMessage(p, client)
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/ws", WsEndpoint)
	log.Fatal(http.ListenAndServe(":8080", r))
}
