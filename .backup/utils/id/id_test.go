package id

import (
	"testing"
)

func TestUUID(t *testing.T) {
	id := UUID()
	if len(id) == 0 {
		t.Error("Expected non-empty UUID")
	}
	if len(id) != 36 {
		t.Errorf("Expected UUID length 36, got %d", len(id))
	}
}

func TestShortUUID(t *testing.T) {
	id := ShortUUID()
	if len(id) == 0 {
		t.Error("Expected non-empty short UUID")
	}
	if len(id) != 22 {
		t.Errorf("Expected short UUID length 22, got %d", len(id))
	}
}

func TestNanoID(t *testing.T) {
	id := NanoID()
	if len(id) == 0 {
		t.Error("Expected non-empty NanoID")
	}
	if len(id) != 21 {
		t.Errorf("Expected NanoID length 21, got %d", len(id))
	}
}

func TestNanoIDWithSize(t *testing.T) {
	id := NanoIDWithSize(10)
	if len(id) != 10 {
		t.Errorf("Expected NanoID length 10, got %d", len(id))
	}
}

func TestTimestampID(t *testing.T) {
	id := TimestampID()
	if len(id) == 0 {
		t.Error("Expected non-empty timestamp ID")
	}
}

func TestRandomHex(t *testing.T) {
	id := RandomHex()
	if len(id) != 32 {
		t.Errorf("Expected random hex length 32, got %d", len(id))
	}
}

func TestRandomHexWithLength(t *testing.T) {
	id := RandomHexWithLength(16)
	if len(id) != 16 {
		t.Errorf("Expected random hex length 16, got %d", len(id))
	}
}

func TestRandomBase64(t *testing.T) {
	id := RandomBase64()
	if len(id) == 0 {
		t.Error("Expected non-empty random base64 ID")
	}
}

func TestSequentialID(t *testing.T) {
	id1 := SequentialID("TEST")
	id2 := SequentialID("TEST")
	if id1 == id2 {
		t.Error("Expected different sequential IDs")
	}
}

func TestSnowflakeGenerator(t *testing.T) {
	generator := NewSnowflakeGenerator(1, 1)
	id1 := generator.Generate()
	id2 := generator.Generate()

	if id1 == id2 {
		t.Error("Expected different snowflake IDs")
	}
	if len(id1) == 0 {
		t.Error("Expected non-empty snowflake ID")
	}
}

func TestUUIDGenerator(t *testing.T) {
	generator := NewUUIDGenerator()
	id := generator.Generate()
	if len(id) != 36 {
		t.Errorf("Expected UUID length 36, got %d", len(id))
	}
}

func TestShortUUIDGenerator(t *testing.T) {
	generator := NewShortUUIDGenerator()
	id := generator.Generate()
	if len(id) != 22 {
		t.Errorf("Expected short UUID length 22, got %d", len(id))
	}
}

func TestNanoIDGenerator(t *testing.T) {
	generator := NewNanoIDGenerator(21)
	id := generator.Generate()
	if len(id) != 21 {
		t.Errorf("Expected NanoID length 21, got %d", len(id))
	}
}

func TestTimestampIDGenerator(t *testing.T) {
	generator := NewTimestampIDGenerator()
	id := generator.Generate()
	if len(id) == 0 {
		t.Error("Expected non-empty timestamp ID")
	}
}

func TestRandomHexGenerator(t *testing.T) {
	generator := NewRandomHexGenerator(32)
	id := generator.Generate()
	if len(id) != 32 {
		t.Errorf("Expected random hex length 32, got %d", len(id))
	}
}

func TestSequentialIDGenerator(t *testing.T) {
	generator := NewSequentialIDGenerator("TEST")
	id1 := generator.Generate()
	id2 := generator.Generate()

	if id1 == id2 {
		t.Error("Expected different sequential IDs")
	}
	if len(id1) == 0 {
		t.Error("Expected non-empty sequential ID")
	}
}
