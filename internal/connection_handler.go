package internal

import (
	"github.com/ChronosX88/zirconium/internal/models"
	"github.com/google/logger"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type ConnectionHandler struct {
	connections map[string]*OriginC2S
}

func NewConnectionHandler() *ConnectionHandler {
	return &ConnectionHandler{
		connections: make(map[string]*OriginC2S),
	}
}

func (ch *ConnectionHandler) HandleNewConnection(wsocket *websocket.Conn) {
	uuid, _ := uuid.NewRandom()
	uuidStr := uuid.String()
	o := &OriginC2S{
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
				o.wsConn.Close()
				logger.Infof("Connection %s was closed. (Reason: %s)", o.connID, err.Error())
				break
			}
			router.RouteMessage(o, msg)
		}
	}()
}
