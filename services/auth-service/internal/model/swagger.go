package model

type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

type AuthSuccessResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Data    AuthResponse  `json:"data"`
}
