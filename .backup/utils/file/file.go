package file

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Exists 检查文件或目录是否存在
func Exists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

// IsFile 检查路径是否为文件
func IsFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// IsDir 检查路径是否为目录
func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// ReadFile 读取文件内容
func ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

// ReadFileString 读取文件内容为字符串
func ReadFileString(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// WriteFile 写入文件内容
func WriteFile(filename string, data []byte, perm os.FileMode) error {
	return os.WriteFile(filename, data, perm)
}

// WriteFileString 写入字符串到文件
func WriteFileString(filename string, content string, perm os.FileMode) error {
	return os.WriteFile(filename, []byte(content), perm)
}

// AppendFile 追加内容到文件
func AppendFile(filename string, data []byte, perm os.FileMode) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, perm)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	return err
}

// AppendFileString 追加字符串到文件
func AppendFileString(filename string, content string, perm os.FileMode) error {
	return AppendFile(filename, []byte(content), perm)
}

// CopyFile 复制文件
func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

// MoveFile 移动文件
func MoveFile(src, dst string) error {
	return os.Rename(src, dst)
}

// DeleteFile 删除文件
func DeleteFile(filename string) error {
	return os.Remove(filename)
}

// DeleteDir 删除目录（递归）
func DeleteDir(dirname string) error {
	return os.RemoveAll(dirname)
}

// CreateDir 创建目录
func CreateDir(dirname string, perm os.FileMode) error {
	return os.MkdirAll(dirname, perm)
}

// GetFileSize 获取文件大小
func GetFileSize(filename string) (int64, error) {
	info, err := os.Stat(filename)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// GetFileMode 获取文件权限
func GetFileMode(filename string) (os.FileMode, error) {
	info, err := os.Stat(filename)
	if err != nil {
		return 0, err
	}
	return info.Mode(), nil
}

// ListFiles 列出目录中的文件
func ListFiles(dirname string) ([]string, error) {
	entries, err := os.ReadDir(dirname)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}

	return files, nil
}

// ListDirs 列出目录中的子目录
func ListDirs(dirname string) ([]string, error) {
	entries, err := os.ReadDir(dirname)
	if err != nil {
		return nil, err
	}

	var dirs []string
	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, entry.Name())
		}
	}

	return dirs, nil
}

// ListAll 列出目录中的所有条目
func ListAll(dirname string) ([]string, error) {
	entries, err := os.ReadDir(dirname)
	if err != nil {
		return nil, err
	}

	var items []string
	for _, entry := range entries {
		items = append(items, entry.Name())
	}

	return items, nil
}

// WalkFiles 遍历目录中的所有文件
func WalkFiles(root string, fn func(string) error) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return fn(path)
		}
		return nil
	})
}

// WalkDirs 遍历目录中的所有子目录
func WalkDirs(root string, fn func(string) error) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && path != root {
			return fn(path)
		}
		return nil
	})
}

// GetExt 获取文件扩展名
func GetExt(filename string) string {
	return filepath.Ext(filename)
}

// GetBaseName 获取文件名（不含路径）
func GetBaseName(path string) string {
	return filepath.Base(path)
}

// GetDirName 获取目录名
func GetDirName(path string) string {
	return filepath.Dir(path)
}

// JoinPath 连接路径
func JoinPath(elem ...string) string {
	return filepath.Join(elem...)
}

// CleanPath 清理路径
func CleanPath(path string) string {
	return filepath.Clean(path)
}

// AbsPath 获取绝对路径
func AbsPath(path string) (string, error) {
	return filepath.Abs(path)
}

// RelPath 获取相对路径
func RelPath(basepath, targpath string) (string, error) {
	return filepath.Rel(basepath, targpath)
}

// MatchPattern 匹配文件模式
func MatchPattern(pattern, name string) (bool, error) {
	return filepath.Match(pattern, name)
}

// Glob 匹配文件模式（支持通配符）
func Glob(pattern string) ([]string, error) {
	return filepath.Glob(pattern)
}

// EnsureDir 确保目录存在，如果不存在则创建
func EnsureDir(dirname string, perm os.FileMode) error {
	if !Exists(dirname) {
		return CreateDir(dirname, perm)
	}
	if !IsDir(dirname) {
		return errors.New("path exists but is not a directory")
	}
	return nil
}

// EnsureFileDir 确保文件所在目录存在
func EnsureFileDir(filename string, perm os.FileMode) error {
	dir := GetDirName(filename)
	return EnsureDir(dir, perm)
}

// ReadLines 读取文件的所有行
func ReadLines(filename string) ([]string, error) {
	content, err := ReadFileString(filename)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(content, "\n")
	// 移除最后的空行（如果存在）
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	return lines, nil
}

// WriteLines 写入多行到文件
func WriteLines(filename string, lines []string, perm os.FileMode) error {
	content := strings.Join(lines, "\n")
	return WriteFileString(filename, content, perm)
}

// Chmod 修改文件权限
func Chmod(filename string, mode os.FileMode) error {
	return os.Chmod(filename, mode)
}

// Chown 修改文件所有者（Unix系统）
func Chown(filename string, uid, gid int) error {
	return os.Chown(filename, uid, gid)
}
