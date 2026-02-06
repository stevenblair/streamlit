package main

import (
	st "github.com/streamlit/streamlit/go-backend"
)

func main() {
	st.Run(app)
}

func app() {
	st.Title("🎮 Widget Examples")

	st.Write("Test various interactive widgets with the Go backend.")

	// Text input
	name := st.TextInput("Enter your name", "World")
	st.Writef("Hello, %s!", name)

	// Slider
	temperature := st.Slider("Temperature", -10, 40, 20)
	st.Writef("Current temperature: %.0f°C", temperature)

	// Selectbox
	option := st.Selectbox(
		"Choose your favorite color",
		[]string{"Red", "Green", "Blue", "Yellow"},
		0,
	)
	st.Writef("You selected: %s", option)

	// Checkbox
	agree := st.Checkbox("I agree to the terms and conditions", false)
	if agree {
		st.Write("Thank you for agreeing!")
	}

	// Button
	if st.Button("Click me!") {
		st.Balloons()
		st.Write("🎉 Button clicked!")
	}

	st.Divider()

	// Metrics in columns
	st.Subheader("Metrics")
	cols := st.Columns(3)
	cols[0].Metric("Users", "1,234", "+12%")
	cols[1].Metric("Revenue", "$5.6M", "+23%")
	cols[2].Metric("Uptime", "99.9%", "-0.1%")
}
