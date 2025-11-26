package validator

import (
	"testing"
)

func TestIsEmail(t *testing.T) {
	tests := []struct {
		email string
		want  bool
	}{
		{"test@example.com", true},
		{"user.name@example.com", true},
		{"user+tag@example.co.uk", true},
		{"invalid", false},
		{"@example.com", false},
		{"test@", false},
		{"", false},
	}

	for _, tt := range tests {
		if got := IsEmail(tt.email); got != tt.want {
			t.Errorf("IsEmail(%q) = %v, want %v", tt.email, got, tt.want)
		}
	}
}

func TestIsPhone(t *testing.T) {
	tests := []struct {
		phone string
		want  bool
	}{
		{"13800138000", true},
		{"15912345678", true},
		{"12345678901", false},
		{"1380013800", false},
		{"", false},
	}

	for _, tt := range tests {
		if got := IsPhone(tt.phone); got != tt.want {
			t.Errorf("IsPhone(%q) = %v, want %v", tt.phone, got, tt.want)
		}
	}
}

func TestIsIDCard(t *testing.T) {
	tests := []struct {
		idCard string
		want   bool
	}{
		{"110101199003075132", true}, // 示例身份证号
		{"123456789012345678", false},
		{"", false},
	}

	for _, tt := range tests {
		if got := IsIDCard(tt.idCard); got != tt.want {
			t.Errorf("IsIDCard(%q) = %v, want %v", tt.idCard, got, tt.want)
		}
	}
}

func TestIsURL(t *testing.T) {
	tests := []struct {
		url  string
		want bool
	}{
		{"https://example.com", true},
		{"http://example.com/path", true},
		{"ftp://example.com", true},
		{"invalid", false},
		{"", false},
	}

	for _, tt := range tests {
		if got := IsURL(tt.url); got != tt.want {
			t.Errorf("IsURL(%q) = %v, want %v", tt.url, got, tt.want)
		}
	}
}

func TestIsIPv4(t *testing.T) {
	tests := []struct {
		ip   string
		want bool
	}{
		{"192.168.1.1", true},
		{"255.255.255.255", true},
		{"0.0.0.0", true},
		{"256.1.1.1", false},
		{"192.168.1", false},
		{"", false},
	}

	for _, tt := range tests {
		if got := IsIPv4(tt.ip); got != tt.want {
			t.Errorf("IsIPv4(%q) = %v, want %v", tt.ip, got, tt.want)
		}
	}
}

func TestIsUUID(t *testing.T) {
	tests := []struct {
		uuid string
		want bool
	}{
		{"550e8400-e29b-41d4-a716-446655440000", true},
		{"550E8400-E29B-41D4-A716-446655440000", true},
		{"invalid", false},
		{"", false},
	}

	for _, tt := range tests {
		if got := IsUUID(tt.uuid); got != tt.want {
			t.Errorf("IsUUID(%q) = %v, want %v", tt.uuid, got, tt.want)
		}
	}
}

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		s    string
		want bool
	}{
		{"", true},
		{"   ", true},
		{"test", false},
	}

	for _, tt := range tests {
		if got := IsEmpty(tt.s); got != tt.want {
			t.Errorf("IsEmpty(%q) = %v, want %v", tt.s, got, tt.want)
		}
	}
}

func TestIsNumeric(t *testing.T) {
	tests := []struct {
		s    string
		want bool
	}{
		{"123", true},
		{"0", true},
		{"abc", false},
		{"12a", false},
		{"", false},
	}

	for _, tt := range tests {
		if got := IsNumeric(tt.s); got != tt.want {
			t.Errorf("IsNumeric(%q) = %v, want %v", tt.s, got, tt.want)
		}
	}
}

func TestIsAlpha(t *testing.T) {
	tests := []struct {
		s    string
		want bool
	}{
		{"abc", true},
		{"ABC", true},
		{"abc123", false},
		{"", false},
	}

	for _, tt := range tests {
		if got := IsAlpha(tt.s); got != tt.want {
			t.Errorf("IsAlpha(%q) = %v, want %v", tt.s, got, tt.want)
		}
	}
}

func TestHasLength(t *testing.T) {
	tests := []struct {
		s    string
		min  int
		max  int
		want bool
	}{
		{"test", 3, 5, true},
		{"test", 5, 10, false},
		{"test", 1, 3, false},
	}

	for _, tt := range tests {
		if got := HasLength(tt.s, tt.min, tt.max); got != tt.want {
			t.Errorf("HasLength(%q, %d, %d) = %v, want %v", tt.s, tt.min, tt.max, got, tt.want)
		}
	}
}

func TestIsIn(t *testing.T) {
	tests := []struct {
		value int
		slice []int
		want  bool
	}{
		{1, []int{1, 2, 3}, true},
		{4, []int{1, 2, 3}, false},
		{1, []int{}, false},
	}

	for _, tt := range tests {
		if got := IsIn(tt.value, tt.slice); got != tt.want {
			t.Errorf("IsIn(%d, %v) = %v, want %v", tt.value, tt.slice, got, tt.want)
		}
	}
}

func TestIsBetween(t *testing.T) {
	tests := []struct {
		value int
		min   int
		max   int
		want  bool
	}{
		{5, 1, 10, true},
		{1, 1, 10, true},
		{10, 1, 10, true},
		{0, 1, 10, false},
		{11, 1, 10, false},
	}

	for _, tt := range tests {
		if got := IsBetween(tt.value, tt.min, tt.max); got != tt.want {
			t.Errorf("IsBetween(%d, %d, %d) = %v, want %v", tt.value, tt.min, tt.max, got, tt.want)
		}
	}
}

func TestIsStrongPassword(t *testing.T) {
	tests := []struct {
		password string
		want     bool
	}{
		{"Password123!", true},
		{"password123!", false}, // 缺少大写
		{"PASSWORD123!", false}, // 缺少小写
		{"Password!", false},    // 缺少数字
		{"Password123", false},  // 缺少特殊字符
		{"Pass1!", false},       // 长度不足
		{"", false},
	}

	for _, tt := range tests {
		if got := IsStrongPassword(tt.password); got != tt.want {
			t.Errorf("IsStrongPassword(%q) = %v, want %v", tt.password, got, tt.want)
		}
	}
}

func TestIsChinese(t *testing.T) {
	tests := []struct {
		s    string
		want bool
	}{
		{"中文", true},
		{"中文abc", false},
		{"", false},
	}

	for _, tt := range tests {
		if got := IsChinese(tt.s); got != tt.want {
			t.Errorf("IsChinese(%q) = %v, want %v", tt.s, got, tt.want)
		}
	}
}

func TestIsDate(t *testing.T) {
	tests := []struct {
		date string
		want bool
	}{
		{"2025-11-11", true},
		{"2025-1-1", false},
		{"invalid", false},
		{"", false},
	}

	for _, tt := range tests {
		if got := IsDate(tt.date); got != tt.want {
			t.Errorf("IsDate(%q) = %v, want %v", tt.date, got, tt.want)
		}
	}
}
