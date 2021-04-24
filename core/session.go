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

func (s *Session) Receive() (models.BaseMessage, error) {
	var msg models.BaseMessage
	err := s.wsConn.ReadJSON(&msg)
	return msg, err
}

func (s *Session) Close() error {
	return s.wsConn.Close()
}
