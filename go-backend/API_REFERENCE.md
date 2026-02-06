# Streamlit Go API Reference

Complete API reference for writing Streamlit apps in Go.

## Import

```go
import st "github.com/stevenblair/streamlit/go-backend"
```

## App Lifecycle

### st.Run(appFunc)

Executes your Streamlit app function.

```go
func main() {
    st.Run(myApp)
}

func myApp() {
    st.Title("My App")
    // Your app code here
}
```

**Parameters:**
- `appFunc func()` - Your app function

**Returns:** `error` - Error if app fails

## Text Elements

### st.Write(args...)

Display anything. Accepts any number of arguments.

```go
st.Write("Hello, World!")
st.Write("Temperature:", 72, "°F")
```

**Parameters:**
- `args ...interface{}` - Values to display

### st.Writef(format, args...)

Display formatted text (like fmt.Sprintf).

```go
st.Writef("Hello, %s! You are %d years old.", name, age)
```

**Parameters:**
- `format string` - Format string
- `args ...interface{}` - Values for formatting

### st.Title(text)

Display a title (largest heading).

```go
st.Title("My Application")
st.Title("🎈 Welcome!")
```

**Parameters:**
- `text string` - Title text

### st.Header(text)

Display a header (medium heading).

```go
st.Header("Section 1")
```

**Parameters:**
- `text string` - Header text

### st.Subheader(text)

Display a subheader (smaller heading).

```go
st.Subheader("Data Analysis")
```

**Parameters:**
- `text string` - Subheader text

### st.Markdown(content)

Display markdown content.

```go
st.Markdown(`
## Features
- **Bold** text
- *Italic* text
- ` + "`code`" + `
`)
```

**Parameters:**
- `content string` - Markdown text

### st.Code(code, language)

Display a code block with syntax highlighting.

```go
st.Code(`func hello() {
    fmt.Println("Hello!")
}`, "go")
```

**Parameters:**
- `code string` - Code content
- `language string` - Language ("go", "python", "javascript", etc.)

### st.Divider()

Display a horizontal divider line.

```go
st.Divider()
```

## Status Messages

### st.Success(message)

Display a success message (green).

```go
st.Success("✅ Operation completed successfully!")
```

**Parameters:**
- `message string` - Success message

### st.Info(message)

Display an info message (blue).

```go
st.Info("ℹ️ Did you know...")
```

**Parameters:**
- `message string` - Info message

### st.Warning(message)

Display a warning message (yellow).

```go
st.Warning("⚠️ Please review before proceeding")
```

**Parameters:**
- `message string` - Warning message

### st.Error(message)

Display an error message (red).

```go
st.Error("❌ An error occurred")
```

**Parameters:**
- `message string` - Error message

## Input Widgets

### st.Button(label) bool

Display a button.

```go
if st.Button("Click me") {
    st.Write("Button was clicked!")
}
```

**Parameters:**
- `label string` - Button label

**Returns:** `bool` - True if clicked (this run)

### st.TextInput(label, defaultValue) string

Display a text input field.

```go
name := st.TextInput("Enter your name", "")
st.Writef("Hello, %s!", name)
```

**Parameters:**
- `label string` - Input label
- `defaultValue string` - Default value

**Returns:** `string` - Current input value

### st.Slider(label, min, max, defaultValue) int

Display a slider widget.

```go
age := st.Slider("Select your age", 0, 100, 25)
st.Writef("You are %d years old", age)
```

**Parameters:**
- `label string` - Slider label
- `min int` - Minimum value
- `max int` - Maximum value
- `defaultValue int` - Default value

**Returns:** `int` - Current slider value

### st.Checkbox(label, defaultValue) bool

Display a checkbox.

```go
agree := st.Checkbox("I agree to terms", false)
if agree {
    st.Success("Thank you!")
}
```

**Parameters:**
- `label string` - Checkbox label
- `defaultValue bool` - Default checked state

**Returns:** `bool` - Current checked state

### st.Selectbox(label, options, defaultIndex) string

Display a selectbox (dropdown).

```go
color := st.Selectbox(
    "Choose a color",
    []string{"Red", "Green", "Blue"},
    0,
)
st.Writef("You selected: %s", color)
```

