package main

import (
	st "github.com/stevenblair/streamlit/go-backend"
)

func main() {
	st.Run(app)
}

func app() {
	st.Title("🎈 Hello Streamlit (Go Backend)")

	st.Write(`
This app is being served by the **Streamlit Go Backend** with a native Go API!

The Go backend provides:
- Native Go API (no Python required!)
- Faster performance
- Type safety
- Compiled binary deployment
`)

	st.Header("Basic Elements")

	st.Write("Here's some basic text rendering.")

	st.Subheader("Subheader Example")

	st.Markdown(`
### Markdown Support

You can use **bold**, *italic*, and ` + "`code`" + ` formatting.

- List item 1
- List item 2
- List item 3
`)

	st.Code(`// Go code example
func hello() {
    fmt.Println("Hello from Go backend!")
    return 42
}`, "go")

	st.Divider()

	st.Success("✅ Success message")
	st.Info("ℹ️ Info message")
	st.Warning("⚠️ Warning message")
	st.Error("❌ Error message")

	st.Divider()

	// Simple metrics
	cols := st.Columns(3)
	cols[0].Metric("Temperature", "70°F", "1.2°F")
	cols[1].Metric("Wind", "9 mph", "-8%")
	cols[2].Metric("Humidity", "86%", "4%")

	st.Balloons()
}
