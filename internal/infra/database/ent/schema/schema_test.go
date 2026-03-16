// Package schema provides tests for ent schema definitions.
package schema

import (
	"testing"

	"entgo.io/ent"
	"github.com/stretchr/testify/assert"
)

func TestSchemasSlice(t *testing.T) {
	// Test that Schemas is defined and not empty
	assert.NotNil(t, Schemas)
	assert.GreaterOrEqual(t, len(Schemas), 1)
}

func TestSchemasContainsUser(t *testing.T) {
	// Test that Schemas contains User schema
	found := false
	for _, s := range Schemas {
		if _, ok := s.(*User); ok {
			found = true
			break
		}
	}
	assert.True(t, found, "Schemas should contain User schema")
}

func TestSchemasType(t *testing.T) {
	// Test that all schemas implement ent.Interface
	for _, s := range Schemas {
		assert.Implements(t, (*ent.Interface)(nil), s)
	}
}

func TestSchemasElements(t *testing.T) {
	// Test that we can access individual schemas
	assert.GreaterOrEqual(t, len(Schemas), 1)

	// First schema should be User
	userSchema, ok := Schemas[0].(*User)
	assert.True(t, ok, "First schema should be User")
	assert.NotNil(t, userSchema)
}

func TestUserSchemaImplementsInterface(t *testing.T) {
	// Test that User struct implements ent.Interface
	user := &User{}
	assert.Implements(t, (*ent.Interface)(nil), user)
}

func TestSchemaCount(t *testing.T) {
	// Test that we have the expected number of schemas
	// This test ensures we update tests when adding new schemas
	expectedSchemaCount := 1 // Currently only User schema
	assert.Equal(t, expectedSchemaCount, len(Schemas))
}
