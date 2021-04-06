package core

import (
	"github.com/cadmium-im/zirconium-go/core/models"
	"github.com/gorilla/websocket"
)

type Session struct {
	wsConn *websocket.Conn
	connID string
	Claims *JWTCustomClaims
}

func (s *Session) Send(message models.BaseMessage) error {
	return s.wsConn.WriteJSON(message)
}

func (s *Session) Close() error {
	return s.wsConn.Close()
}
