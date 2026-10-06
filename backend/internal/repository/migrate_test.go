package repository

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMigrateAndConnect(t *testing.T) {
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		t.Skip("TEST_DATABASE_URL is not set, skipping database integration tests")
	}

	// 1. Run migrations up
	err := Migrate(dbURL)
	require.NoError(t, err, "Migration should succeed")

	// 2. Connect to the migrated DB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := ConnectDB(ctx, dbURL)
	require.NoError(t, err, "Should connect to the database")
	defer pool.Close()

	// 3. Test that the constraints exist and work
	// A latitude > 90 should be rejected by the chk_latitude constraint
	_, err = pool.Exec(ctx, `
		INSERT INTO entities (name, type, status, latitude, longitude)
		VALUES ('Test', 'vehicle', 'active', 91.0, 100.0)
	`)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "chk_latitude", "Constraint should prevent invalid latitude")

	// A wrong type should be rejected
	_, err = pool.Exec(ctx, `
		INSERT INTO entities (name, type, status, latitude, longitude)
		VALUES ('Test', 'spaceship', 'active', -7.0, 110.0)
	`)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "chk_type", "Constraint should prevent invalid type")
}
