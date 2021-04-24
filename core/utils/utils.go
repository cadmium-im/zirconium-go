package utils

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/alexedwards/argon2id"
	"github.com/cadmium-im/zirconium-go/core/models"
	"github.com/fatih/structs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrPasswordIsEmpty    = fmt.Errorf("the password is empty")
	ErrPasswordIsTooShort = fmt.Errorf("the password is too short")
)

func GenRandomBytes(size int) (blk []byte, err error) {
	blk = make([]byte, size)
	_, err = rand.Read(blk)
	return
}

func PrepareMessageUnauthorized(msg models.BaseMessage, serverDomain string) models.BaseMessage {
	protocolError := models.ProtocolError{
		ErrCode:    "unauthorized",
		ErrText:    "Unauthorized access",
		ErrPayload: make(map[string]interface{}),
	}
	errMsg := models.NewBaseMessage(msg.ID, msg.MessageType, serverDomain, msg.From, false, structs.Map(protocolError))
	return errMsg
}

func PrepareMessageInternalServerError(msg models.BaseMessage, err error, serverID string) models.BaseMessage {
	protocolError := models.ProtocolError{
		ErrCode:    "internal-server-error",
		ErrText:    err.Error(),
		ErrPayload: nil,
	}
	errMsg := models.NewBaseMessage(msg.ID, msg.MessageType, serverID, msg.From, false, structs.Map(protocolError))
	return errMsg
}

func PrepareErrorMessage(msg models.BaseMessage, errorType string, errorText string, serverID string) models.BaseMessage {
	protocolError := models.ProtocolError{
		ErrCode:    errorType,
		ErrText:    errorText,
		ErrPayload: nil,
	}
	errMsg := models.NewBaseMessage(msg.ID, msg.MessageType, serverID, msg.From, false, structs.Map(protocolError))
	return errMsg
}

func HashPassword(password string) (string, error) {
	return argon2id.CreateHash(password, argon2id.DefaultParams)
}

func ValidatePassword(password string) error {
	if password == "" {
		return ErrPasswordIsEmpty
	}
	if len(password) < 4 {
		return ErrPasswordIsTooShort
	}
	return nil
}

func IsCollectionExists(ctx context.Context, db *mongo.Database, collectionName string) (bool, error) {
	isExists := false
	names, err := db.ListCollectionNames(ctx, bson.D{{"name", collectionName}})
	if err != nil {
		return false, err
	}

	for _, name := range names {
		if name == collectionName {
			isExists = true
			break
		}
	}
	return isExists, nil
}

func InStringArray(val string, array []string) (ok bool) {
	for i := range array {
		if ok = array[i] == val; ok {
			return
		}
	}
	return
}
