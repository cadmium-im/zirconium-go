package models

type ProtocolError struct {
	ErrCode    string                 `structs:"code"`
	ErrText    string                 `structs:"text"`
	ErrPayload map[string]interface{} `structs:"payload,omitempty"`
}
