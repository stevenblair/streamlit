package main

import (
	"fmt"
	"time"

	st "github.com/stevenblair/streamlit/go-backend"
)

func main() {
	st.Run(app)
}

func app() {
	// Page header
	st.Title("🎨 Complete Widget Showcase")
	st.Markdown("**Demonstration of all available Streamlit widgets in Go**")
	st.Divider()

	// ============================================================================
	// TEXT ELEMENTS SECTION
	// ============================================================================
	st.Header("📝 Text Elements")
	
	st.Write("This is `st.Write()` - displays any text or data")
	st.Writef("This is `st.Writef()` - formatted text with %s", "variables")
	st.Title("This is st.Title()")
	st.Header("This is st.Header()")
	st.Subheader("This is st.Subheader()")
	
	st.Markdown(`
### st.Markdown() - Rich Text Formatting

Supports **bold**, *italic*, ` + "`code`" + `, and more:
- Bullet points
- Nested lists
  - Sub-items
- [Links](https://streamlit.io)

**Features:**
1. Numbered lists
2. Code blocks
3. Tables (coming soon)
`)
	
	st.Code(`func main() {
    fmt.Println("Hello, World!")
    st.Title("Code Display")
}`, "go")

	st.Divider()

	// ============================================================================
	// STATUS/ALERT ELEMENTS SECTION
	// ============================================================================
	st.Header("✅ Status & Alert Messages")
	
	st.Success("✅ Success! Operation completed successfully")
	st.Info("ℹ️ Info: This is an informational message")
	st.Warning("⚠️ Warning: Please review this carefully")
	st.Error("❌ Error: Something went wrong")
	
	st.Divider()

	// ============================================================================
	// INPUT WIDGETS SECTION
	// ============================================================================
	st.Header("🎛️ Input Widgets")
	st.Subheader("Text Input")
	
	name := st.TextInput("Enter your name", "")
	if name != "" {
		st.Success(fmt.Sprintf("👋 Hello, %s!", name))
	}

	st.Subheader("Number Slider")
	
	temperature := st.Slider("Temperature (°C)", -20.0, 50.0, 20.0)
	if temperature < 0 {
		st.Info(fmt.Sprintf("🥶 It's freezing at %.1f°C", temperature))
	} else if temperature > 30 {
		st.Warning(fmt.Sprintf("🔥 It's hot at %.1f°C", temperature))
	} else {
		st.Success(fmt.Sprintf("😊 Comfortable at %.1f°C", temperature))
	}

	st.Subheader("Checkbox")
	
	agree := st.Checkbox("I agree to the terms and conditions", false)
	if agree {
		st.Success("✅ Thank you for accepting!")
	} else {
		st.Warning("⚠️ Please accept the terms to continue")
	}

	notifications := st.Checkbox("Enable notifications", true)
	darkMode := st.Checkbox("Dark mode", false)
	
	st.Write(fmt.Sprintf("Settings: Notifications=%v, Dark Mode=%v", notifications, darkMode))

	st.Subheader("Selectbox (Dropdown)")
	
	color := st.Selectbox(
		"Choose your favorite color",
		[]string{"🔴 Red", "🟢 Green", "🔵 Blue", "🟡 Yellow", "🟣 Purple"},
		0,
	)
	st.Write(fmt.Sprintf("You selected: %s", color))

	language := st.Selectbox(
		"Programming Language",
		[]string{"Go", "Python", "JavaScript", "Rust", "TypeScript"},
		0,
	)
	st.Code(getCodeSample(language), getLanguageCode(language))

	st.Divider()

	// ============================================================================
	// BUTTON WIDGETS SECTION
	// ============================================================================
	st.Header("🔘 Buttons & Actions")
	
	st.Write("Click the buttons below to trigger actions:")
	
	if st.Button("🎉 Celebrate") {
		st.Balloons()
		st.Success("🎊 Celebration time!")
	}
	
	if st.Button("❄️ Snow") {
		st.Snow()
		st.Info("❄️ Let it snow!")
	}
	
	if st.Button("ℹ️ Show Info") {
		st.Info("This is additional information that appears when clicked!")
	}

	if st.Button("🔄 Refresh Data") {
		st.Success(fmt.Sprintf("Data refreshed at %s", time.Now().Format("15:04:05")))
	}

	st.Divider()

	// ============================================================================
	// LAYOUT ELEMENTS SECTION
	// ============================================================================
	st.Header("📐 Layout & Containers")
	
	st.Subheader("Multi-Column Layout")
	st.Write("Organize content side-by-side using columns:")
	
	cols := st.Columns(3)
	
	cols[0].Write("### Column 1")
	cols[0].Write("🎯 **Left Column**")
	cols[0].Write("Content in the first column")
	
	cols[1].Write("### Column 2")
	cols[1].Write("⚡ **Middle Column**")
	cols[1].Write("Content in the second column")
	
	cols[2].Write("### Column 3")
	cols[2].Write("🚀 **Right Column**")
	cols[2].Write("Content in the third column")

	st.Divider()

	// ============================================================================
	// METRICS & DASHBOARDS SECTION
	// ============================================================================
	st.Header("📊 Metrics & KPIs")
	st.Write("Display key performance indicators with delta changes:")
	
	metricCols := st.Columns(4)
	
	metricCols[0].Metric("Revenue", "$45,231", "+12.5%")
	metricCols[1].Metric("Active Users", "1,234", "+8.2%")
	metricCols[2].Metric("Response Time", "245ms", "-15.3%")
	metricCols[3].Metric("Success Rate", "99.2%", "+0.3%")

	st.Divider()
	
	// Second row of metrics
	st.Subheader("Real-time Monitoring")
	
	monitorCols := st.Columns(3)
	
	monitorCols[0].Metric("CPU Usage", "45%", "+5%")
	monitorCols[1].Metric("Memory", "2.3GB", "-0.1GB")
	monitorCols[2].Metric("Disk I/O", "125MB/s", "+25MB/s")

	st.Divider()

	// ============================================================================
	// DATA VISUALIZATION SECTION
	// ============================================================================
	st.Header("📈 Data Visualization")
	st.Write("Create interactive charts from your data:")
	
	st.Subheader("Line Chart Example")
	
	// Generate sample data based on slider
	numPoints := int(st.Slider("Number of data points", 10.0, 1000.0, 100.0))
	amplitude := st.Slider("Wave amplitude", 1.0, 10.0, 5.0)
	
	data := generateSineWave(numPoints, amplitude)
	st.LineChart(data)
	
	st.Info(fmt.Sprintf("Displaying %d data points with amplitude %.1f", numPoints, amplitude))

	st.Divider()

	// ============================================================================
	// SIDEBAR SECTION
	// ============================================================================
	sidebar := st.Sidebar()
	
	sidebar.Title("⚙️ Sidebar")
	sidebar.Write("This is the sidebar - perfect for controls and settings")
	
	sidebar.Divider()
	
	sidebar.Header("Settings")
	
	theme := sidebar.Selectbox(
		"Theme",
		[]string{"Light", "Dark", "Auto"},
		0,
	)
	
	fontSize := sidebar.Slider("Font Size", 10.0, 24.0, 14.0)
	
	showAdvanced := sidebar.Checkbox("Show advanced options", false)
	
	if showAdvanced {
		sidebar.Subheader("Advanced Settings")
		sidebar.Write("• Cache duration")
		sidebar.Write("• Debug mode")
		sidebar.Write("• API endpoint")
	}
	
	sidebar.Divider()
	
	sidebar.Header("About")
	sidebar.Write("**Streamlit Go Backend**")
	sidebar.Write("Version: 1.0.0")
	sidebar.Write("Get started:")
	sidebar.Write("go get github.com/...")
	
	if sidebar.Button("❓ Help") {
		st.Info("📚 Documentation available at streamlit.io")
	}

	// Display sidebar selections in main area
	st.Divider()
	st.Header("📋 Current Configuration")
	
	configCols := st.Columns(2)
	configCols[0].Write(fmt.Sprintf("**Theme:** %s", theme))
	configCols[0].Write(fmt.Sprintf("**Font Size:** %.0fpt", fontSize))
	configCols[1].Write(fmt.Sprintf("**Name:** %s", ifEmpty(name, "Not set")))
	configCols[1].Write(fmt.Sprintf("**Language:** %s", language))

	// ============================================================================
	// FOOTER
	// ============================================================================
	st.Divider()
	st.Markdown("---")
	st.Markdown(`
### 🎯 Summary of Available Widgets

**Text Display:**
- ` + "`st.Write()`, `st.Writef()`, `st.Title()`, `st.Header()`, `st.Subheader()`" + `
- ` + "`st.Markdown()`, `st.Code()`, `st.Divider()`" + `

**Status Messages:**
- ` + "`st.Success()`, `st.Info()`, `st.Warning()`, `st.Error()`" + `

**Input Widgets:**
- ` + "`st.TextInput()`, `st.Slider()`, `st.Checkbox()`, `st.Selectbox()`" + `
- ` + "`st.Button()`" + `

**Layout:**
- ` + "`st.Columns()`, `st.Metric()`, `st.Sidebar()`" + `

**Visualizations:**
- ` + "`st.LineChart()`" + `

**Animations:**
- ` + "`st.Balloons()`, `st.Snow()`" + `
`)
	
	st.Info("💡 All widgets are fully functional and support state management!")
}

