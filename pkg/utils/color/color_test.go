package color

import (
	"testing"
)

func TestColorize(t *testing.T) {
	Enable()
	result := Colorize("test", Red)
	if result == "" {
		t.Error("Expected non-empty result")
	}
}

func TestBlack(t *testing.T) {
	Enable()
	result := Black("test")
	if result == "" {
		t.Error("Expected non-empty result")
	}
}

func TestRed(t *testing.T) {
	Enable()
	result := Red("test")
	if result == "" {
		t.Error("Expected non-empty result")
	}
}

func TestGreen(t *testing.T) {
	Enable()
	result := Green("test")
	if result == "" {
		t.Error("Expected non-empty result")
	}
}

func TestBold(t *testing.T) {
	Enable()
	result := Bold("test")
	if result == "" {
		t.Error("Expected non-empty result")
	}
}

func TestSuccess(t *testing.T) {
	Enable()
	result := Success("test")
	if result == "" {
		t.Error("Expected non-empty result")
	}
}

func TestError(t *testing.T) {
	Enable()
	result := Error("test")
	if result == "" {
		t.Error("Expected non-empty result")
	}
}

func TestWarning(t *testing.T) {
	Enable()
	result := Warning("test")
	if result == "" {
		t.Error("Expected non-empty result")
	}
}

func TestInfo(t *testing.T) {
	Enable()
	result := Info("test")
	if result == "" {
		t.Error("Expected non-empty result")
	}
}

func TestSetEnabled(t *testing.T) {
	SetEnabled(true)
	if !IsEnabled() {
		t.Error("Expected enabled")
	}
	SetEnabled(false)
	if IsEnabled() {
		t.Error("Expected disabled")
	}
}

func TestDisable(t *testing.T) {
	Disable()
	result := Red("test")
	if result != "test" {
		t.Error("Expected plain text when disabled")
	}
}
