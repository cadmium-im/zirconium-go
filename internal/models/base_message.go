package models

import "github.com/google/uuid"

// BaseMessage is a basic message model, basis of the whole protocol. It is used for a very easy protocol extension process.
type BaseMessage struct {
	ID          string                 `json:"id"`
	MessageType string                 `json:"type"`
	From        string                 `json:"from"`
	To          string                 `json:"to"`
	Ok          bool                   `json:"ok"`
	AuthToken   string                 `json:"authToken"`
	Payload     map[string]interface{} `json:"payload"`
}

func NewBaseMessage(messageType string, from string, to string, ok bool, payload map[string]interface{}) BaseMessage {
	uuid, _ := uuid.NewRandom()
	uuidStr := uuid.String()
	return BaseMessage{
		ID:          uuidStr,
		MessageType: messageType,
		From:        from,
		To:          to,
		Ok:          ok,
		Payload:     payload,
	}
}
