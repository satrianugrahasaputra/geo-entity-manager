package validation

import (
	"testing"

	"geo-entity-manager/backend/internal/apperror"
	"geo-entity-manager/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testStruct is used to test the rules without relying entirely on the DTO.
type testStruct struct {
	Name        string             `validate:"required,notblank,min=3,max=100"`
	Type        model.EntityType   `validate:"required,oneof=vehicle iot_device facility other"`
	Status      model.EntityStatus `validate:"required,oneof=active inactive maintenance offline"`
	Description string             `validate:"max=500"`
	Latitude    float64            `validate:"finite,min=-90,max=90"`
	Longitude   float64            `validate:"finite,min=-180,max=180"`
}

func TestValidation(t *testing.T) {
	valid := testStruct{
		Name:      "Valid Name",
		Type:      model.TypeVehicle,
		Status:    model.StatusActive,
		Latitude:  -7.123,
		Longitude: 110.123,
	}

	tests := []struct {
		name        string
		mutator     func(*testStruct)
		wantErrKeys []string
	}{
		{"valid", func(s *testStruct) {}, nil},
		{"name empty", func(s *testStruct) { s.Name = "" }, []string{"name"}},
		{"name spaces only", func(s *testStruct) { s.Name = "   " }, []string{"name"}},
		{"name too short", func(s *testStruct) { s.Name = "ab" }, []string{"name"}},
		{"name 100 chars ok", func(s *testStruct) {
			b := make([]byte, 100)
			for i := range b {
				b[i] = 'a'
			}
			s.Name = string(b)
		}, nil},
		{"name too long", func(s *testStruct) {
			b := make([]byte, 101)
			for i := range b {
				b[i] = 'a'
			}
			s.Name = string(b)
		}, []string{"name"}},
		{"type invalid", func(s *testStruct) { s.Type = "unknown" }, []string{"type"}},
		{"status invalid", func(s *testStruct) { s.Status = "unknown" }, []string{"status"}},
		{"lat out of bounds", func(s *testStruct) { s.Latitude = 91 }, []string{"latitude"}},
		{"lng out of bounds", func(s *testStruct) { s.Longitude = -181 }, []string{"longitude"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := valid
			tt.mutator(&s)
			err := Struct(s)
			if len(tt.wantErrKeys) > 0 {
				require.Error(t, err)
				valErr, ok := err.(*apperror.ValidationError)
				require.True(t, ok, "expected ValidationError")

				var gotFields []string
				for _, d := range valErr.Details {
					gotFields = append(gotFields, d.Field)
				}
				assert.ElementsMatch(t, tt.wantErrKeys, gotFields)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
