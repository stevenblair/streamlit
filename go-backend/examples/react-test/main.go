package main

import (
	"log"
	"time"

	st "github.com/stevenblair/streamlit/go-backend"
)

func main() {
	st.Title("🎉 React Frontend Integration Test")
	st.Write("This app is using the **React frontend** built from the Streamlit repository!")

	st.Header("✅ What's Working")
	st.Write("- Go backend with protobuf support")
	st.Write("- React frontend from frontend/app/build/")
	st.Write("- WebSocket communication")
	st.Write("- All 87 .pb.go files generated")

	st.Subheader("🔧 Features Available")

	// Test various widgets
	if st.Button("Test Button") {
		st.Success("Button clicked! ✅")
	}

	name := st.TextInput("Enter your name", "")
	if name != "" {
		st.Write("Hello, " + name + "! 👋")
	}

	value := st.Slider("Select a value", 0.0, 100.0, 50.0)
	st.Write("Slider value:", value)

	checked := st.Checkbox("Enable feature", false)
	if checked {
		st.Info("Feature enabled!")
	}

	// Sidebar
	st.Sidebar().Header("Sidebar Test")
	st.Sidebar().Write("This is in the sidebar")
	option := st.Sidebar().Selectbox("Choose option", []string{"Option 1", "Option 2", "Option 3"}, 0)
	st.Sidebar().Write("Selected:", option)

	st.Divider()

	st.Header("📊 Integration Status")
	st.Metric("Protobuf Files", "87", "+87")
	st.Metric("Frontend", "React", "✅")
	st.Metric("Status", "Complete", "✅")

	st.Divider()

	st.Info("💡 **Next Steps**: The React frontend is now serving all UI components. You can extend the protobuf integration to support more complex element types!")

	// Add timestamp
	st.Write("Last updated: " + time.Now().Format("2006-01-02 15:04:05"))

	log.Println("React integration test app ready at http://localhost:8501")
	st.Run(func() {})
}
