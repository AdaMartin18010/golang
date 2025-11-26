package network

import (
	"fmt"
	"net"
	"os"
	"strings"
)

// IsValidIP 检查是否为有效的IP地址
func IsValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

// IsIPv4 检查是否为IPv4地址
func IsIPv4(ip string) bool {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return false
	}
	return parsed.To4() != nil
}

// IsIPv6 检查是否为IPv6地址
func IsIPv6(ip string) bool {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return false
	}
	return parsed.To4() == nil && parsed.To16() != nil
}

// IsPrivateIP 检查是否为私有IP地址
func IsPrivateIP(ip string) bool {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return false
	}
	return parsed.IsPrivate()
}

// IsLoopback 检查是否为回环地址
func IsLoopback(ip string) bool {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return false
	}
	return parsed.IsLoopback()
}

// IsMulticast 检查是否为多播地址
func IsMulticast(ip string) bool {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return false
	}
	return parsed.IsMulticast()
}

// IsUnspecified 检查是否为未指定地址
func IsUnspecified(ip string) bool {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return false
	}
	return parsed.IsUnspecified()
}

// ParseIP 解析IP地址
func ParseIP(ip string) net.IP {
	return net.ParseIP(ip)
}

// IPToInt 将IPv4地址转换为整数
func IPToInt(ip string) (uint32, error) {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return 0, fmt.Errorf("invalid IP address: %s", ip)
	}
	ipv4 := parsed.To4()
	if ipv4 == nil {
		return 0, fmt.Errorf("not an IPv4 address: %s", ip)
	}
	return uint32(ipv4[0])<<24 | uint32(ipv4[1])<<16 | uint32(ipv4[2])<<8 | uint32(ipv4[3]), nil
}

// IntToIP 将整数转换为IPv4地址
func IntToIP(ip uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d",
		byte(ip>>24),
		byte(ip>>16),
		byte(ip>>8),
		byte(ip))
}

// GetLocalIP 获取本地IP地址
func GetLocalIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}

// GetLocalIPs 获取所有本地IP地址
func GetLocalIPs() ([]string, error) {
	var ips []string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP.String())
			}
		}
	}
	return ips, nil
}

// GetHostname 获取主机名
func GetHostname() (string, error) {
	return os.Hostname()
}

// ResolveIP 解析主机名到IP地址
func ResolveIP(hostname string) ([]string, error) {
	ips, err := net.LookupIP(hostname)
	if err != nil {
		return nil, err
	}
	var result []string
	for _, ip := range ips {
		result = append(result, ip.String())
	}
	return result, nil
}

// ResolveHostname 解析IP地址到主机名
func ResolveHostname(ip string) (string, error) {
	names, err := net.LookupAddr(ip)
	if err != nil {
		return "", err
	}
	if len(names) == 0 {
		return "", fmt.Errorf("no hostname found for IP: %s", ip)
	}
	return strings.TrimSuffix(names[0], "."), nil
}

// IsPortOpen 检查端口是否开放
func IsPortOpen(host string, port int) bool {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// IsPortOpenTimeout 检查端口是否开放（带超时）
func IsPortOpenTimeout(host string, port int, timeoutSeconds int) bool {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", address, timeoutSeconds)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// ParseCIDR 解析CIDR
func ParseCIDR(cidr string) (*net.IPNet, error) {
	_, ipnet, err := net.ParseCIDR(cidr)
	return ipnet, err
}

// IsIPInCIDR 检查IP是否在CIDR范围内
func IsIPInCIDR(ip, cidr string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}
	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return false
	}
	return ipnet.Contains(parsedIP)
}

// GetNetworkInfo 获取网络信息
func GetNetworkInfo() ([]NetworkInterface, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var result []NetworkInterface
	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		var ips []string
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok {
				ips = append(ips, ipnet.IP.String())
			}
		}

		result = append(result, NetworkInterface{
			Name:         iface.Name,
			Index:        iface.Index,
			MTU:          iface.MTU,
			HardwareAddr: iface.HardwareAddr.String(),
			Flags:        iface.Flags.String(),
			IPs:          ips,
		})
	}
	return result, nil
}

// NetworkInterface 网络接口信息
type NetworkInterface struct {
	Name         string
	Index        int
	MTU          int
	HardwareAddr string
	Flags        string
	IPs          []string
}

// ValidatePort 验证端口号
func ValidatePort(port int) bool {
	return port > 0 && port <= 65535
}

// ValidateHost 验证主机名或IP
func ValidateHost(host string) bool {
	if IsValidIP(host) {
		return true
	}
	// 简单的主机名验证
	if len(host) == 0 || len(host) > 253 {
		return false
	}
	parts := strings.Split(host, ".")
	for _, part := range parts {
		if len(part) == 0 || len(part) > 63 {
			return false
		}
	}
	return true
}

// FormatMAC 格式化MAC地址
func FormatMAC(mac string) string {
	mac = strings.ReplaceAll(mac, ":", "")
	mac = strings.ReplaceAll(mac, "-", "")
	if len(mac) != 12 {
		return mac
	}
	return fmt.Sprintf("%s:%s:%s:%s:%s:%s",
		mac[0:2], mac[2:4], mac[4:6],
		mac[6:8], mac[8:10], mac[10:12])
}

// IsValidMAC 检查是否为有效的MAC地址
func IsValidMAC(mac string) bool {
	_, err := net.ParseMAC(mac)
	return err == nil
}

// GetFreePort 获取可用端口
func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

// Ping 简单的ping实现（TCP连接测试）
func Ping(host string, port int, timeoutSeconds int) bool {
	return IsPortOpenTimeout(host, port, timeoutSeconds)
}

// IsReachable 检查主机是否可达
func IsReachable(host string, port int) bool {
	return IsPortOpen(host, port)
}
