package models

type ProtocolError struct {
	ErrCode    string                 `json:"code"`
	ErrText    string                 `json:"text"`
	ErrPayload map[string]interface{} `json:"payload"`
}
