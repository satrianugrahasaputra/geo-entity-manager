package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"geo-entity-manager/backend/internal/apperror"
	"geo-entity-manager/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeAndValidate(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		wantError   bool
		wantErrMsg  string
		wantDetails []string
	}{
		{
			name:      "valid request",
			body:      `{"name":"test","type":"vehicle","status":"active","latitude":1,"longitude":1}`,
			wantError: false,
		},
		{
			name:       "empty body",
			body:       ``,
			wantError:  true,
			wantErrMsg: "Body JSON kosong",
		},
		{
			name:       "bad json",
			body:       `{name:"test"`,
			wantError:  true,
			wantErrMsg: "Format JSON tidak valid",
		},
		{
			name:        "unknown field",
			body:        `{"name":"test","type":"vehicle","status":"active","latitude":1,"longitude":1,"extra":true}`,
			wantError:   true,
			wantErrMsg:  "Terdapat field yang tidak dikenal",
			wantDetails: []string{"extra"},
		},
		{
			name:        "wrong type",
			body:        `{"name":"test","type":"vehicle","status":"active","latitude":"abc","longitude":1}`,
			wantError:   true,
			wantErrMsg:  "Tipe data tidak sesuai",
			wantDetails: []string{"latitude"},
		},
		{
			name:       "trailing garbage",
			body:       `{"name":"test","type":"vehicle","status":"active","latitude":1,"longitude":1} extra`,
			wantError:  true,
			wantErrMsg: "Body JSON mengandung data ekstra",
		},
		{
			name:        "validation error",
			body:        `{"name":"","type":"vehicle","status":"active","latitude":1,"longitude":1}`,
			wantError:   true,
			wantErrMsg:  "Input tidak valid",
			wantDetails: []string{"name"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(tt.body))

			var v CreateEntityRequest
			err := DecodeAndValidate(req, &v)

			if tt.wantError {
				require.Error(t, err)
				valErr, ok := err.(*apperror.ValidationError)
				require.True(t, ok, "expected ValidationError")
				assert.Equal(t, tt.wantErrMsg, valErr.Message)

				var gotFields []string
				for _, d := range valErr.Details {
					gotFields = append(gotFields, d.Field)
				}
				assert.ElementsMatch(t, tt.wantDetails, gotFields)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "test", v.Name)
				assert.Equal(t, model.TypeVehicle, v.Type)
			}
		})
	}
}
