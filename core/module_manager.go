package core

import (
	"reflect"
	"sync"
	"time"

	"github.com/ChronosX88/zirconium/core/models"
)

type C2SMessageHandler struct {
	HandlerFunc      func(origin *OriginC2S, message models.BaseMessage)
	AnonymousAllowed bool
}

type ModuleManager struct {
	moduleMutex        sync.Mutex
	c2sMessageHandlers map[string][]*C2SMessageHandler
	coreEventHandlers  map[string][]func(sourceModuleName string, event map[string]interface{})
}

type ModuleAPI interface {
	Hook(messageType string, anonymousAllowed bool, handlerFunc func(origin *OriginC2S, message models.BaseMessage))
	HookInternalEvent(eventName string, handlerFunc func(sourceModuleName string, event map[string]interface{}))
	Unhook(messageType string, handlerFunc func(origin *OriginC2S, message models.BaseMessage))
	UnhookInternalEvent(eventName string, handlerFunc func(sourceModuleName string, event map[string]interface{}))
	FireEvent(sourceModuleName string, eventName string, eventPayload map[string]interface{})
	GenerateToken(entityID, deviceID string, tokenExpireTimeDuration time.Duration) (string, error)
	GetServerDomain() string
}

func NewModuleManager() (ModuleAPI, error) {
	var mm = &ModuleManager{
		c2sMessageHandlers: make(map[string][]*C2SMessageHandler),
		coreEventHandlers:  make(map[string][]func(sourceModuleName string, event map[string]interface{})),
	}
	return mm, nil
}

func (mm *ModuleManager) Hook(messageType string, anonymousAllowed bool, handlerFunc func(origin *OriginC2S, message models.BaseMessage)) {
	mm.moduleMutex.Lock()
	mm.c2sMessageHandlers[messageType] = append(mm.c2sMessageHandlers[messageType], &C2SMessageHandler{
		HandlerFunc:      handlerFunc,
		AnonymousAllowed: anonymousAllowed,
	})
	mm.moduleMutex.Unlock()
}

func (mm *ModuleManager) HookInternalEvent(eventName string, handlerFunc func(sourceModuleName string, event map[string]interface{})) {
	mm.moduleMutex.Lock()
	mm.coreEventHandlers[eventName] = append(mm.coreEventHandlers[eventName], handlerFunc)
	mm.moduleMutex.Unlock()
}

func (mm *ModuleManager) Unhook(messageType string, handlerFunc func(origin *OriginC2S, message models.BaseMessage)) {
	mm.moduleMutex.Lock()
	defer mm.moduleMutex.Unlock()
	var handlers = mm.c2sMessageHandlers[messageType]
	if handlers != nil {
		for i, v := range handlers {
			if reflect.ValueOf(v.HandlerFunc).Pointer() == reflect.ValueOf(handlerFunc).Pointer() {
				handlers[i] = handlers[len(handlers)-1]
				handlers[len(handlers)-1] = nil
				handlers = handlers[:len(handlers)-1]
				mm.c2sMessageHandlers[messageType] = handlers
				break
			}
		}
	}
}

func (mm *ModuleManager) UnhookInternalEvent(eventName string, handlerFunc func(sourceModuleName string, event map[string]interface{})) {
	mm.moduleMutex.Lock()
	defer mm.moduleMutex.Unlock()
	var handlers = mm.coreEventHandlers[eventName]
	if handlers != nil {
		for i, v := range handlers {
			if reflect.ValueOf(v).Pointer() == reflect.ValueOf(handlerFunc).Pointer() {
				handlers[i] = handlers[len(handlers)-1]
				handlers[len(handlers)-1] = nil
				handlers = handlers[:len(handlers)-1]
				mm.coreEventHandlers[eventName] = handlers
				break
			}
		}
	}
}

func (mm *ModuleManager) FireEvent(sourceModuleName string, eventName string, eventPayload map[string]interface{}) {
	router.RoutecoreEvent(sourceModuleName, eventName, eventPayload)
}

func (mm *ModuleManager) GenerateToken(entityID, deviceID string, tokenExpireTimeDuration time.Duration) (string, error) {
	token, err := authManager.CreateNewToken(entityID, deviceID, tokenExpireTimeDuration)
	return token, err
}

func (mm *ModuleManager) GetServerDomain() string {
	return serverDomain
}
