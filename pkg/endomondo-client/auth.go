package endomondo

// AuthParams stores authorization parameters
type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	DeviceID string `json:"deviceId"`
	Country  string `json:"country"`
	Action   string `json:"action"`
}

// AuthResponse dto for authorization response data
type AuthResponse struct {
	ID int64 `json:"id"`
}
