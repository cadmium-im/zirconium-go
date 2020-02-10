package internal

import (
	"crypto/rand"
	"reflect"

	"github.com/ChronosX88/zirconium/internal/models"
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
	errMsg := models.NewBaseMessage(msg.MessageType, serverDomain, msg.From, false, StructToMap(protocolError))
	return errMsg
}
