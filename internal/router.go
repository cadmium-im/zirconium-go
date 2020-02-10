package internal

import (
	"github.com/ChronosX88/zirconium/internal/models"
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
		moduleManager: mm,
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
		logger.Infof("Drop message with type %s because server hasn't proper handlers", message.MessageType)
	}
}

func (r *Router) RouteInternalEvent(sourceModuleName string, eventName string, eventPayload map[string]interface{}) {
	handlers := r.moduleManager.internalEventHandlers[eventName]
	if handlers != nil {
		for _, v := range handlers {
			go v(sourceModuleName, eventPayload)
		}
	} else {
		logger.Infof("Drop event %s because server hasn't proper handlers", eventName)
	}
}
