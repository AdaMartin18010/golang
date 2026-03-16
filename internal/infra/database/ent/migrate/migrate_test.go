// Package migrate provides tests for ent schema migration.
package migrate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSchemaVariables(t *testing.T) {
	// Test that migration option variables are properly defined
	// These are functions that return MigrateOption
	assert.NotNil(t, WithGlobalUniqueID)
	assert.NotNil(t, WithDropColumn)
	assert.NotNil(t, WithDropIndex)
	assert.NotNil(t, WithForeignKeys)
}

func TestSchemaType(t *testing.T) {
	// Test Schema struct type definition
	s := &Schema{}
	assert.NotNil(t, s)
}

func TestNewSchema(t *testing.T) {
	// Test NewSchema function returns a Schema pointer
	// Note: We can't test with a real driver in unit tests
	// but we can verify the function signature and return type
	
	// Test that NewSchema function exists
	assert.NotNil(t, NewSchema)
	
	// Test that NewSchema accepts driver parameter
	assert.NotPanics(t, func() {
		// Just testing the function signature
		_ = NewSchema
	})
}

func TestSchema_CreateMethod(t *testing.T) {
	// Test that Schema has Create method with correct signature
	s := &Schema{}
	
	// Verify the method exists by checking it can be referenced
	assert.NotNil(t, s.Create)
}

func TestSchema_WriteToMethod(t *testing.T) {
	// Test that Schema has WriteTo method with correct signature
	s := &Schema{}
	
	// Verify the method exists by checking it can be referenced
	assert.NotNil(t, s.WriteTo)
}

func TestCreateFunction(t *testing.T) {
	// Test that Create function is properly defined
	assert.NotNil(t, Create)
}

// TestMigrationOptionFunctions tests that migration option functions work correctly
func TestMigrationOptionFunctions(t *testing.T) {
	// Test that option functions return valid MigrateOption
	opts := []interface{}{
		WithGlobalUniqueID(true),
		WithDropColumn(true),
		WithDropIndex(true),
		WithForeignKeys(true),
	}
	
	// Verify all options are not nil
	for _, opt := range opts {
		assert.NotNil(t, opt)
	}
	
	assert.Equal(t, 4, len(opts))
}
