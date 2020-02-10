package internal

import (
	"github.com/ChronosX88/zirconium/internal/models"
	"github.com/gorilla/websocket"
)

type OriginC2S struct {
	wsConn   *websocket.Conn
	connID   string
	entityID *models.EntityID
	deviceID *string
}

func (o *OriginC2S) Send(message models.BaseMessage) error {
	return o.wsConn.WriteJSON(message)
}
