package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]string)
var clientsReverse = make(map[string]*websocket.Conn)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		response.Write([]byte("Zirconium server is up and running!"))
	}).Methods("GET")
	router.HandleFunc("/ws", wsHandler)

	log.Println("Zirconium successfully started!")
	log.Fatal(http.ListenAndServe(":8844", router))
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	// register client
	clients[ws] = uuid.New().String()
	clientsReverse[clients[ws]] = ws
	go readLoop(ws)
	log.Printf("Connection %s created!", clients[ws])
}

func readLoop(c *websocket.Conn) {
	for {
		if _, _, err := c.NextReader(); err != nil {
			connectionID := clients[c]
			if connectionID != "" {
				delete(clients, c)
				delete(clientsReverse, connectionID)
				log.Printf("Connection %s closed!", connectionID)
			} else {
				log.Println("connection wasn't found")
			}
			c.Close()
			break
		}
	}
}
