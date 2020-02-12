package core

import (
	"crypto/rand"
	"log"
	"reflect"

	"github.com/ChronosX88/zirconium/core/models"
	"github.com/google/logger"
	"github.com/hashicorp/go-plugin"
)

func GenRandomBytes(size int) (blk []byte, err error) {
	blk = make([]byte, size)
	_, err = rand.Read(blk)
	return
}

func StructToMap(item interface{}) map[string]interface{} {
	res := map[string]interface{}{}
	if item == nil {
		return res
	}
	v := reflect.TypeOf(item)
	reflectValue := reflect.ValueOf(item)
	reflectValue = reflect.Indirect(reflectValue)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		tag := v.Field(i).Tag.Get("json")
		field := reflectValue.Field(i).Interface()
		if tag != "" && tag != "-" {
			if v.Field(i).Type.Kind() == reflect.Struct {
				res[tag] = StructToMap(field)
			} else {
				res[tag] = field
			}
		}
	}
	return res
}

func PrepareMessageUnauthorized(msg models.BaseMessage) models.BaseMessage {
	protocolError := models.ProtocolError{
		ErrCode:    "unauthorized",
		ErrText:    "Unauthorized access",
		ErrPayload: make(map[string]interface{}),
	}
	errMsg := models.NewBaseMessage(msg.ID, msg.MessageType, serverDomain, msg.From, false, StructToMap(protocolError))
	return errMsg
}

var HandshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "TEST_TEST",
	MagicCookieValue: "qwerty",
}

type PluginOpts struct {
	Module      ModuleFunc
	RunAsPlugin bool
}

// GetPluginMap returns the plugin map defined Hashicorp go-plugin.
// The reserved parameter should only be used by the RPC receiver (the plugin).
// Otherwise, reserved should be nil for the RPC sender (the mainapp).
func GetPluginMap(reserved *PluginOpts) map[string]plugin.Plugin {
	var moduleObj ModuleRef

	if reserved != nil {
		moduleObj.F = reserved.Module
	}

	return map[string]plugin.Plugin{
		ModuleInterfaceName: &moduleObj,
	}
}

func StartPlugin(opts *PluginOpts, quit chan bool) {
	if opts.RunAsPlugin {
		go func() {
			logger.Info("Starting plugin communication...")

			plugin.Serve(&plugin.ServeConfig{
				HandshakeConfig: HandshakeConfig,
				Plugins:         GetPluginMap(opts),
			})

			logger.Info("Exiting plugin communication...")

			quit <- true
			logger.Info("Exiting plugin...")
		}()
	} else {
		log.Println("Starting in standalone mode...")
	}
}
