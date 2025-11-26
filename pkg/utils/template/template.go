package template

import (
	"bytes"
	"fmt"
	html "html/template"
	"io"
	"os"
	"strings"
	text "text/template"
)

// TextTemplate 文本模板
type TextTemplate struct {
	tmpl *text.Template
}

// NewTextTemplate 创建新的文本模板
func NewTextTemplate(name string) *TextTemplate {
	return &TextTemplate{
		tmpl: text.New(name),
	}
}

// Parse 解析模板字符串
func (t *TextTemplate) Parse(text string) (*TextTemplate, error) {
	var err error
	t.tmpl, err = t.tmpl.Parse(text)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// ParseFiles 解析模板文件
func (t *TextTemplate) ParseFiles(filenames ...string) (*TextTemplate, error) {
	var err error
	t.tmpl, err = t.tmpl.ParseFiles(filenames...)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// ParseGlob 解析匹配的模板文件
func (t *TextTemplate) ParseGlob(pattern string) (*TextTemplate, error) {
	var err error
	t.tmpl, err = t.tmpl.ParseGlob(pattern)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// Execute 执行模板
func (t *TextTemplate) Execute(wr io.Writer, data interface{}) error {
	return t.tmpl.Execute(wr, data)
}

// ExecuteToString 执行模板并返回字符串
func (t *TextTemplate) ExecuteToString(data interface{}) (string, error) {
	var buf bytes.Buffer
	err := t.tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// ExecuteToBytes 执行模板并返回字节数组
func (t *TextTemplate) ExecuteToBytes(data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	err := t.tmpl.Execute(&buf, data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// ExecuteToFile 执行模板并写入文件
func (t *TextTemplate) ExecuteToFile(data interface{}, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	return t.tmpl.Execute(file, data)
}

// HTMLTemplate HTML模板
type HTMLTemplate struct {
	tmpl *html.Template
}

// NewHTMLTemplate 创建新的HTML模板
func NewHTMLTemplate(name string) *HTMLTemplate {
	return &HTMLTemplate{
		tmpl: html.Template.New(name),
	}
}

// Parse 解析模板字符串
func (t *HTMLTemplate) Parse(text string) (*HTMLTemplate, error) {
	var err error
	t.tmpl, err = t.tmpl.Parse(text)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// ParseFiles 解析模板文件
func (t *HTMLTemplate) ParseFiles(filenames ...string) (*HTMLTemplate, error) {
	var err error
	t.tmpl, err = t.tmpl.ParseFiles(filenames...)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// ParseGlob 解析匹配的模板文件
func (t *HTMLTemplate) ParseGlob(pattern string) (*HTMLTemplate, error) {
	var err error
	t.tmpl, err = t.tmpl.ParseGlob(pattern)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// Execute 执行模板
func (t *HTMLTemplate) Execute(wr io.Writer, data interface{}) error {
	return t.tmpl.Execute(wr, data)
}

// ExecuteToString 执行模板并返回字符串
func (t *HTMLTemplate) ExecuteToString(data interface{}) (string, error) {
	var buf bytes.Buffer
	err := t.tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// ExecuteToBytes 执行模板并返回字节数组
func (t *HTMLTemplate) ExecuteToBytes(data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	err := t.tmpl.Execute(&buf, data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// ExecuteToFile 执行模板并写入文件
func (t *HTMLTemplate) ExecuteToFile(data interface{}, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	return t.tmpl.Execute(file, data)
}

// Render 渲染文本模板
func Render(templateText string, data interface{}) (string, error) {
	tmpl, err := text.New("").Parse(templateText)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// RenderHTML 渲染HTML模板
func RenderHTML(templateText string, data interface{}) (string, error) {
	tmpl, err := html.New("").Parse(templateText)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// RenderFile 从文件渲染模板
func RenderFile(filename string, data interface{}) (string, error) {
	tmpl, err := text.ParseFiles(filename)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// RenderHTMLFile 从文件渲染HTML模板
func RenderHTMLFile(filename string, data interface{}) (string, error) {
	tmpl, err := html.ParseFiles(filename)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// Validate 验证模板是否有效
func Validate(templateText string) error {
	_, err := text.New("").Parse(templateText)
	return err
}

// ValidateHTML 验证HTML模板是否有效
func ValidateHTML(templateText string) error {
	_, err := html.New("").Parse(templateText)
	return err
}

// AddFunc 添加自定义函数到文本模板
func (t *TextTemplate) AddFunc(name string, fn interface{}) *TextTemplate {
	t.tmpl = t.tmpl.Funcs(text.FuncMap{name: fn})
	return t
}

// AddFuncs 添加多个自定义函数到文本模板
func (t *TextTemplate) AddFuncs(funcMap text.FuncMap) *TextTemplate {
	t.tmpl = t.tmpl.Funcs(funcMap)
	return t
}

// AddFunc 添加自定义函数到HTML模板
func (t *HTMLTemplate) AddFunc(name string, fn interface{}) *HTMLTemplate {
	t.tmpl = t.tmpl.Funcs(html.FuncMap{name: fn})
	return t
}

// AddFuncs 添加多个自定义函数到HTML模板
func (t *HTMLTemplate) AddFuncs(funcMap html.FuncMap) *HTMLTemplate {
	t.tmpl = t.tmpl.Funcs(funcMap)
	return t
}

// Clone 克隆模板
func (t *TextTemplate) Clone() (*TextTemplate, error) {
	tmpl, err := t.tmpl.Clone()
	if err != nil {
		return nil, err
	}
	return &TextTemplate{tmpl: tmpl}, nil
}

// Clone 克隆HTML模板
func (t *HTMLTemplate) Clone() (*HTMLTemplate, error) {
	tmpl, err := t.tmpl.Clone()
	if err != nil {
		return nil, err
	}
	return &HTMLTemplate{tmpl: tmpl}, nil
}

// Lookup 查找命名模板
func (t *TextTemplate) Lookup(name string) *TextTemplate {
	tmpl := t.tmpl.Lookup(name)
	if tmpl == nil {
		return nil
	}
	return &TextTemplate{tmpl: tmpl}
}

// Lookup 查找命名HTML模板
func (t *HTMLTemplate) Lookup(name string) *HTMLTemplate {
	tmpl := t.tmpl.Lookup(name)
	if tmpl == nil {
		return nil
	}
	return &HTMLTemplate{tmpl: tmpl}
}

// DefinedTemplates 获取所有定义的模板名称
func (t *TextTemplate) DefinedTemplates() string {
	return t.tmpl.DefinedTemplates()
}

// DefinedTemplates 获取所有定义的HTML模板名称
func (t *HTMLTemplate) DefinedTemplates() string {
	return t.tmpl.DefinedTemplates()
}

// CommonFuncMap 常用函数映射
var CommonFuncMap = text.FuncMap{
	"add": func(a, b int) int {
		return a + b
	},
	"sub": func(a, b int) int {
		return a - b
	},
	"mul": func(a, b int) int {
		return a * b
	},
	"div": func(a, b int) int {
		if b == 0 {
			return 0
		}
		return a / b
	},
	"mod": func(a, b int) int {
		if b == 0 {
			return 0
		}
		return a % b
	},
	"eq": func(a, b interface{}) bool {
		return a == b
	},
	"ne": func(a, b interface{}) bool {
		return a != b
	},
	"gt": func(a, b int) bool {
		return a > b
	},
	"ge": func(a, b int) bool {
		return a >= b
	},
	"lt": func(a, b int) bool {
		return a < b
	},
	"le": func(a, b int) bool {
		return a <= b
	},
	"and": func(a, b bool) bool {
		return a && b
	},
	"or": func(a, b bool) bool {
		return a || b
	},
	"not": func(a bool) bool {
		return !a
	},
	"upper": func(s string) string {
		return strings.ToUpper(s)
	},
	"lower": func(s string) string {
		return strings.ToLower(s)
	},
	"title": func(s string) string {
		return strings.Title(s)
	},
	"trim": func(s string) string {
		return strings.TrimSpace(s)
	},
	"join": func(sep string, elems []string) string {
		return strings.Join(elems, sep)
	},
	"split": func(sep, s string) []string {
		return strings.Split(s, sep)
	},
	"contains": func(s, substr string) bool {
		return strings.Contains(s, substr)
	},
	"hasPrefix": func(s, prefix string) bool {
		return strings.HasPrefix(s, prefix)
	},
	"hasSuffix": func(s, suffix string) bool {
		return strings.HasSuffix(s, suffix)
	},
	"replace": func(old, new, s string) string {
		return strings.ReplaceAll(s, old, new)
	},
	"len": func(s interface{}) int {
		switch v := s.(type) {
		case string:
			return len(v)
		case []interface{}:
			return len(v)
		case map[string]interface{}:
			return len(v)
		default:
			return 0
		}
	},
	"default": func(def, val interface{}) interface{} {
		if val == nil || val == "" {
			return def
		}
		return val
	},
	"format": func(format string, args ...interface{}) string {
		return fmt.Sprintf(format, args...)
	},
}

// HTMLCommonFuncMap HTML常用函数映射
var HTMLCommonFuncMap = html.FuncMap{
	"add": func(a, b int) int {
		return a + b
	},
	"sub": func(a, b int) int {
		return a - b
	},
	"mul": func(a, b int) int {
		return a * b
	},
	"div": func(a, b int) int {
		if b == 0 {
			return 0
		}
		return a / b
	},
	"mod": func(a, b int) int {
		if b == 0 {
			return 0
		}
		return a % b
	},
	"eq": func(a, b interface{}) bool {
		return a == b
	},
	"ne": func(a, b interface{}) bool {
		return a != b
	},
	"gt": func(a, b int) bool {
		return a > b
	},
	"ge": func(a, b int) bool {
		return a >= b
	},
	"lt": func(a, b int) bool {
		return a < b
	},
	"le": func(a, b int) bool {
		return a <= b
	},
	"and": func(a, b bool) bool {
		return a && b
	},
	"or": func(a, b bool) bool {
		return a || b
	},
	"not": func(a bool) bool {
		return !a
	},
	"upper": func(s string) string {
		return strings.ToUpper(s)
	},
	"lower": func(s string) string {
		return strings.ToLower(s)
	},
	"title": func(s string) string {
		return strings.Title(s)
	},
	"trim": func(s string) string {
		return strings.TrimSpace(s)
	},
	"join": func(sep string, elems []string) string {
		return strings.Join(elems, sep)
	},
	"split": func(sep, s string) []string {
		return strings.Split(s, sep)
	},
	"contains": func(s, substr string) bool {
		return strings.Contains(s, substr)
	},
	"hasPrefix": func(s, prefix string) bool {
		return strings.HasPrefix(s, prefix)
	},
	"hasSuffix": func(s, suffix string) bool {
		return strings.HasSuffix(s, suffix)
	},
	"replace": func(old, new, s string) string {
		return strings.ReplaceAll(s, old, new)
	},
	"len": func(s interface{}) int {
		switch v := s.(type) {
		case string:
			return len(v)
		case []interface{}:
			return len(v)
		case map[string]interface{}:
			return len(v)
		default:
			return 0
		}
	},
	"default": func(def, val interface{}) interface{} {
		if val == nil || val == "" {
			return def
		}
		return val
	},
	"format": func(format string, args ...interface{}) string {
		return fmt.Sprintf(format, args...)
	},
}
