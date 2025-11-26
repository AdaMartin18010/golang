package network

import (
	"testing"
)

func TestIsValidIP(t *testing.T) {
	if !IsValidIP("192.168.1.1") {
		t.Error("Expected valid IP")
	}
	if IsValidIP("invalid") {
		t.Error("Expected invalid IP")
	}
}

func TestIsIPv4(t *testing.T) {
	if !IsIPv4("192.168.1.1") {
		t.Error("Expected IPv4")
	}
	if IsIPv4("2001:db8::1") {
		t.Error("Expected not IPv4")
	}
}

func TestIsIPv6(t *testing.T) {
	if !IsIPv6("2001:db8::1") {
		t.Error("Expected IPv6")
	}
	if IsIPv6("192.168.1.1") {
		t.Error("Expected not IPv6")
	}
}

func TestIsPrivateIP(t *testing.T) {
	if !IsPrivateIP("192.168.1.1") {
		t.Error("Expected private IP")
	}
	if IsPrivateIP("8.8.8.8") {
		t.Error("Expected not private IP")
	}
}

func TestIPToInt(t *testing.T) {
	ip, err := IPToInt("192.168.1.1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if ip == 0 {
		t.Error("Expected non-zero IP")
	}
}

func TestIntToIP(t *testing.T) {
	ip := IntToIP(3232235777) // 192.168.1.1
	if ip != "192.168.1.1" {
		t.Errorf("Expected '192.168.1.1', got %s", ip)
	}
}

func TestGetLocalIPs(t *testing.T) {
	ips, err := GetLocalIPs()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(ips) == 0 {
		t.Error("Expected at least one local IP")
	}
}

func TestIsIPInCIDR(t *testing.T) {
	if !IsIPInCIDR("192.168.1.1", "192.168.1.0/24") {
		t.Error("Expected IP in CIDR")
	}
	if IsIPInCIDR("10.0.0.1", "192.168.1.0/24") {
		t.Error("Expected IP not in CIDR")
	}
}

func TestValidatePort(t *testing.T) {
	if !ValidatePort(8080) {
		t.Error("Expected valid port")
	}
	if ValidatePort(0) {
		t.Error("Expected invalid port")
	}
	if ValidatePort(65536) {
		t.Error("Expected invalid port")
	}
}

func TestValidateHost(t *testing.T) {
	if !ValidateHost("localhost") {
		t.Error("Expected valid host")
	}
	if !ValidateHost("192.168.1.1") {
		t.Error("Expected valid host")
	}
}

func TestFormatMAC(t *testing.T) {
	formatted := FormatMAC("00:11:22:33:44:55")
	if formatted == "" {
		t.Error("Expected formatted MAC")
	}
}

func TestIsValidMAC(t *testing.T) {
	if !IsValidMAC("00:11:22:33:44:55") {
		t.Error("Expected valid MAC")
	}
	if IsValidMAC("invalid") {
		t.Error("Expected invalid MAC")
	}
}

func TestGetFreePort(t *testing.T) {
	port, err := GetFreePort()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if port <= 0 || port > 65535 {
		t.Errorf("Expected valid port, got %d", port)
	}
}
