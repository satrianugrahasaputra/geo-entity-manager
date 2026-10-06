package handler

import (
	"github.com/gin-gonic/gin"
)

// Error codes returned in the uniform error body.
const (
	CodeBadRequest       = "BAD_REQUEST"
	CodeValidation       = "VALIDATION_ERROR"
	CodeNotFound         = "NOT_FOUND"
	CodeMethodNotAllowed = "METHOD_NOT_ALLOWED"
	CodePayloadTooLarge  = "PAYLOAD_TOO_LARGE"
	CodeInternal         = "INTERNAL_ERROR"
)

// ErrorDetail describes a problem with a single field.
type ErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ErrorBody is the content of the "error" key.
type ErrorBody struct {
	Code    string        `json:"code"`
	Message string        `json:"message"`
	Details []ErrorDetail `json:"details"`
}

// ErrorResponse is the uniform error envelope: {"error":{"code","message","details"}}.
type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

// abortWithError writes the uniform error body and stops the handler chain.
// Details is always serialized as an array (never null) so clients can iterate safely.
func abortWithError(c *gin.Context, status int, code, message string, details []ErrorDetail) {
	if details == nil {
		details = []ErrorDetail{}
	}
	c.AbortWithStatusJSON(status, ErrorResponse{
		Error: ErrorBody{Code: code, Message: message, Details: details},
	})
}
