package config

import (
	"os"
	"testing"
)

func TestFileLoader(t *testing.T) {
	// 创建临时JSON文件
	tmpfile, err := os.CreateTemp("", "config*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	content := `{"name":"test","port":8080}`
	tmpfile.WriteString(content)
	tmpfile.Close()

	var config struct {
		Name string `json:"name"`
		Port int    `json:"port"`
	}

	loader := NewFileLoader(tmpfile.Name())
	err = loader.Load(&config)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if config.Name != "test" {
		t.Errorf("Expected 'test', got %s", config.Name)
	}
	if config.Port != 8080 {
		t.Errorf("Expected 8080, got %d", config.Port)
	}
}

func TestEnvLoader(t *testing.T) {
	os.Setenv("APP_NAME", "test")
	os.Setenv("APP_PORT", "8080")
	defer os.Unsetenv("APP_NAME")
	defer os.Unsetenv("APP_PORT")

	var config struct {
		Name string `env:"name"`
		Port int    `env:"port"`
	}

	loader := NewEnvLoader("APP")
	err := loader.Load(&config)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if config.Name != "test" {
		t.Errorf("Expected 'test', got %s", config.Name)
	}
	if config.Port != 8080 {
		t.Errorf("Expected 8080, got %d", config.Port)
	}
}

func TestMapLoader(t *testing.T) {
	data := map[string]interface{}{
		"name": "test",
		"port": 8080,
	}

	var config struct {
		Name string `map:"name"`
		Port int    `map:"port"`
	}

	loader := NewMapLoader(data)
	err := loader.Load(&config)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if config.Name != "test" {
		t.Errorf("Expected 'test', got %s", config.Name)
	}
	if config.Port != 8080 {
		t.Errorf("Expected 8080, got %d", config.Port)
	}
}

func TestMultiLoader(t *testing.T) {
	// 文件配置
	tmpfile, err := os.CreateTemp("", "config*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	tmpfile.WriteString(`{"name":"file","port":8080}`)
	tmpfile.Close()

	// 环境变量配置
	os.Setenv("APP_NAME", "env")
	defer os.Unsetenv("APP_NAME")

	var config struct {
		Name string `json:"name" env:"name"`
		Port int    `json:"port" env:"port"`
	}

	fileLoader := NewFileLoader(tmpfile.Name())
	envLoader := NewEnvLoader("APP")
	multiLoader := NewMultiLoader(fileLoader, envLoader)

	err = multiLoader.Load(&config)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// 环境变量应该覆盖文件配置
	if config.Name != "env" {
		t.Errorf("Expected 'env', got %s", config.Name)
	}
	// 端口应该保持文件配置的值（因为环境变量中没有设置）
	if config.Port != 8080 {
		t.Errorf("Expected 8080, got %d", config.Port)
	}
}
