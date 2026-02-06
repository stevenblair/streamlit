package main

import (
	st "github.com/streamlit/streamlit/go-backend"
)

func main() {
	st.Run(app)
}

func app() {
	// Title and header
	st.Title("🚀 Complex Streamlit Go Example")
	st.Header("Testing Multiple Element Types")

	// Text elements
	st.Subheader("📝 Text Elements")
	st.Write("This is plain text using st.Write()")
	st.Markdown("This is **bold** and *italic* markdown with `code`")

	// Code block
	st.Subheader("💻 Code Example")
	st.Code(`func HelloWorld() string {
    message := "Hello from Go!"
    return message
}`, "go")

	st.Code(`def python_example():
    """Python code example"""
    return "Hello from Python!"`, "python")

	// Divider
	st.Divider()

	// Data display
	st.Subheader("📊 Data Display")
	st.Code(`{
    "app": "Streamlit Go Backend",
    "version": "1.0.0",
    "features": ["Fast", "Type-safe", "Compiled"],
    "active": true
}`, "json")

	// Columns layout
	st.Divider()
	st.Subheader("📐 Multi-Column Layout")
	cols := st.Columns(3)
	cols[0].Write("**Column 1**")
	cols[0].Write("Left column content")
	cols[1].Write("**Column 2**")
	cols[1].Write("Middle column content")
	cols[2].Write("**Column 3**")
	cols[2].Write("Right column content")

	// Interactive elements
	st.Divider()
	st.Subheader("🎮 Interactive Elements")
	st.Button("Click Me!")
	st.Checkbox("Enable feature", false)
	st.TextInput("Enter your name", "")
	st.Slider("Select a value", 0, 100, 50)
	st.Selectbox("Choose option", []string{"Option 1", "Option 2", "Option 3"}, 0)

	// Status elements
	st.Divider()
	st.Subheader("⚡ Status Messages")
	st.Success("✅ Operation completed successfully!")
	st.Info("ℹ️ This is an informational message")
	st.Warning("⚠️ Warning: This is a test warning")
	st.Error("❌ Error: This is a test error")

	// Metrics
	st.Divider()
	st.Subheader("📈 Dashboard Metrics")
	metricCols := st.Columns(4)
	metricCols[0].Metric("Revenue", "$45K", "+12%")
	metricCols[1].Metric("Users", "1,234", "+5%")
	metricCols[2].Metric("Sales", "567", "-2%")
	metricCols[3].Metric("Rating", "4.8⭐", "+0.3")

	// More content
	st.Divider()
	st.Header("🎯 Advanced Features")

	st.Markdown(`
### Lists and Formatting

This demonstrates rich markdown support:

1. **Numbered lists**
2. *Styled text*
3. ` + "`Inline code`" + `

#### Unordered lists:
- Fast compilation
- Low memory usage
- Easy deployment
	`)

	st.Divider()
	st.Balloons()
	st.Write("🎈 **Balloons animation!** Try refreshing to see it again.")
}
