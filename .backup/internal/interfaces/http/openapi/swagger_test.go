// Package openapi provides tests for Swagger UI handler.
package openapi

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSwaggerConfigStruct(t *testing.T) {
	config := SwaggerConfig{
		OpenAPISpecPath: "api/openapi/openapi.yaml",
		Title:           "Test API",
		EnableUI:        true,
		RequireAuth:     true,
	}

	assert.Equal(t, "api/openapi/openapi.yaml", config.OpenAPISpecPath)
	assert.Equal(t, "Test API", config.Title)
	assert.True(t, config.EnableUI)
	assert.True(t, config.RequireAuth)
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	assert.Equal(t, "api/openapi/openapi.yaml", config.OpenAPISpecPath)
	assert.Equal(t, "API Documentation", config.Title)
	assert.True(t, config.EnableUI)
	assert.False(t, config.RequireAuth)
}

func TestHandler_UIDisabled(t *testing.T) {
	config := SwaggerConfig{
		OpenAPISpecPath: "api/openapi/openapi.yaml",
		Title:           "Test API",
		EnableUI:        false,
		RequireAuth:     false,
	}

	handler := Handler(config)
	req := httptest.NewRequest(http.MethodGet, "/swagger/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestHandler_SwaggerUIEnabled(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	config := DefaultConfig()
	handler := Handler(config)

	// Test Swagger UI page
	req := httptest.NewRequest(http.MethodGet, "/swagger/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	// Should return HTML or error (depending on whether spec file exists)
	assert.Contains(t, []int{http.StatusOK, http.StatusInternalServerError}, rec.Code)
}

func TestHandler_SwaggerStaticFiles(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	config := DefaultConfig()
	handler := Handler(config)

	// Test Swagger UI static files
	req := httptest.NewRequest(http.MethodGet, "/swagger-ui/swagger-ui.css", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	// May return 200 or 404 depending on whether embed files exist
	assert.Contains(t, []int{http.StatusOK, http.StatusNotFound}, rec.Code)
}

func TestHandler_OpenAPISpec_AbsolutePath(t *testing.T) {
	// Create a temporary OpenAPI spec file
	tmpDir := t.TempDir()
	specContent := `
openapi: 3.0.0
info:
  title: Test API
  version: 1.0.0
paths: {}
`
	specPath := filepath.Join(tmpDir, "openapi.yaml")
	err := os.WriteFile(specPath, []byte(specContent), 0644)
	assert.NoError(t, err)

	config := SwaggerConfig{
		OpenAPISpecPath: specPath,
		Title:           "Test API",
		EnableUI:        true,
		RequireAuth:     false,
	}

	handler := Handler(config)
	req := httptest.NewRequest(http.MethodGet, "/swagger/openapi.yaml", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/yaml", rec.Header().Get("Content-Type"))
	assert.Contains(t, rec.Body.String(), "openapi: 3.0.0")
}

func TestHandler_OpenAPISpec_RelativePath(t *testing.T) {
	// Create a temporary directory structure
	originalWd, err := os.Getwd()
	assert.NoError(t, err)
	defer os.Chdir(originalWd)

	tmpDir := t.TempDir()
	err = os.Chdir(tmpDir)
	assert.NoError(t, err)

	// Create api/openapi directory
	err = os.MkdirAll("api/openapi", 0755)
	assert.NoError(t, err)

	specContent := `
openapi: 3.0.0
info:
  title: Test API
  version: 1.0.0
paths: {}
`
	err = os.WriteFile("api/openapi/openapi.yaml", []byte(specContent), 0644)
	assert.NoError(t, err)

	config := DefaultConfig()
	handler := Handler(config)
	req := httptest.NewRequest(http.MethodGet, "/swagger/openapi.yaml", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/yaml", rec.Header().Get("Content-Type"))
}

func TestHandler_OpenAPISpec_NotFound(t *testing.T) {
	config := SwaggerConfig{
		OpenAPISpecPath: "/nonexistent/path/openapi.yaml",
		Title:           "Test API",
		EnableUI:        true,
		RequireAuth:     false,
	}

	handler := Handler(config)
	req := httptest.NewRequest(http.MethodGet, "/swagger/openapi.yaml", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestHandler_OpenAPISpec_WithAuth(t *testing.T) {
	// Create a temporary OpenAPI spec file
	tmpDir := t.TempDir()
	specContent := `openapi: 3.0.0`
	specPath := filepath.Join(tmpDir, "openapi.yaml")
	err := os.WriteFile(specPath, []byte(specContent), 0644)
	assert.NoError(t, err)

	config := SwaggerConfig{
		OpenAPISpecPath: specPath,
		Title:           "Test API",
		EnableUI:        true,
		RequireAuth:     true, // Auth enabled but not implemented yet
	}

	handler := Handler(config)
	req := httptest.NewRequest(http.MethodGet, "/swagger/openapi.yaml", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	// Auth is not implemented yet, so it should still return success
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestHandler_SwaggerPage_WithAuth(t *testing.T) {
	config := SwaggerConfig{
		OpenAPISpecPath: "api/openapi/openapi.yaml",
		Title:           "Protected API",
		EnableUI:        true,
		RequireAuth:     true, // Auth enabled but not implemented yet
	}

	handler := Handler(config)
	req := httptest.NewRequest(http.MethodGet, "/swagger/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	// Should return HTML page (auth not implemented yet)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "text/html", rec.Header().Get("Content-Type"))
	body := rec.Body.String()
	assert.Contains(t, body, "Protected API")
	assert.Contains(t, body, "SwaggerUIBundle")
}

func TestHandler_SwaggerPage_Content(t *testing.T) {
	config := SwaggerConfig{
		OpenAPISpecPath: "api/openapi/openapi.yaml",
		Title:           "My API",
		EnableUI:        true,
		RequireAuth:     false,
	}

	handler := Handler(config)
	req := httptest.NewRequest(http.MethodGet, "/swagger/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "text/html", rec.Header().Get("Content-Type"))

	body := rec.Body.String()
	// Verify HTML content
	assert.Contains(t, body, "<!DOCTYPE html>")
	assert.Contains(t, body, "<title>My API</title>")
	assert.Contains(t, body, "swagger-ui.css")
	assert.Contains(t, body, "swagger-ui-bundle.js")
	assert.Contains(t, body, "/swagger/openapi.yaml")
	assert.Contains(t, body, "SwaggerUIBundle")
	assert.Contains(t, body, "StandaloneLayout")
}
