package models

type ProtocolError struct {
	ErrCode    string                 `json:"errCode"`
	ErrText    string                 `json:"errText"`
	ErrPayload map[string]interface{} `json:"errPayload"`
}