**Parameters:**
- `label string` - Selectbox label
- `options []string` - List of options
- `defaultIndex int` - Index of default option (0-based)

**Returns:** `string` - Selected option

## Layouts

### st.Columns(count) []*Column

Create a multi-column layout.

```go
cols := st.Columns(3)
cols[0].Write("Column 1")
cols[1].Write("Column 2")
cols[2].Write("Column 3")
```

**Parameters:**
- `count int` - Number of columns

**Returns:** `[]*Column` - Array of column objects

### Column Methods

Each column object returned by `st.Columns()` has these methods:

#### col.Write(args...)

Write to the column.

```go
col.Write("Content in column")
```

#### col.Metric(label, value, delta)

Display a metric in the column.

```go
col.Metric("Temperature", "70°F", "+2°F")
```

## Metrics

### st.Metric(label, value, delta)

Display a metric with label, value, and optional delta.

```go
st.Metric("Revenue", "$1.2M", "+12%")
st.Metric("Users", "1,234", "-5%")
```

**Parameters:**
- `label string` - Metric label
- `value string` - Metric value
- `delta string` - Change indicator (can be empty)

## Animations

### st.Balloons()

Trigger balloons animation.

```go
if st.Button("Celebrate") {
    st.Balloons()
}
```

### st.Snow()

Trigger snow animation.

```go
if st.Button("Let it snow") {
    st.Snow()
}
```

## Complete Example

```go
package main

import st "github.com/stevenblair/streamlit/go-backend"

func main() {
    st.Run(app)
}

func app() {
    // Title
    st.Title("🎯 Complete API Demo")

    // Text elements
    st.Header("Welcome")
    st.Write("This demonstrates all API features")
    st.Markdown("With **markdown** support")

    st.Code(`func main() {
    fmt.Println("Hello!")
}`, "go")

    st.Divider()

    // Status messages
    st.Success("Success message")
    st.Info("Info message")
    st.Warning("Warning message")
    st.Error("Error message")

    st.Divider()

    // Input widgets
    name := st.TextInput("Name", "User")
    age := st.Slider("Age", 0, 100, 25)
    color := st.Selectbox("Color", []string{"Red", "Green", "Blue"}, 0)
    agree := st.Checkbox("Agree", false)

    if st.Button("Submit") {
        st.Writef("Name: %s, Age: %d, Color: %s, Agreed: %v",
            name, age, color, agree)
        st.Balloons()
    }

    st.Divider()

    // Metrics in columns
    cols := st.Columns(3)
    cols[0].Metric("Users", "1,234", "+12%")
    cols[1].Metric("Revenue", "$5.6M", "+23%")
    cols[2].Metric("Uptime", "99.9%", "-0.1%")
}
```

## Type Signatures

For reference, here are the complete type signatures:

```go
// Lifecycle
func Run(appFunc func()) error

// Text
func Write(args ...interface{})
func Writef(format string, args ...interface{})
func Title(text string)
func Header(text string)
func Subheader(text string)
func Markdown(content string)
func Code(code string, language string)
func Divider()

// Status
func Success(message string)
func Info(message string)
func Warning(message string)
func Error(message string)

// Widgets
func Button(label string) bool
func TextInput(label string, defaultValue string) string
func Slider(label string, min, max, defaultValue int) int
func Checkbox(label string, defaultValue bool) bool
func Selectbox(label string, options []string, defaultIndex int) string

// Layout
func Columns(count int) []*Column
func Metric(label, value, delta string)

// Animations
func Balloons()
func Snow()

// Column type
type Column struct {
    Write(args ...interface{})
    Metric(label, value, delta string)
}
```

## State Management (Coming Soon)

Future API will include:

```go
// Session state (planned)
st.SessionState.Get(key string) interface{}
st.SessionState.Set(key string, value interface{})

// Caching (planned)
st.Cache(func() interface{}) interface{}
```

## Notes

- **Widget keys**: Currently auto-generated from labels
- **Reruns**: Triggered automatically on widget interaction
- **State**: Widget state persists during session (planned)
- **Type safety**: Go's type system catches errors at compile time

## Getting Help

- Check [examples/](../examples/) for working code
- Read [QUICKSTART.md](QUICKSTART.md) for setup
- File issues with `go-backend` label

---

**Version**: 0.1.0 (Experimental)
