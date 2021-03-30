package auth

type AuthRequest struct {
	Type   string                 `json:"type"`
	Fields map[string]interface{} `json:"fields"`
}

type AuthResponse struct {
	Token    string `json:"token"`
	DeviceID string `json:"deviceID"`
}