// Helper functions

func generateSineWave(points int, amplitude float64) []float64 {
	data := make([]float64, points)
	for i := 0; i < points; i++ {
		x := float64(i) / float64(points) * 2.0 * 3.14159 * 3 // 3 full waves
		data[i] = amplitude * (0.5 + 0.5*sinApprox(x))
	}
	return data
}

// Fast sine approximation using Bhaskara I's formula
func sinApprox(x float64) float64 {
	// Normalize to [-π, π]
	for x > 3.14159 {
		x -= 2 * 3.14159
	}
	for x < -3.14159 {
		x += 2 * 3.14159
	}
	
	// Bhaskara I's sine approximation
	if x < 0 {
		return -sinApprox(-x)
	}
	return (16*x*(3.14159-x)) / (5*3.14159*3.14159 - 4*x*(3.14159-x))
}

func getCodeSample(language string) string {
	samples := map[string]string{
		"Go": `func main() {
    fmt.Println("Hello from Go!")
    st.Title("Streamlit + Go")
}`,
		"Python": `def main():
    print("Hello from Python!")
    st.title("Streamlit + Python")`,
		"JavaScript": `function main() {
    console.log("Hello from JavaScript!");
    st.title("Streamlit + JavaScript");
}`,
		"Rust": `fn main() {
    println!("Hello from Rust!");
    st::title("Streamlit + Rust");
}`,
		"TypeScript": `function main(): void {
    console.log("Hello from TypeScript!");
    st.title("Streamlit + TypeScript");
}`,
	}
	
	if sample, ok := samples[language]; ok {
		return sample
	}
	return "// Code sample not available"
}

func getLanguageCode(language string) string {
	langCodes := map[string]string{
		"Go":         "go",
		"Python":     "python",
		"JavaScript": "javascript",
		"Rust":       "rust",
		"TypeScript": "typescript",
	}
	
	if code, ok := langCodes[language]; ok {
		return code
	}
	return ""
}

func ifEmpty(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}
