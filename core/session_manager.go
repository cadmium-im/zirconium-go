package core

import (
	"github.com/cadmium-im/zirconium-go/core/models"
	"github.com/google/logger"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type SessionManager struct {
	router      *Router
	connections map[string]*Session
}

func NewSessionManager(r *Router) *SessionManager {
	return &SessionManager{
		router:      r,
		connections: make(map[string]*Session),
	}
}

func (ch *SessionManager) HandleNewConnection(wsocket *websocket.Conn) {
	randomUUID := uuid.New().String()
	o := &Session{
		wsConn: wsocket,
		connID: randomUUID,
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
