package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"geo-entity-manager/backend/internal/apperror"
)

// HandleError maps domain and application errors to the uniform HTTP error response.
func HandleError(c *gin.Context, err error) {
	var maxBytesErr *http.MaxBytesError
	var valErr *apperror.ValidationError

	switch {
	case errors.Is(err, apperror.ErrNotFound):
		abortWithError(c, http.StatusNotFound, CodeNotFound, "Entitas tidak ditemukan", nil)
	case errors.As(err, &maxBytesErr):
		abortWithError(c, http.StatusRequestEntityTooLarge, CodePayloadTooLarge, "Ukuran payload terlalu besar (maksimal 1 MB)", nil)
	case errors.As(err, &valErr):
		var details []ErrorDetail
		for _, d := range valErr.Details {
			details = append(details, ErrorDetail{
				Field:   d.Field,
				Message: d.Message,
			})
		}
		// PRD spec says validation failure gets 422
		abortWithError(c, http.StatusUnprocessableEntity, CodeValidation, valErr.Message, details)
	default:
		// Also covers format json not valid, which the decode function wraps.
		// Wait, PRD says "400 JSON rusak". The decoder returns 400 for parsing errors?
		// Currently decoder returns ValidationError (which maps to 422).
		// Let's distinguish Syntax errors. If it's a syntax error, we want 400.
		
		if err.Error() == "json parse error: EOF" || err.Error() == "Body JSON kosong" || err.Error() == "Format JSON tidak valid" {
			abortWithError(c, http.StatusBadRequest, CodeBadRequest, "Format JSON tidak valid", nil)
		} else {
			// For unknown errors, we panic so the Recovery middleware can log it and return 500
			panic(err)
		}
	}
}
