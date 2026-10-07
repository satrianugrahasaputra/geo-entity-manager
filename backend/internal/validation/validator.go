package validation

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"geo-entity-manager/backend/internal/apperror"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())

	// Register custom validations
	_ = validate.RegisterValidation("notblank", func(fl validator.FieldLevel) bool {
		return strings.TrimSpace(fl.Field().String()) != ""
	})

	_ = validate.RegisterValidation("finite", func(fl validator.FieldLevel) bool {
		val := fl.Field().Float()
		return !math.IsNaN(val) && !math.IsInf(val, 0)
	})
}

// Struct validates a struct and translates errors to apperror.ValidationError.
func Struct(s interface{}) error {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	var valErrs validator.ValidationErrors
	if !errors.As(err, &valErrs) {
		return err // Should not happen with validator
	}

	var details []apperror.FieldDetail
	for _, e := range valErrs {
		details = append(details, apperror.FieldDetail{
			Field:   strings.ToLower(e.Field()),
			Message: translateError(e),
		})
	}

	return apperror.NewValidationError("Input tidak valid", details)
}

func translateError(e validator.FieldError) string {
	switch e.Tag() {
	case "required", "notblank":
		return "tidak boleh kosong"
	case "min":
		if e.Type().Kind().String() == "float64" {
			return fmt.Sprintf("harus lebih besar atau sama dengan %s", e.Param())
		}
		return fmt.Sprintf("minimal %s karakter", e.Param())
	case "max":
		if e.Type().Kind().String() == "float64" {
			return fmt.Sprintf("harus lebih kecil atau sama dengan %s", e.Param())
		}
		return fmt.Sprintf("maksimal %s karakter", e.Param())
	case "oneof":
		return "harus salah satu dari: " + strings.ReplaceAll(e.Param(), " ", ", ")
	case "finite":
		return "angka tidak valid"
	default:
		return "nilai tidak valid"
	}
}
