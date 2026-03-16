// Package schema provides tests for the User ent schema.
package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"entgo.io/ent"
)

func TestUserStruct(t *testing.T) {
	// Test that User struct is properly defined
	user := &User{}
	assert.NotNil(t, user)
}

func TestUserFields(t *testing.T) {
	// Test Fields method returns expected fields
	user := &User{}
	fields := user.Fields()
	
	assert.NotNil(t, fields)
	assert.Equal(t, 5, len(fields))
	
	// Get field map for easier testing
	fieldMap := make(map[string]ent.Field)
	for _, f := range fields {
		// Use reflection or type assertion to get field name
		// Since ent.Field is an interface, we test by checking the descriptor
		desc := f.Descriptor()
		fieldMap[desc.Name] = f
	}
	
	// Test that all expected fields exist
	expectedFields := []string{"id", "email", "name", "created_at", "updated_at"}
	for _, name := range expectedFields {
		assert.Contains(t, fieldMap, name, "Field %s should exist", name)
	}
}

func TestUserFieldsDescriptor(t *testing.T) {
	user := &User{}
	fields := user.Fields()
	
	// Test id field
	idField := findFieldByName(fields, "id")
	assert.NotNil(t, idField)
	idDesc := idField.Descriptor()
	assert.Equal(t, "id", idDesc.Name)
	assert.NotNil(t, idDesc.Info)
	
	// Test email field
	emailField := findFieldByName(fields, "email")
	assert.NotNil(t, emailField)
	emailDesc := emailField.Descriptor()
	assert.Equal(t, "email", emailDesc.Name)
	assert.True(t, emailDesc.Unique, "email should be unique")
	assert.NotEmpty(t, emailDesc.Validators, "email should have validators")
	
	// Test name field
	nameField := findFieldByName(fields, "name")
	assert.NotNil(t, nameField)
	nameDesc := nameField.Descriptor()
	assert.Equal(t, "name", nameDesc.Name)
	assert.NotEmpty(t, nameDesc.Validators, "name should have validators")
	
	// Test created_at field
	createdAtField := findFieldByName(fields, "created_at")
	assert.NotNil(t, createdAtField)
	createdAtDesc := createdAtField.Descriptor()
	assert.Equal(t, "created_at", createdAtDesc.Name)
	assert.True(t, createdAtDesc.Immutable, "created_at should be immutable")
	assert.NotNil(t, createdAtDesc.Default, "created_at should have default value")
	
	// Test updated_at field
	updatedAtField := findFieldByName(fields, "updated_at")
	assert.NotNil(t, updatedAtField)
	updatedAtDesc := updatedAtField.Descriptor()
	assert.Equal(t, "updated_at", updatedAtDesc.Name)
	assert.NotNil(t, updatedAtDesc.Default, "updated_at should have default value")
	assert.NotNil(t, updatedAtDesc.UpdateDefault, "updated_at should have update default value")
}

// Helper function to find field by name
func findFieldByName(fields []ent.Field, name string) ent.Field {
	for _, f := range fields {
		if f.Descriptor().Name == name {
			return f
		}
	}
	return nil
}

func TestUserEdges(t *testing.T) {
	// Test Edges method
	user := &User{}
	edges := user.Edges()
	
	// User currently has no edges defined
	assert.Nil(t, edges)
}

func TestEmailRegex(t *testing.T) {
	// Test that email regex is defined and works correctly
	assert.NotNil(t, emailRegex)
	
	// Test valid emails
	validEmails := []string{
		"user@example.com",
		"test.email@domain.org",
		"user123@test.co.uk",
		"first.last@company.io",
		"user+tag@example.com",
	}
	
	for _, email := range validEmails {
		assert.True(t, emailRegex.MatchString(email), "Email %s should be valid", email)
	}
	
	// Test invalid emails
	invalidEmails := []string{
		"invalid-email",
		"@example.com",
		"user@",
		"user@.com",
		"",
	}
	
	for _, email := range invalidEmails {
		assert.False(t, emailRegex.MatchString(email), "Email %s should be invalid", email)
	}
}

func TestEmailRegexPattern(t *testing.T) {
	// Verify the regex pattern is correct
	expectedPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	assert.Equal(t, expectedPattern, emailRegex.String())
}

func TestDefaultValues(t *testing.T) {
	user := &User{}
	fields := user.Fields()
	
	// Test created_at has time.Now as default
	createdAtField := findFieldByName(fields, "created_at")
	assert.NotNil(t, createdAtField)
	createdAtDesc := createdAtField.Descriptor()
	assert.NotNil(t, createdAtDesc.Default)
	
	// Test updated_at has time.Now as default and update default
	updatedAtField := findFieldByName(fields, "updated_at")
	assert.NotNil(t, updatedAtField)
	updatedAtDesc := updatedAtField.Descriptor()
	assert.NotNil(t, updatedAtDesc.Default)
	assert.NotNil(t, updatedAtDesc.UpdateDefault)
}
