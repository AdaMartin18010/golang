package url

import (
	"net/url"
	"strings"
)

// Parse 解析URL
func Parse(rawURL string) (*url.URL, error) {
	return url.Parse(rawURL)
}

// ParseRequestURI 解析请求URI
func ParseRequestURI(rawURL string) (*url.URL, error) {
	return url.ParseRequestURI(rawURL)
}

// BuildURL 构建URL
func BuildURL(baseURL string, path string, params map[string]string) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	// 设置路径
	if path != "" {
		u.Path = strings.TrimSuffix(u.Path, "/") + "/" + strings.TrimPrefix(path, "/")
	}

	// 设置查询参数
	if len(params) > 0 {
		q := u.Query()
		for k, v := range params {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
	}

	return u.String(), nil
}

// AddQuery 添加查询参数
func AddQuery(rawURL string, key, value string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	q := u.Query()
	q.Set(key, value)
	u.RawQuery = q.Encode()

	return u.String(), nil
}

// AddQueries 批量添加查询参数
func AddQueries(rawURL string, params map[string]string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	return u.String(), nil
}

// RemoveQuery 移除查询参数
func RemoveQuery(rawURL string, key string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	q := u.Query()
	q.Del(key)
	u.RawQuery = q.Encode()

	return u.String(), nil
}

// GetQuery 获取查询参数值
func GetQuery(rawURL string, key string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	return u.Query().Get(key), nil
}

// GetAllQueries 获取所有查询参数
func GetAllQueries(rawURL string) (map[string][]string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	return u.Query(), nil
}

// Encode 编码URL
func Encode(s string) string {
	return url.QueryEscape(s)
}

// Decode 解码URL
func Decode(s string) (string, error) {
	return url.QueryUnescape(s)
}

// JoinPath 连接路径
func JoinPath(baseURL string, paths ...string) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	path := u.Path
	for _, p := range paths {
		path = strings.TrimSuffix(path, "/") + "/" + strings.TrimPrefix(p, "/")
	}
	u.Path = path

	return u.String(), nil
}

// SetScheme 设置URL协议
func SetScheme(rawURL string, scheme string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	u.Scheme = scheme
	return u.String(), nil
}

// SetHost 设置URL主机
func SetHost(rawURL string, host string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	u.Host = host
	return u.String(), nil
}

// SetPath 设置URL路径
func SetPath(rawURL string, path string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	u.Path = path
	return u.String(), nil
}

// GetScheme 获取URL协议
func GetScheme(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	return u.Scheme, nil
}

// GetHost 获取URL主机
func GetHost(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	return u.Host, nil
}

// GetPath 获取URL路径
func GetPath(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	return u.Path, nil
}

// IsValid 检查URL是否有效
func IsValid(rawURL string) bool {
	_, err := url.ParseRequestURI(rawURL)
	return err == nil
}

// IsAbsolute 检查URL是否为绝对路径
func IsAbsolute(rawURL string) bool {
	u, err := url.Parse(rawURL)
	if err != nil {
		return false
	}
	return u.IsAbs()
}

// Normalize 规范化URL
func Normalize(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	// 规范化路径
	u.Path = strings.TrimSuffix(u.Path, "/")
	if u.Path == "" {
		u.Path = "/"
	}

	return u.String(), nil
}

// Resolve 解析相对URL
func Resolve(baseURL, relativeURL string) (string, error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	rel, err := url.Parse(relativeURL)
	if err != nil {
		return "", err
	}

	return base.ResolveReference(rel).String(), nil
}

// ReplacePath 替换URL路径
func ReplacePath(rawURL string, newPath string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	u.Path = newPath
	return u.String(), nil
}

// ReplaceHost 替换URL主机
func ReplaceHost(rawURL string, newHost string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	u.Host = newHost
	return u.String(), nil
}

// ReplaceScheme 替换URL协议
func ReplaceScheme(rawURL string, newScheme string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	u.Scheme = newScheme
	return u.String(), nil
}

// GetDomain 获取域名（不含端口）
func GetDomain(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	host := u.Host
	if idx := strings.Index(host, ":"); idx != -1 {
		host = host[:idx]
	}

	return host, nil
}

// GetPort 获取端口号
func GetPort(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	host := u.Host
	if idx := strings.Index(host, ":"); idx != -1 {
		return host[idx+1:], nil
	}

	// 根据协议返回默认端口
	switch u.Scheme {
	case "http":
		return "80", nil
	case "https":
		return "443", nil
	default:
		return "", nil
	}
}

// BuildQueryString 构建查询字符串
func BuildQueryString(params map[string]string) string {
	if len(params) == 0 {
		return ""
	}

	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}

	return values.Encode()
}

// ParseQueryString 解析查询字符串
func ParseQueryString(queryString string) (map[string]string, error) {
	values, err := url.ParseQuery(queryString)
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for k, v := range values {
		if len(v) > 0 {
			result[k] = v[0]
		}
	}

	return result, nil
}

// MaskURL 掩码URL（隐藏敏感信息）
func MaskURL(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	// 掩码用户名和密码
	if u.User != nil {
		u.User = nil
	}

	// 掩码查询参数中的敏感字段
	sensitiveKeys := []string{"password", "token", "secret", "key", "api_key"}
	q := u.Query()
	for _, key := range sensitiveKeys {
		if q.Has(key) {
			q.Set(key, "***")
		}
	}
	u.RawQuery = q.Encode()

	return u.String(), nil
}

// IsHTTPS 检查是否为HTTPS
func IsHTTPS(rawURL string) bool {
	scheme, err := GetScheme(rawURL)
	if err != nil {
		return false
	}
	return scheme == "https"
}

// IsHTTP 检查是否为HTTP
func IsHTTP(rawURL string) bool {
	scheme, err := GetScheme(rawURL)
	if err != nil {
		return false
	}
	return scheme == "http"
}

// ToHTTPS 转换为HTTPS
func ToHTTPS(rawURL string) (string, error) {
	return ReplaceScheme(rawURL, "https")
}

// ToHTTP 转换为HTTP
func ToHTTP(rawURL string) (string, error) {
	return ReplaceScheme(rawURL, "http")
}
