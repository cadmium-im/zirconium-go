package models

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

func NewBaseMessage(id, messageType, from, to string, ok bool, payload map[string]interface{}) BaseMessage {
	return BaseMessage{
		ID:          id,
		MessageType: messageType,
		From:        from,
		To:          to,
		Ok:          ok,
		Payload:     payload,
	}
}
