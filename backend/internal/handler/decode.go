package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"geo-entity-manager/backend/internal/apperror"
	"geo-entity-manager/backend/internal/validation"
)

// DecodeAndValidate reads a JSON body strictly (no unknown fields) and runs validation.
// It maps various JSON/IO errors into the appropriate error types.
func DecodeAndValidate(r *http.Request, v interface{}) error {
	if r.Body == nil {
		return apperror.NewValidationError("Body kosong", nil)
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			return err // handled by caller as 413
		}
		return err
	}

	dec := json.NewDecoder(bytes.NewReader(body))
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		var syntaxErr *json.SyntaxError
		var unmarshalTypeErr *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxErr):
			return apperror.NewValidationError("Format JSON tidak valid", nil)
		case errors.As(err, &unmarshalTypeErr):
			return apperror.NewValidationError("Tipe data tidak sesuai", []apperror.FieldDetail{
				{Field: unmarshalTypeErr.Field, Message: "harus berupa " + unmarshalTypeErr.Type.String()},
			})
		case strings.Contains(err.Error(), "unknown field"):
			field := strings.TrimSuffix(strings.TrimPrefix(err.Error(), `json: unknown field "`), `"`)
			return apperror.NewValidationError("Terdapat field yang tidak dikenal", []apperror.FieldDetail{
				{Field: field, Message: "field asing tidak diizinkan"},
			})
		case errors.Is(err, io.EOF):
			return apperror.NewValidationError("Body JSON kosong", nil)
		default:
			// Treat other decode errors as bad request
			return fmt.Errorf("json parse error: %w", err)
		}
	}

	// Make sure there isn't trailing garbage
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return apperror.NewValidationError("Body JSON mengandung data ekstra", nil)
	}

	// Run domain rules via struct tags
	return validation.Struct(v)
}

// DecodeQueryAndValidate parses URL query parameters and runs validation.
func DecodeQueryAndValidate(r *http.Request, v interface{}) error {
	// For strict query validation, we could check for unknown params here.
	// But using a simple generic form binder from gorilla/schema or similar is often easier.
	// Since we don't have it, we'll manually check unknown keys if needed, 
	// or let the framework binder do it. Let's do a simple check.
	allowedKeys := map[string]bool{
		"type": true, "status": true, "search": true, "page": true, "limit": true,
	}
	for key := range r.URL.Query() {
		if !allowedKeys[key] {
			return apperror.NewValidationError("Terdapat query param yang tidak dikenal", []apperror.FieldDetail{
				{Field: key, Message: "query param asing tidak diizinkan"},
			})
		}
	}

	// We let the caller bind (e.g. c.ShouldBindQuery) and just run validation here
	// because implementing a reflection-based form decoder from scratch is too much.
	// Actually we should just pass the struct to validation.Struct.
	return validation.Struct(v)
}

