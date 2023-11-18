package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func main() {
	fs := http.FileServer(http.Dir("./static"))
	clients := []*websocket.Conn{}
	http.Handle("/", fs)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		defer c.Close()
		clients = append(clients, c)
		for {
			mt, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
			for _, client := range clients {
				client.WriteMessage(mt, message)
			}
		}
		newClients := make([]*websocket.Conn, len(clients)-1)
		i := 0
		for _, client := range clients {
			if client != c {
				newClients[i] = client
				i += 1
			}
		}
		clients = newClients
	})

	log.Print("Listening on :3001")
	err := http.ListenAndServe(":3001", nil)
	if err != nil {
		log.Fatal(err)
	}
}
