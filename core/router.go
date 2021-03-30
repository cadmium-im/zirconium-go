package core

import (
	"github.com/cadmium-im/zirconium-go/core/models"
	"github.com/google/logger"
)

type Router struct {
	appContext  *AppContext
	handlers    map[string][]C2SMessageHandler
	connections []*Session
}

type C2SMessageHandler interface {
	HandleMessage(s *Session, message models.BaseMessage)
	IsAuthorizationRequired() bool
	HandlingType() string
}

func NewRouter(ctx *AppContext) (*Router, error) {
	r := &Router{
		appContext: ctx,
		handlers:   map[string][]C2SMessageHandler{},
	}
	return r, nil
}

func (r *Router) RouteMessage(origin *Session, message models.BaseMessage) {
	handlers := r.handlers[message.MessageType]
	if handlers != nil {
		for _, v := range handlers {
			if v.IsAuthorizationRequired() {
				if len(origin.entityID) == 0 {
					logger.Warningf("Connection %s isn't authorized", origin.connID)

					msg := PrepareMessageUnauthorized(message, r.appContext.cfg.ServerDomains[0]) // fixme: domain
					_ = origin.Send(msg)
				}
			}
			go v.HandleMessage(origin, message)
		}
	} else {
		protocolError := models.ProtocolError{
			ErrCode:    "unhandled",
			ErrText:    "Server doesn't implement message type " + message.MessageType,
			ErrPayload: make(map[string]interface{}),
		}
		errMsg := models.NewBaseMessage(message.ID, message.MessageType, r.appContext.cfg.ServerID, []string{message.From}, false, StructToMap(protocolError))
		logger.Infof("Drop message with type %s because server hasn't proper handlers", message.MessageType)
		_ = origin.Send(errMsg)
	}
}

func (r *Router) RegisterC2SHandler(c C2SMessageHandler) {
	r.handlers[c.HandlingType()] = append(r.handlers[c.HandlingType()], c)
}
