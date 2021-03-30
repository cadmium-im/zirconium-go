package core

import (
	"github.com/cadmium-im/zirconium-go/core/models"
	"github.com/google/logger"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type ConnectionHandler struct {
	router      *Router
	connections map[string]*Session
}

func NewConnectionHandler(r *Router) *ConnectionHandler {
	return &ConnectionHandler{
		router:      r,
		connections: make(map[string]*Session),
	}
}

func (ch *ConnectionHandler) HandleNewConnection(wsocket *websocket.Conn) {
	randomUUID, _ := uuid.NewRandom()
	uuidStr := randomUUID.String()
	o := &Session{
		wsConn: wsocket,
		connID: uuidStr,
	}
	ch.connections[o.connID] = o
	go func() {
		for {
			var msg models.BaseMessage
			err := o.wsConn.ReadJSON(&msg)
			if err != nil {
				delete(ch.connections, o.connID)
				_ = o.wsConn.Close()
				logger.Infof("Connection %s was closed. (Reason: %s)", o.connID, err.Error())
				break
			}
			ch.router.RouteMessage(o, msg)
		}
	}()
	logger.Infof("Connection %s was created", o.connID)
}
