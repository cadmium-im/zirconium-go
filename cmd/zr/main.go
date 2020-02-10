package main

import (
	"log"
	"net/http"

	"github.com/ChronosX88/zirconium/internal"
	"github.com/google/logger"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var connectionHandler = internal.NewConnectionHandler()
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	internal.InitializeContext("localhost")
	router := mux.NewRouter()
	router.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		response.Write([]byte("Zirconium server is up and running!"))
	}).Methods("GET")
	router.HandleFunc("/ws", wsHandler)

	logger.Info("Zirconium successfully started!")
	logger.Fatal(http.ListenAndServe(":8844", router))
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	connectionHandler.HandleNewConnection(ws)
}
