package core

import (
	"github.com/cadmium-im/zirconium-go/core/models"
	"github.com/cadmium-im/zirconium-go/core/models/auth"
	"github.com/cadmium-im/zirconium-go/core/utils"
	"github.com/fatih/structs"
	"github.com/google/logger"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	serverID    string
	authManager *AuthManager
}

func NewAuthHandler(am *AuthManager, serverID string) *AuthHandler {
	return &AuthHandler{
		authManager: am,
		serverID:    serverID,
	}
}

func (ah *AuthHandler) HandleMessage(s *Session, message models.BaseMessage) {
	var authRequest auth.AuthRequest
	err := mapstructure.Decode(message.Payload, &authRequest)
	if err != nil {
		logger.Errorf(err.Error())
		return
	}

	switch authRequest.Type {
	case "urn:cadmium:auth:simple":
		{
			token, claims, err := ah.authManager.HandleSimpleAuth(authRequest.Fields["username"].(string), authRequest.Fields["password"].(string))
			if err != nil {
				if mongo.ErrNoDocuments == err {
					msg := utils.PrepareErrorMessage(message, "auth-failed", "invalid username", ah.serverID)
					_ = s.Send(msg)
					return
				}
				logger.Errorf(err.Error())
				msg := utils.PrepareMessageInternalServerError(message, err, ah.serverID)
				_ = s.Send(msg)
				return
			}
			ar := auth.AuthResponse{
				Token:    token,
				DeviceID: claims.DeviceID,
			}
			payload := structs.Map(ar)
			msg := models.NewBaseMessage(message.ID, message.MessageType, ah.serverID, nil, true, payload)
			_ = s.Send(msg)
			s.Claims = claims
			return
		}
	case "urn:cadmium:auth:token":
		{
			claims, err := ah.authManager.ValidateToken(authRequest.Fields["token"].(string))
			if err != nil {
				logger.Errorf(err.Error())
				msg := utils.PrepareMessageInternalServerError(message, err, ah.serverID)
				_ = s.Send(msg)
				return
			}
			s.Claims = claims
			msg := models.NewBaseMessage(message.ID, message.MessageType, ah.serverID, nil, true, nil)
			_ = s.Send(msg)
			return
		}
	}
}

func (ah *AuthHandler) IsAuthorizationRequired() bool {
	return false
}

func (ah *AuthHandler) HandlingType() string {
	return "urn:cadmium:auth"
}
