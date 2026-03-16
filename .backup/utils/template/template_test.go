package template

import (
	"strings"
	"testing"
)

func TestTextTemplate(t *testing.T) {
	tmpl := NewTextTemplate("test")
	tmpl, err := tmpl.Parse("Hello, {{.Name}}!")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	data := map[string]string{"Name": "World"}
	result, err := tmpl.ExecuteToString(data)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if result != "Hello, World!" {
		t.Errorf("Expected 'Hello, World!', got %s", result)
	}
}

func TestHTMLTemplate(t *testing.T) {
	tmpl := NewHTMLTemplate("test")
	tmpl, err := tmpl.Parse("<h1>Hello, {{.Name}}!</h1>")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	data := map[string]string{"Name": "World"}
	result, err := tmpl.ExecuteToString(data)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if result != "<h1>Hello, World!</h1>" {
		t.Errorf("Expected '<h1>Hello, World!</h1>', got %s", result)
	}
}

func TestRender(t *testing.T) {
	templateText := "Hello, {{.Name}}!"
	data := map[string]string{"Name": "World"}
	result, err := Render(templateText, data)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if result != "Hello, World!" {
		t.Errorf("Expected 'Hello, World!', got %s", result)
	}
}

func TestRenderHTML(t *testing.T) {
	templateText := "<h1>Hello, {{.Name}}!</h1>"
	data := map[string]string{"Name": "World"}
	result, err := RenderHTML(templateText, data)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if result != "<h1>Hello, World!</h1>" {
		t.Errorf("Expected '<h1>Hello, World!</h1>', got %s", result)
	}
}

func TestValidate(t *testing.T) {
	err := Validate("Hello, {{.Name}}!")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	err = Validate("Hello, {{.Name}")
	if err == nil {
		t.Error("Expected error for invalid template")
	}
}

func TestValidateHTML(t *testing.T) {
	err := ValidateHTML("<h1>Hello, {{.Name}}!</h1>")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestAddFunc(t *testing.T) {
	tmpl := NewTextTemplate("test")
	tmpl = tmpl.AddFunc("upper", func(s string) string {
		return strings.ToUpper(s)
	})
	tmpl, err := tmpl.Parse("Hello, {{upper .Name}}!")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	data := map[string]string{"Name": "world"}
	result, err := tmpl.ExecuteToString(data)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if result != "Hello, WORLD!" {
		t.Errorf("Expected 'Hello, WORLD!', got %s", result)
	}
}

func TestCommonFuncMap(t *testing.T) {
	tmpl := NewTextTemplate("test")
	tmpl = tmpl.AddFuncs(CommonFuncMap)
	tmpl, err := tmpl.Parse("{{add 1 2}}")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	result, err := tmpl.ExecuteToString(nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if result != "3" {
		t.Errorf("Expected '3', got %s", result)
	}
}
