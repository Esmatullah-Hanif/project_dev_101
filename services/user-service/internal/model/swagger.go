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

type UserSuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    User   `json:"data"`
}

type UserListSuccessResponse struct {
	Success bool                `json:"success"`
	Message string              `json:"message"`
	Data    []User              `json:"data"`
	Meta    map[string]interface{} `json:"meta"`
}
