package framework

import (
	"testing"
	"time"
)

func TestFactory_String(t *testing.T) {
	f := NewFactory()

	str := f.String(10)
	if len(str) != 10 {
		t.Errorf("Expected string length 10, got %d", len(str))
	}

	str2 := f.String(10)
	if str == str2 {
		t.Error("Generated strings should be different")
	}
}

func TestFactory_Int(t *testing.T) {
	f := NewFactory()

	for i := 0; i < 100; i++ {
		val := f.Int(1, 100)
		if val < 1 || val >= 100 {
			t.Errorf("Value %d is out of range [1, 100)", val)
		}
	}
}

func TestFactory_Email(t *testing.T) {
	f := NewFactory()

	email := f.Email()
	if email == "" {
		t.Error("Email should not be empty")
	}

	// 简单验证邮箱格式
	if len(email) < 5 {
		t.Error("Email should have reasonable length")
	}
}

func TestFactory_Phone(t *testing.T) {
	f := NewFactory()

	phone := f.Phone()
	if len(phone) != 11 {
		t.Errorf("Expected phone length 11, got %d", len(phone))
	}
}

func TestFactory_Time(t *testing.T) {
	f := NewFactory()

	start := time.Now().Add(-24 * time.Hour)
	end := time.Now()

	tm := f.Time(start, end)
	if tm.Before(start) || tm.After(end) {
		t.Errorf("Time %v is out of range [%v, %v]", tm, start, end)
	}
}

func TestFactory_Date(t *testing.T) {
	f := NewFactory()

	date := f.Date()
	now := time.Now()

	if date.After(now) {
		t.Error("Date should be in the past")
	}

	if date.Before(now.AddDate(-2, 0, 0)) {
		t.Error("Date should be within reasonable range")
	}
}

func TestFactory_Slice(t *testing.T) {
	f := NewFactory()

	slice := f.Slice(5, func() interface{} {
		return f.String(10)
	})

	if len(slice) != 5 {
		t.Errorf("Expected slice length 5, got %d", len(slice))
	}
}

func TestFactory_StringSlice(t *testing.T) {
	f := NewFactory()

	slice := f.StringSlice(10)
	if len(slice) != 10 {
		t.Errorf("Expected slice length 10, got %d", len(slice))
	}
}

func TestUserFactory_User(t *testing.T) {
	f := NewUserFactory()

	user := f.User()
	if user["id"] == nil {
		t.Error("User should have id")
	}
	if user["email"] == nil {
		t.Error("User should have email")
	}
}

func TestOAuth2ClientFactory_Client(t *testing.T) {
	f := NewOAuth2ClientFactory()

	client := f.Client()
	if client["id"] == nil {
		t.Error("Client should have id")
	}
	if client["secret"] == nil {
		t.Error("Client should have secret")
	}
}

func TestTokenFactory_Token(t *testing.T) {
	f := NewTokenFactory()

	token := f.Token()
	if token["access_token"] == nil {
		t.Error("Token should have access_token")
	}
	if token["expires_at"] == nil {
		t.Error("Token should have expires_at")
	}
}

func TestAuditLogFactory_AuditLog(t *testing.T) {
	f := NewAuditLogFactory()

	log := f.AuditLog()
	if log["id"] == nil {
		t.Error("Audit log should have id")
	}
	if log["action"] == nil {
		t.Error("Audit log should have action")
	}
}
