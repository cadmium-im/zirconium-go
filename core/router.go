package core

import (
	"github.com/ChronosX88/zirconium/core/models"
	"github.com/google/logger"
)

type Router struct {
	moduleManager *ModuleManager
	connections   []*OriginC2S
}

func NewRouter() (*Router, error) {
	mm, err := NewModuleManager()
	if err != nil {
		return nil, err
	}
	r := &Router{
		moduleManager: mm.(*ModuleManager),
	}
	return r, nil
}

func (r *Router) RouteMessage(origin *OriginC2S, message models.BaseMessage) {
	handlers := r.moduleManager.c2sMessageHandlers[message.MessageType]
	if handlers != nil {
		for _, v := range handlers {
			if !v.AnonymousAllowed {
				var entityID, deviceID string
				var isValid bool
				var err error
				if message.AuthToken != "" {
					isValid, entityID, deviceID, err = authManager.ValidateToken(message.AuthToken)
					if err != nil || !isValid {
						logger.Warningf("Connection %s isn't authorized", origin.connID)
						msg := PrepareMessageUnauthorized(message)
						origin.Send(msg)
					}
				} else {
					logger.Warningf("Connection %s isn't authorized", origin.connID)

					msg := PrepareMessageUnauthorized(message)
					origin.Send(msg)
				}

				if origin.entityID == nil {
					origin.entityID = models.NewEntityID(entityID)
				}
				if origin.deviceID == nil {
					origin.deviceID = &deviceID
				}
			}
			go v.HandlerFunc(origin, message)
		}
	} else {
		protocolError := models.ProtocolError{
			ErrCode:    "unhandled",
			ErrText:    "Server doesn't implement message type " + message.MessageType,
			ErrPayload: make(map[string]interface{}),
		}
		errMsg := models.NewBaseMessage(message.ID, message.MessageType, serverDomain, message.From, false, StructToMap(protocolError))
		logger.Infof("Drop message with type %s because server hasn't proper handlers", message.MessageType)
		origin.Send(errMsg)
	}
}

func (r *Router) RoutecoreEvent(sourceModuleName string, eventName string, eventPayload map[string]interface{}) {
	handlers := r.moduleManager.coreEventHandlers[eventName]
	if handlers != nil {
		for _, v := range handlers {
			go v(sourceModuleName, eventPayload)
		}
	} else {
		logger.Infof("Drop event %s because server hasn't proper handlers", eventName)
	}
}
