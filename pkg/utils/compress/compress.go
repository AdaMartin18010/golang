package compress

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"io"
	"os"
)

// GzipCompress gzip压缩
func GzipCompress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	_, err := writer.Write(data)
	if err != nil {
		writer.Close()
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GzipDecompress gzip解压
func GzipDecompress(data []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return io.ReadAll(reader)
}

// GzipCompressToFile gzip压缩到文件
func GzipCompressToFile(data []byte, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := gzip.NewWriter(file)
	defer writer.Close()

	_, err = writer.Write(data)
	return err
}

// GzipDecompressFromFile 从文件gzip解压
func GzipDecompressFromFile(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	return io.ReadAll(reader)
}

// ZlibCompress zlib压缩
func ZlibCompress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := zlib.NewWriter(&buf)
	_, err := writer.Write(data)
	if err != nil {
		writer.Close()
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// ZlibDecompress zlib解压
func ZlibDecompress(data []byte) ([]byte, error) {
	reader, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return io.ReadAll(reader)
}

// CompressLevel gzip压缩（指定压缩级别）
func CompressLevel(data []byte, level int) ([]byte, error) {
	var buf bytes.Buffer
	writer, err := gzip.NewWriterLevel(&buf, level)
	if err != nil {
		return nil, err
	}
	_, err = writer.Write(data)
	if err != nil {
		writer.Close()
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// CompressBest gzip压缩（最佳压缩率）
func CompressBest(data []byte) ([]byte, error) {
	return CompressLevel(data, gzip.BestCompression)
}

// CompressFast gzip压缩（最快速度）
func CompressFast(data []byte) ([]byte, error) {
	return CompressLevel(data, gzip.BestSpeed)
}

// CompressDefault gzip压缩（默认压缩率）
func CompressDefault(data []byte) ([]byte, error) {
	return CompressLevel(data, gzip.DefaultCompression)
}

// CompressNoCompression gzip压缩（不压缩）
func CompressNoCompression(data []byte) ([]byte, error) {
	return CompressLevel(data, gzip.NoCompression)
}

// IsGzip 检查数据是否为gzip格式
func IsGzip(data []byte) bool {
	if len(data) < 2 {
		return false
	}
	return data[0] == 0x1f && data[1] == 0x8b
}

// GetCompressionRatio 获取压缩率
func GetCompressionRatio(originalSize, compressedSize int) float64 {
	if originalSize == 0 {
		return 0
	}
	return float64(compressedSize) / float64(originalSize) * 100
}

// GetCompressionSavings 获取压缩节省的字节数
func GetCompressionSavings(originalSize, compressedSize int) int {
	return originalSize - compressedSize
}

// CompressStream gzip压缩流
func CompressStream(reader io.Reader, writer io.Writer) error {
	gzipWriter := gzip.NewWriter(writer)
	defer gzipWriter.Close()

	_, err := io.Copy(gzipWriter, reader)
	return err
}

// DecompressStream gzip解压流
func DecompressStream(reader io.Reader, writer io.Writer) error {
	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	_, err = io.Copy(writer, gzipReader)
	return err
}
