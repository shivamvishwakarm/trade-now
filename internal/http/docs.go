// Package http contains the HTTP layer for the Trade Now API.
// This file defines shared response models used across Swagger annotations.
package http

// ErrorResponse is the standard error envelope returned on failure.
//
//	@Description	Standard error response body.
type ErrorResponse struct {
	Error string `json:"error" example:"something went wrong"`
}

// MessageResponse is the standard success envelope for operations that return no data.
//
//	@Description	Generic success message response.
type MessageResponse struct {
	Message string `json:"message" example:"operation successful"`
}
