package utils

import (
	"fmt"
	"math"
	"strings"
	"sync"
	"time"
)

// Size represents file size with human-readable formatting
type Size struct {
	bytes int64
}

var sizeUnits = []string{"B", "KB", "MB", "GB", "TB", "PB"}

func NewSize(bytes int64) *Size {
	return &Size{bytes: bytes}
}

func (s *Size) String() string {
	if s.bytes == 0 {
		return "0 B"
	}
	
	size := float64(s.bytes)
	unitIndex := 0
	
	for size >= 1024 && unitIndex < len(sizeUnits)-1 {
		size /= 1024
		unitIndex++
	}
	
	if size < 10 {
		return fmt.Sprintf("%.1f %s", size, sizeUnits[unitIndex])
	}
	return fmt.Sprintf("%.0f %s", size, sizeUnits[unitIndex])
}

// Table represents a formatted table
type Table struct {
	headers []string
	rows    [][]string
	widths  []int
}

func NewTable(headers []string) *Table {
	t := &Table{
		headers: headers,
		rows:    [][]string{},
		widths:  make([]int, len(headers)),
	}
	
	// Initialize column widths with header lengths
	for i, header := range headers {
		t.widths[i] = len(header)
	}
	
	return t
}

func (t *Table) AddRow(row []string) {
	if len(row) != len(t.headers) {
		return
	}
	
	t.rows = append(t.rows, row)
	
	// Update column widths
	for i, cell := range row {
		if len(cell) > t.widths[i] {
			t.widths[i] = len(cell)
		}
	}
}

func (t *Table) String() string {
	var result strings.Builder
	
	// Print headers
	t.printRow(&result, t.headers)
	
	// Print separator
	for i, width := range t.widths {
		if i > 0 {
			result.WriteString("+")
		}
		result.WriteString(strings.Repeat("-", width+2))
	}
	result.WriteString("\n")
	
	// Print rows
	for _, row := range t.rows {
		t.printRow(&result, row)
	}
	
	return result.String()
}

func (t *Table) printRow(builder *strings.Builder, row []string) {
	for i, cell := range row {
		if i > 0 {
			builder.WriteString("|")
		}
		builder.WriteString(" ")
		builder.WriteString(cell)
		builder.WriteString(strings.Repeat(" ", t.widths[i]-len(cell)))
		builder.WriteString(" ")
	}
	builder.WriteString("\n")
}

// ProgressBar represents a progress bar
type ProgressBar struct {
	end      int64
	current  int64
	caption  string
	barlen   int
	start    time.Time
	mu       sync.Mutex
}

func NewProgressBar(end int64, caption string, barlen int) *ProgressBar {
	if barlen <= 0 {
		barlen = 50
	}
	
	return &ProgressBar{
		end:     end,
		caption: caption,
		barlen:  barlen,
		start:   time.Now(),
	}
}

func (p *ProgressBar) Update(current int64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.current = current
}

func (p *ProgressBar) String() string {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	if p.end <= 0 {
		return fmt.Sprintf("%s: 0%%", p.caption)
	}
	
	percent := float64(p.current) / float64(p.end) * 100
	filled := int(math.Round(float64(p.barlen) * float64(p.current) / float64(p.end)))
	
	bar := "[" + strings.Repeat("=", filled) + strings.Repeat(" ", p.barlen-filled) + "]"
	
	return fmt.Sprintf("%s: %s %.1f%% (%s/%s)", 
		p.caption, bar, percent, 
		NewSize(p.current).String(), NewSize(p.end).String())
}

// Paint provides text coloring functionality
type Paint struct {
	text   string
	colors []string
}

var colorCodes = map[string]string{
	"red":     "\033[31m",
	"green":   "\033[32m",
	"yellow":  "\033[33m",
	"blue":    "\033[34m",
	"magenta": "\033[35m",
	"cyan":    "\033[36m",
	"white":   "\033[37m",
	"gray":    "\033[90m",
	"reset":   "\033[0m",
}

func NewPaint(text string, colors []string) *Paint {
	return &Paint{text: text, colors: colors}
}

func (p *Paint) String() string {
	if len(p.colors) == 0 {
		return p.text
	}
	
	var result strings.Builder
	for _, color := range p.colors {
		if code, exists := colorCodes[color]; exists {
			result.WriteString(code)
		}
	}
	
	result.WriteString(p.text)
	result.WriteString(colorCodes["reset"])
	
	return result.String()
}

// Helper functions
func Red(text string) string {
	return NewPaint(text, []string{"red"}).String()
}

func Green(text string) string {
	return NewPaint(text, []string{"green"}).String()
}

func Yellow(text string) string {
	return NewPaint(text, []string{"yellow"}).String()
}

func Blue(text string) string {
	return NewPaint(text, []string{"blue"}).String()
}

func Cyan(text string) string {
	return NewPaint(text, []string{"cyan"}).String()
}

func Gray(text string) string {
	return NewPaint(text, []string{"gray"}).String()
}