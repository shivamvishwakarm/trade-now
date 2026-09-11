package user

// ErrorResponse is the standard error envelope for user endpoints.
type ErrorResponse struct {
	Error string `json:"error" example:"something went wrong"`
}
