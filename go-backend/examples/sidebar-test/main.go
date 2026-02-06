package main

import (
	"fmt"
	st "github.com/streamlit/streamlit/go-backend"
)

func main() {
	st.Run(app)
}

func app() {
	st.Title("Sidebar Test")
	st.Write("Main area content")
	
	// Test sidebar
	sidebar := st.Sidebar()
	sidebar.Title("Sidebar Title")
	sidebar.Write("Sidebar content")
	sidebar.Header("Settings")
	
	value := sidebar.Slider("Test Slider", 0.0, 100.0, 50.0)
	sidebar.Write(fmt.Sprintf("Slider value: %.0f", value))
	
	name := sidebar.TextInput("Your name", "")
	if name != "" {
		sidebar.Write(fmt.Sprintf("Hello, %s!", name))
	}
	
	checked := sidebar.Checkbox("Enable feature", false)
	sidebar.Write(fmt.Sprintf("Checkbox: %v", checked))
	
	// Back to main area
	st.Write("Back in main area")
	st.Write(fmt.Sprintf("Sidebar slider value: %.0f", value))
}
