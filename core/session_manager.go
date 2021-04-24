package core

import (
	"github.com/cadmium-im/zirconium-go/core/models"
	"github.com/cadmium-im/zirconium-go/core/utils"
	"github.com/google/logger"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type SessionManager struct {
	domains     []string
	router      *Router
	connections map[string]*Session
}

func NewSessionManager(r *Router, domains []string) *SessionManager {
	return &SessionManager{
		domains:     domains,
		router:      r,
		connections: make(map[string]*Session),
	}
}

func (sm *SessionManager) HandleNewConnection(wsocket *websocket.Conn) {
	randomUUID := uuid.New().String()
	s := &Session{
		wsConn: wsocket,
		connID: randomUUID,
	}

	msg, err := s.Receive()
	if err != nil {
		logger.Infof("Error occurred when tried to receive first message! Closing connection... (%s)", err.Error())
		_ = s.wsConn.Close()
		return
	}
	if msg.MessageType != "urn:cadmium:connection:open" {
		msg.MessageType = "urn:cadmium:connection"
		emsg := utils.PrepareErrorMessage(msg, "invalid-conn-negotiation", "", "")
		s.Send(emsg)
		s.wsConn.Close()
		return
	}

	eid, err := models.NewEntityIDFromString(msg.To)
	if err != nil {
		emsg := utils.PrepareErrorMessage(msg, "invalid-eid", "", "")
		s.Send(emsg)
		s.wsConn.Close()
		return
	}
	if eid.OnlyServerPart && !utils.InStringArray(eid.ServerPart, sm.domains) {
		emsg := utils.PrepareErrorMessage(msg, "unknown-host", "", "")
		s.Send(emsg)
		s.wsConn.Close()
		return
	}

	replyMessage := models.NewEmptyBaseMessage()
	replyMessage.ID = msg.ID
	replyMessage.From = msg.To
	replyMessage.MessageType = "urn:cadmium:connection:open"
	replyMessage.Payload["id"] = randomUUID
	s.Send(replyMessage)

	sm.connections[s.connID] = s
	go func() {
		for {
			msg, err := s.Receive()
			if err != nil {
				delete(sm.connections, s.connID)
				_ = s.wsConn.Close()
				logger.Infof("Connection %s was closed. (Reason: %s)", s.connID, err.Error())
				break
			}
			eid, err := models.NewEntityIDFromString(msg.To)
			if err != nil {
				emsg := utils.PrepareErrorMessage(msg, "invalid-eid", "", "")
				s.Send(emsg)
				continue
			}
			if eid.LocalPart == "" && utils.InStringArray(eid.ServerPart, sm.domains) {
				emsg := utils.PrepareErrorMessage(msg, "unknown-host", "", "")
				s.Send(emsg)
				continue
			}
			sm.router.RouteMessage(s, msg)
		}
	}()
	logger.Infof("Connection %s was created", s.connID)
}
