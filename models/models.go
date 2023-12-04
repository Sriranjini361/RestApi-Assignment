package models

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SuccessResponseSignup struct {
	Message string `json:"message"`
}

type ErrorResponseSignup struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	SID int `json:"sid"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
