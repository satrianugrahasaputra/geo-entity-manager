package config

import (
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func lookupFrom(env map[string]string) LookupFunc {
	return func(key string) (string, bool) {
		v, ok := env[key]
		return v, ok
	}
}

func TestLoadFrom(t *testing.T) {
	tests := []struct {
		name    string
		env     map[string]string
		want    Config
		wantErr []string
	}{
		{
			name: "defaults when env is empty",
			env:  map[string]string{},
			want: Config{
				HTTPPort:        8080,
				CORSOrigins:     []string{"http://localhost:5173"},
				LogLevel:        slog.LevelInfo,
				SeedData:        false,
				ShutdownTimeout: 10 * time.Second,
			},
		},
		{
			name: "all values overridden",
			env: map[string]string{
				"HTTP_PORT":        "9090",
				"DATABASE_URL":     " postgres://u:p@h:5432/d ",
				"CORS_ORIGINS":     " http://a.test , ,http://b.test ",
				"LOG_LEVEL":        "debug",
				"SEED_DATA":        "true",
				"SHUTDOWN_TIMEOUT": "3s",
			},
			want: Config{
				HTTPPort:        9090,
				DatabaseURL:     "postgres://u:p@h:5432/d",
				CORSOrigins:     []string{"http://a.test", "http://b.test"},
				LogLevel:        slog.LevelDebug,
				SeedData:        true,
				ShutdownTimeout: 3 * time.Second,
			},
		},
		{
			name: "empty CORS_ORIGINS disables CORS",
			env:  map[string]string{"CORS_ORIGINS": ""},
			want: Config{
				HTTPPort:        8080,
				CORSOrigins:     []string{},
				LogLevel:        slog.LevelInfo,
				ShutdownTimeout: 10 * time.Second,
			},
		},
		{
			name: "blank values fall back to defaults",
			env:  map[string]string{"HTTP_PORT": "  ", "LOG_LEVEL": ""},
			want: Config{
				HTTPPort:        8080,
				CORSOrigins:     []string{"http://localhost:5173"},
				LogLevel:        slog.LevelInfo,
				ShutdownTimeout: 10 * time.Second,
			},
		},
		{
			name: "invalid values are all reported",
			env: map[string]string{
				"HTTP_PORT":        "70000",
				"LOG_LEVEL":        "verbose",
				"SEED_DATA":        "maybe",
				"SHUTDOWN_TIMEOUT": "-1s",
			},
			wantErr: []string{"HTTP_PORT", "LOG_LEVEL", "SEED_DATA", "SHUTDOWN_TIMEOUT"},
		},
		{
			name:    "non numeric port",
			env:     map[string]string{"HTTP_PORT": "abc"},
			wantErr: []string{"HTTP_PORT"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadFrom(lookupFrom(tt.env))
			if len(tt.wantErr) > 0 {
				require.Error(t, err)
				for _, key := range tt.wantErr {
					assert.Contains(t, err.Error(), key)
				}
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
