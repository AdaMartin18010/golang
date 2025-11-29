package security

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"mime"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	// ErrFileTooLarge 文件过大
	ErrFileTooLarge = errors.New("file too large")
	// ErrInvalidFileType 无效的文件类型
	ErrInvalidFileType = errors.New("invalid file type")
	// ErrInvalidFileName 无效的文件名
	ErrInvalidFileName = errors.New("invalid file name")
)

// FileUploadValidator 文件上传验证器
type FileUploadValidator struct {
	maxSize      int64
	allowedTypes []string
	allowedExts  []string
	scanContent  bool
}

// FileUploadConfig 文件上传配置
type FileUploadConfig struct {
	MaxSize      int64    // 最大文件大小（字节）
	AllowedTypes []string // 允许的 MIME 类型
	AllowedExts  []string // 允许的文件扩展名
	ScanContent  bool     // 是否扫描文件内容
}

// DefaultFileUploadConfig 默认文件上传配置
func DefaultFileUploadConfig() FileUploadConfig {
	return FileUploadConfig{
		MaxSize: 10 * 1024 * 1024, // 10MB
		AllowedTypes: []string{
			"image/jpeg",
			"image/png",
			"image/gif",
			"application/pdf",
		},
		AllowedExts: []string{
			".jpg", ".jpeg", ".png", ".gif", ".pdf",
		},
		ScanContent: true,
	}
}

// NewFileUploadValidator 创建文件上传验证器
func NewFileUploadValidator(config FileUploadConfig) *FileUploadValidator {
	if config.MaxSize == 0 {
		config = DefaultFileUploadConfig()
	}

	return &FileUploadValidator{
		maxSize:      config.MaxSize,
		allowedTypes: config.AllowedTypes,
		allowedExts:  config.AllowedExts,
		scanContent:  config.ScanContent,
	}
}

// ValidateFile 验证文件
func (v *FileUploadValidator) ValidateFile(filename string, contentType string, size int64, content io.Reader) error {
	// 验证文件大小
	if size > v.maxSize {
		return ErrFileTooLarge
	}

	// 验证文件名
	if err := v.validateFileName(filename); err != nil {
		return err
	}

	// 验证文件扩展名
	ext := strings.ToLower(filepath.Ext(filename))
	if !v.isAllowedExt(ext) {
		return ErrInvalidFileType
	}

	// 验证 MIME 类型
	if contentType != "" {
		if !v.isAllowedType(contentType) {
			return ErrInvalidFileType
		}
	}

	// 扫描文件内容（如果启用）
	if v.scanContent && content != nil {
		if err := v.scanFileContent(content, ext); err != nil {
			return err
		}
	}

	return nil
}

// validateFileName 验证文件名
func (v *FileUploadValidator) validateFileName(filename string) error {
	if filename == "" {
		return ErrInvalidFileName
	}

	// 检查路径遍历攻击
	if strings.Contains(filename, "..") {
		return ErrInvalidFileName
	}

	// 检查绝对路径
	if filepath.IsAbs(filename) {
		return ErrInvalidFileName
	}

	// 检查危险字符
	dangerousChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, char := range dangerousChars {
		if strings.Contains(filename, char) {
			return ErrInvalidFileName
		}
	}

	return nil
}

// isAllowedExt 检查扩展名是否允许
func (v *FileUploadValidator) isAllowedExt(ext string) bool {
	if len(v.allowedExts) == 0 {
		return true // 如果没有限制，允许所有
	}

	for _, allowed := range v.allowedExts {
		if ext == allowed {
			return true
		}
	}

	return false
}

// isAllowedType 检查 MIME 类型是否允许
func (v *FileUploadValidator) isAllowedType(contentType string) bool {
	if len(v.allowedTypes) == 0 {
		return true // 如果没有限制，允许所有
	}

	// 解析 MIME 类型（可能包含参数）
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return false
	}

	for _, allowed := range v.allowedTypes {
		if mediaType == allowed {
			return true
		}
	}

	return false
}

// scanFileContent 扫描文件内容
func (v *FileUploadValidator) scanFileContent(content io.Reader, ext string) error {
	// 读取文件头（魔数）
	header := make([]byte, 512)
	n, err := content.Read(header)
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to read file header: %w", err)
	}

	// 验证文件头是否匹配扩展名
	if !v.validateFileHeader(header[:n], ext) {
		return ErrInvalidFileType
	}

	return nil
}

// validateFileHeader 验证文件头
func (v *FileUploadValidator) validateFileHeader(header []byte, ext string) bool {
	if len(header) < 4 {
		return false
	}

	// 常见文件类型的魔数
	fileSignatures := map[string][]byte{
		".jpg":  {0xFF, 0xD8, 0xFF},
		".jpeg": {0xFF, 0xD8, 0xFF},
		".png":  {0x89, 0x50, 0x4E, 0x47},
		".gif":  {0x47, 0x49, 0x46, 0x38},
		".pdf":  {0x25, 0x50, 0x44, 0x46}, // %PDF
	}

	signature, exists := fileSignatures[ext]
	if !exists {
		return true // 未知类型，不验证
	}

	if len(header) < len(signature) {
		return false
	}

	for i, b := range signature {
		if header[i] != b {
			return false
		}
	}

	return true
}

// GenerateSafeFileName 生成安全的文件名
func (v *FileUploadValidator) GenerateSafeFileName(originalName string) string {
	// 移除路径
	filename := filepath.Base(originalName)

	// 移除特殊字符
	filename = strings.ReplaceAll(filename, " ", "_")
	filename = regexp.MustCompile(`[^a-zA-Z0-9._-]`).ReplaceAllString(filename, "")

	// 限制长度
	if len(filename) > 255 {
		ext := filepath.Ext(filename)
		name := filename[:255-len(ext)]
		filename = name + ext
	}

	return filename
}

// CalculateFileHash 计算文件哈希
func (v *FileUploadValidator) CalculateFileHash(content io.Reader) (string, error) {
	hash := md5.New()
	if _, err := io.Copy(hash, content); err != nil {
		return "", fmt.Errorf("failed to calculate hash: %w", err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
