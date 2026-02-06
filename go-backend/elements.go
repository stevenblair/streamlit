package streamlit

import (
	"encoding/json"
	"fmt"
	"sync"
)

// TextElement represents a text element
type TextElement struct {
	Content string
}

func (t *TextElement) Render() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":    "text",
		"content": t.Content,
	})
}

// HeadingElement represents a heading (title, header, subheader)
type HeadingElement struct {
	Content string
	Level   int // 1 = title, 2 = header, 3 = subheader
}

func (h *HeadingElement) Render() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":    "heading",
		"content": h.Content,
		"level":   h.Level,
	})
}

// MarkdownElement represents markdown content
type MarkdownElement struct {
	Content string
}

func (m *MarkdownElement) Render() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":    "markdown",
		"content": m.Content,
	})
}

// CodeElement represents a code block
type CodeElement struct {
	Content  string
	Language string
}

func (c *CodeElement) Render() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":     "code",
		"content":  c.Content,
		"language": c.Language,
	})
}

// DividerElement represents a horizontal divider
type DividerElement struct{}

func (d *DividerElement) Render() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type": "divider",
	})
}

// AlertElement represents a status message (success, info, warning, error)
type AlertElement struct {
	Message string
	Type    string // "success", "info", "warning", "error"
}

func (a *AlertElement) Render() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":   "alert",
		"body":   a.Message,
		"format": a.Type,
	})
}

// ButtonElement represents a button widget
type ButtonElement struct {
	Label string
	Key   string
}

func (b *ButtonElement) Render() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":  "button",
		"label": b.Label,
		"key":   b.Key,
	})
}

// TextInputElement represents a text input widget
type TextInputElement struct {
	Label        string
	DefaultValue string
	Key          string
}

func (t *TextInputElement) Render() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":         "text_input",
		"label":        t.Label,
		"defaultValue": t.DefaultValue,
		"key":          t.Key,
	})
}

// SliderElement represents a slider widget
type SliderElement struct {
	Label        string
	Min          int
	Max          int
	DefaultValue int
	Key          string
}

func (s *SliderElement) Render() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":         "slider",
		"label":        s.Label,
		"min":          s.Min,
		"max":          s.Max,
		"defaultValue": s.DefaultValue,
		"key":          s.Key,
	})
}

// CheckboxElement represents a checkbox widget
type CheckboxElement struct {
	Label        string
	DefaultValue bool
	Key          string
}

func (c *CheckboxElement) Render() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":         "checkbox",
		"label":        c.Label,
		"defaultValue": c.DefaultValue,
		"key":          c.Key,
	})
}

// SelectboxElement represents a selectbox widget
type SelectboxElement struct {
	Label        string
	Options      []string
	DefaultIndex int
	Key          string
}

func (s *SelectboxElement) Render() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":         "selectbox",
		"label":        s.Label,
		"options":      s.Options,
		"defaultIndex": s.DefaultIndex,
		"key":          s.Key,
	})
}

// MetricElement represents a metric display
type MetricElement struct {
	Label string
	Value string
	Delta string
}

func (m *MetricElement) Render() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":  "metric",
		"label": m.Label,
		"value": m.Value,
		"delta": m.Delta,
	})
}

// Column represents a single column in a multi-column layout
type Column struct {
	mu       sync.Mutex
	index    int
	elements []Element
}

// Write writes text to the column
func (c *Column) Write(args ...interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.elements = append(c.elements, &TextElement{Content: fmt.Sprint(args...)})
}

// Metric displays a metric in the column
func (c *Column) Metric(label string, value string, delta string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.elements = append(c.elements, &MetricElement{Label: label, Value: value, Delta: delta})
}

// Slider creates a slider widget in the column
func (c *Column) Slider(label string, min, max, defaultValue float64) float64 {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := fmt.Sprintf("slider_%s_%d", label, c.index)

	// Check for existing state
	widgetStateMu.RLock()
	if val, exists := widgetState[key]; exists {
		widgetStateMu.RUnlock()
		if floatVal, ok := val.(float64); ok {
			c.elements = append(c.elements, &SliderElement{
				Label:        label,
				Min:          int(min),
				Max:          int(max),
				DefaultValue: int(floatVal),
				Key:          key,
			})
			return floatVal
		}
	}
	widgetStateMu.RUnlock()

	c.elements = append(c.elements, &SliderElement{
		Label:        label,
		Min:          int(min),
		Max:          int(max),
		DefaultValue: int(defaultValue),
		Key:          key,
	})
	return defaultValue
}

// ColumnsElement represents a multi-column layout
type ColumnsElement struct {
	Columns []*Column
}

func (c *ColumnsElement) Render() ([]byte, error) {
	cols := make([]map[string]interface{}, len(c.Columns))
	for i, col := range c.Columns {
		col.mu.Lock()
		colElements := make([]map[string]interface{}, len(col.elements))
		for j, elem := range col.elements {
			data, err := elem.Render()
			if err != nil {
				col.mu.Unlock()
				return nil, err
			}
			var elemMap map[string]interface{}
			if err := json.Unmarshal(data, &elemMap); err != nil {
				col.mu.Unlock()
				return nil, err
			}
			colElements[j] = elemMap
		}
		col.mu.Unlock()
		cols[i] = map[string]interface{}{
			"index":    i,
			"elements": colElements,
		}
	}
	return json.Marshal(map[string]interface{}{
		"type":    "columns",
		"count":   len(c.Columns),
		"columns": cols,
	})
}

// BalloonsElement triggers balloons animation
type BalloonsElement struct{}

func (b *BalloonsElement) Render() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type": "balloons",
	})
}

// SnowElement triggers snow animation
type SnowElement struct{}

func (s *SnowElement) Render() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type": "snow",
	})
}

// LineChartElement represents a line chart
type LineChartElement struct {
	Data []float64
}

func (l *LineChartElement) Render() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type": "line_chart",
		"data": l.Data,
	})
}
