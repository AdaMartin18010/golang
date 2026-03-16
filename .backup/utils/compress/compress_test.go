package compress

import (
	"bytes"
	"testing"
)

func TestGzipCompressDecompress(t *testing.T) {
	original := []byte("hello world")
	compressed, err := GzipCompress(original)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	decompressed, err := GzipDecompress(compressed)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !bytes.Equal(original, decompressed) {
		t.Errorf("Expected %s, got %s", string(original), string(decompressed))
	}
}

func TestZlibCompressDecompress(t *testing.T) {
	original := []byte("hello world")
	compressed, err := ZlibCompress(original)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	decompressed, err := ZlibDecompress(compressed)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !bytes.Equal(original, decompressed) {
		t.Errorf("Expected %s, got %s", string(original), string(decompressed))
	}
}

func TestCompressLevel(t *testing.T) {
	original := []byte("hello world")
	compressed, err := CompressLevel(original, 6)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	decompressed, err := GzipDecompress(compressed)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !bytes.Equal(original, decompressed) {
		t.Errorf("Expected %s, got %s", string(original), string(decompressed))
	}
}

func TestCompressBest(t *testing.T) {
	original := []byte("hello world")
	compressed, err := CompressBest(original)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	decompressed, err := GzipDecompress(compressed)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !bytes.Equal(original, decompressed) {
		t.Errorf("Expected %s, got %s", string(original), string(decompressed))
	}
}

func TestCompressFast(t *testing.T) {
	original := []byte("hello world")
	compressed, err := CompressFast(original)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	decompressed, err := GzipDecompress(compressed)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !bytes.Equal(original, decompressed) {
		t.Errorf("Expected %s, got %s", string(original), string(decompressed))
	}
}

func TestIsGzip(t *testing.T) {
	original := []byte("hello world")
	compressed, _ := GzipCompress(original)

	if !IsGzip(compressed) {
		t.Error("Expected gzip format")
	}

	if IsGzip(original) {
		t.Error("Expected not gzip format")
	}
}

func TestGetCompressionRatio(t *testing.T) {
	ratio := GetCompressionRatio(100, 50)
	if ratio != 50.0 {
		t.Errorf("Expected 50.0, got %f", ratio)
	}
}

func TestGetCompressionSavings(t *testing.T) {
	savings := GetCompressionSavings(100, 50)
	if savings != 50 {
		t.Errorf("Expected 50, got %d", savings)
	}
}

func TestCompressStream(t *testing.T) {
	original := []byte("hello world")
	reader := bytes.NewReader(original)
	var buf bytes.Buffer

	err := CompressStream(reader, &buf)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	compressed := buf.Bytes()
	if len(compressed) == 0 {
		t.Error("Expected compressed data")
	}
}

func TestDecompressStream(t *testing.T) {
	original := []byte("hello world")
	compressed, _ := GzipCompress(original)

	reader := bytes.NewReader(compressed)
	var buf bytes.Buffer

	err := DecompressStream(reader, &buf)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	decompressed := buf.Bytes()
	if !bytes.Equal(original, decompressed) {
		t.Errorf("Expected %s, got %s", string(original), string(decompressed))
	}
}
