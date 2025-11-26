package file

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExists(t *testing.T) {
	// 创建临时文件
	tmpfile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()

	if !Exists(tmpfile.Name()) {
		t.Error("Expected file to exist")
	}

	if Exists("/nonexistent/path") {
		t.Error("Expected file not to exist")
	}
}

func TestIsFile(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()

	if !IsFile(tmpfile.Name()) {
		t.Error("Expected path to be a file")
	}

	tmpdir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpdir)

	if IsFile(tmpdir) {
		t.Error("Expected path not to be a file")
	}
}

func TestIsDir(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpdir)

	if !IsDir(tmpdir) {
		t.Error("Expected path to be a directory")
	}

	tmpfile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()

	if IsDir(tmpfile.Name()) {
		t.Error("Expected path not to be a directory")
	}
}

func TestReadWriteFile(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()

	content := "test content"
	err = WriteFileString(tmpfile.Name(), content, 0644)
	if err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	readContent, err := ReadFileString(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if readContent != content {
		t.Errorf("Expected %s, got %s", content, readContent)
	}
}

func TestCopyFile(t *testing.T) {
	src, err := os.CreateTemp("", "src")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(src.Name())
	src.WriteString("test content")
	src.Close()

	dst := filepath.Join(filepath.Dir(src.Name()), "dst")
	defer os.Remove(dst)

	err = CopyFile(src.Name(), dst)
	if err != nil {
		t.Fatalf("Failed to copy file: %v", err)
	}

	if !Exists(dst) {
		t.Error("Expected destination file to exist")
	}

	content, err := ReadFileString(dst)
	if err != nil {
		t.Fatalf("Failed to read copied file: %v", err)
	}

	if content != "test content" {
		t.Errorf("Expected 'test content', got %s", content)
	}
}

func TestCreateDir(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpdir)

	newDir := filepath.Join(tmpdir, "newdir")
	err = CreateDir(newDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	if !IsDir(newDir) {
		t.Error("Expected directory to exist")
	}
}

func TestGetExt(t *testing.T) {
	tests := []struct {
		filename string
		expected string
	}{
		{"test.txt", ".txt"},
		{"test.go", ".go"},
		{"test", ""},
		{"test.tar.gz", ".gz"},
	}

	for _, tt := range tests {
		result := GetExt(tt.filename)
		if result != tt.expected {
			t.Errorf("GetExt(%s) = %s, expected %s", tt.filename, result, tt.expected)
		}
	}
}

func TestGetBaseName(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"/path/to/file.txt", "file.txt"},
		{"file.txt", "file.txt"},
		{"/path/to/", "to"},
	}

	for _, tt := range tests {
		result := GetBaseName(tt.path)
		if result != tt.expected {
			t.Errorf("GetBaseName(%s) = %s, expected %s", tt.path, result, tt.expected)
		}
	}
}

func TestReadLines(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	content := "line1\nline2\nline3"
	tmpfile.WriteString(content)
	tmpfile.Close()

	lines, err := ReadLines(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to read lines: %v", err)
	}

	if len(lines) != 3 {
		t.Errorf("Expected 3 lines, got %d", len(lines))
	}

	if lines[0] != "line1" {
		t.Errorf("Expected 'line1', got %s", lines[0])
	}
}

func TestWriteLines(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()

	lines := []string{"line1", "line2", "line3"}
	err = WriteLines(tmpfile.Name(), lines, 0644)
	if err != nil {
		t.Fatalf("Failed to write lines: %v", err)
	}

	readLines, err := ReadLines(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to read lines: %v", err)
	}

	if len(readLines) != len(lines) {
		t.Errorf("Expected %d lines, got %d", len(lines), len(readLines))
	}
}
