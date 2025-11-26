package table

import (
	"fmt"
	"strings"
)

// Table 表格
type Table struct {
	headers []string
	rows    [][]string
	widths  []int
}

// NewTable 创建表格
func NewTable(headers ...string) *Table {
	t := &Table{
		headers: headers,
		rows:    make([][]string, 0),
		widths:  make([]int, len(headers)),
	}
	t.calculateWidths()
	return t
}

// AddRow 添加行
func (t *Table) AddRow(row ...string) {
	// 确保行长度与表头一致
	for len(row) < len(t.headers) {
		row = append(row, "")
	}
	if len(row) > len(t.headers) {
		row = row[:len(t.headers)]
	}
	t.rows = append(t.rows, row)
	t.calculateWidths()
}

// AddRows 添加多行
func (t *Table) AddRows(rows [][]string) {
	for _, row := range rows {
		t.AddRow(row...)
	}
}

// calculateWidths 计算列宽
func (t *Table) calculateWidths() {
	// 重置宽度
	for i := range t.widths {
		t.widths[i] = 0
	}

	// 计算表头宽度
	for i, header := range t.headers {
		if len(header) > t.widths[i] {
			t.widths[i] = len(header)
		}
	}

	// 计算行宽度
	for _, row := range t.rows {
		for i, cell := range row {
			if i < len(t.widths) {
				if len(cell) > t.widths[i] {
					t.widths[i] = len(cell)
				}
			}
		}
	}
}

// Render 渲染表格
func (t *Table) Render() string {
	var buf strings.Builder

	// 渲染分隔线
	separator := t.renderSeparator()
	buf.WriteString(separator)
	buf.WriteString("\n")

	// 渲染表头
	buf.WriteString(t.renderRow(t.headers))
	buf.WriteString("\n")
	buf.WriteString(separator)
	buf.WriteString("\n")

	// 渲染数据行
	for _, row := range t.rows {
		buf.WriteString(t.renderRow(row))
		buf.WriteString("\n")
	}

	// 渲染结束分隔线
	buf.WriteString(separator)

	return buf.String()
}

// renderRow 渲染行
func (t *Table) renderRow(row []string) string {
	var cells []string
	for i, cell := range row {
		if i < len(t.widths) {
			cells = append(cells, fmt.Sprintf(" %-*s ", t.widths[i], cell))
		}
	}
	return "|" + strings.Join(cells, "|") + "|"
}

// renderSeparator 渲染分隔线
func (t *Table) renderSeparator() string {
	var parts []string
	for _, width := range t.widths {
		parts = append(parts, strings.Repeat("-", width+2))
	}
	return "+" + strings.Join(parts, "+") + "+"
}

// Print 打印表格
func (t *Table) Print() {
	fmt.Print(t.Render())
}

// String 返回表格字符串
func (t *Table) String() string {
	return t.Render()
}

// SimpleTable 简单表格（无边框）
type SimpleTable struct {
	headers []string
	rows    [][]string
	widths  []int
	sep     string
}

// NewSimpleTable 创建简单表格
func NewSimpleTable(headers ...string) *SimpleTable {
	return &SimpleTable{
		headers: headers,
		rows:    make([][]string, 0),
		widths:  make([]int, len(headers)),
		sep:     "  ",
	}
}

// SetSeparator 设置分隔符
func (st *SimpleTable) SetSeparator(sep string) {
	st.sep = sep
}

// AddRow 添加行
func (st *SimpleTable) AddRow(row ...string) {
	for len(row) < len(st.headers) {
		row = append(row, "")
	}
	if len(row) > len(st.headers) {
		row = row[:len(st.headers)]
	}
	st.rows = append(st.rows, row)
	st.calculateWidths()
}

// calculateWidths 计算列宽
func (st *SimpleTable) calculateWidths() {
	for i := range st.widths {
		st.widths[i] = 0
	}

	for i, header := range st.headers {
		if len(header) > st.widths[i] {
			st.widths[i] = len(header)
		}
	}

	for _, row := range st.rows {
		for i, cell := range row {
			if i < len(st.widths) {
				if len(cell) > st.widths[i] {
					st.widths[i] = len(cell)
				}
			}
		}
	}
}

// Render 渲染表格
func (st *SimpleTable) Render() string {
	var buf strings.Builder

	// 渲染表头
	buf.WriteString(st.renderRow(st.headers))
	buf.WriteString("\n")

	// 渲染数据行
	for _, row := range st.rows {
		buf.WriteString(st.renderRow(row))
		buf.WriteString("\n")
	}

	return buf.String()
}

// renderRow 渲染行
func (st *SimpleTable) renderRow(row []string) string {
	var cells []string
	for i, cell := range row {
		if i < len(st.widths) {
			cells = append(cells, fmt.Sprintf("%-*s", st.widths[i], cell))
		}
	}
	return strings.Join(cells, st.sep)
}

// Print 打印表格
func (st *SimpleTable) Print() {
	fmt.Print(st.Render())
}

// String 返回表格字符串
func (st *SimpleTable) String() string {
	return st.Render()
}

// FormatTable 格式化表格（快捷函数）
func FormatTable(headers []string, rows [][]string) string {
	t := NewTable(headers...)
	t.AddRows(rows)
	return t.Render()
}

// PrintTable 打印表格（快捷函数）
func PrintTable(headers []string, rows [][]string) {
	t := NewTable(headers...)
	t.AddRows(rows)
	t.Print()
}

// FormatSimpleTable 格式化简单表格（快捷函数）
func FormatSimpleTable(headers []string, rows [][]string) string {
	st := NewSimpleTable(headers...)
	st.AddRows(rows)
	return st.Render()
}

// PrintSimpleTable 打印简单表格（快捷函数）
func PrintSimpleTable(headers []string, rows [][]string) {
	st := NewSimpleTable(headers...)
	st.AddRows(rows)
	st.Print()
}
