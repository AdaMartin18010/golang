package process

import (
	"os"
	"testing"
)

func TestGetPID(t *testing.T) {
	pid := GetPID()
	if pid <= 0 {
		t.Error("Expected positive PID")
	}
}

func TestGetPPID(t *testing.T) {
	ppid := GetPPID()
	if ppid <= 0 {
		t.Error("Expected positive PPID")
	}
}

func TestGetExecutable(t *testing.T) {
	executable, err := GetExecutable()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if executable == "" {
		t.Error("Expected non-empty executable path")
	}
}

func TestGetArgs(t *testing.T) {
	args := GetArgs()
	if len(args) == 0 {
		t.Error("Expected at least one argument")
	}
}

func TestGetEnv(t *testing.T) {
	os.Setenv("TEST_VAR", "test_value")
	value := GetEnv("TEST_VAR")
	if value != "test_value" {
		t.Errorf("Expected 'test_value', got %s", value)
	}
}

func TestSetEnv(t *testing.T) {
	err := SetEnv("TEST_VAR", "test_value")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	value := GetEnv("TEST_VAR")
	if value != "test_value" {
		t.Errorf("Expected 'test_value', got %s", value)
	}
}

func TestGetWorkingDir(t *testing.T) {
	dir, err := GetWorkingDir()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if dir == "" {
		t.Error("Expected non-empty directory")
	}
}

func TestGetProcessInfo(t *testing.T) {
	info, err := GetProcessInfo()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if info.PID <= 0 {
		t.Error("Expected positive PID")
	}
}

func TestRunCommand(t *testing.T) {
	output, err := RunCommand("echo", "test")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if output == "" {
		t.Error("Expected non-empty output")
	}
}

func TestIsRoot(t *testing.T) {
	_ = IsRoot() // 不应该panic
}

func TestGetHostname(t *testing.T) {
	hostname, err := GetHostname()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if hostname == "" {
		t.Error("Expected non-empty hostname")
	}
}

func TestGetTempDir(t *testing.T) {
	dir := GetTempDir()
	if dir == "" {
		t.Error("Expected non-empty temp directory")
	}
}
