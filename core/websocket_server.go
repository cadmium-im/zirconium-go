package core

import (
	"fmt"
	"github.com/google/logger"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"net/http"
)

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebsocketServer struct {
	r              *mux.Router
	sessionManager *SessionManager
	cfg            *Config
}

func NewWebsocketServer(cfg *Config, sessionManager *SessionManager) *WebsocketServer {
	wss := &WebsocketServer{}

	wss.sessionManager = sessionManager
	wss.cfg = cfg
	r := mux.NewRouter()
	wss.r = r
	r.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		_, _ = response.Write([]byte("Zirconium server is up and running!"))
	}).Methods("GET")
	r.HandleFunc(cfg.Websocket.Endpoint, func(w http.ResponseWriter, r *http.Request) {
		ws, err := wsUpgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.Errorf(err.Error())
			return
		}
		wss.sessionManager.HandleNewConnection(ws)
	})

	return wss
}

func (wss *WebsocketServer) Run() error {
	addr := fmt.Sprintf("%s:%d", wss.cfg.Websocket.Host, wss.cfg.Websocket.Port)

	return http.ListenAndServe(addr, wss.r)
}
