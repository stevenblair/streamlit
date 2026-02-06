package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	st "github.com/streamlit/streamlit/go-backend"
)

func main() {
	fmt.Println("🧪 Testing Streamlit Go app...")
	fmt.Println("This will demonstrate that the app stays running.")
	fmt.Println()

	// Set a custom message handler to show we're connected
	log.SetPrefix("[TestApp] ")

	// Try to run the app
	done := make(chan error, 1)
	go func() {
		done <- st.Run(app)
	}()

	// Also listen for interrupt
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	fmt.Println("Starting app... (Press Ctrl+C to exit)")
	fmt.Println()

	// Wait a bit and show we're still running
	go func() {
		for i := 1; i <= 10; i++ {
			select {
			case <-done:
				return
			case <-time.After(1 * time.Second):
				fmt.Printf("⏱️  App still running... (%ds)\n", i)
			}
		}
		fmt.Println("\n✅ SUCCESS: App is staying alive (not exiting immediately!)")
		fmt.Println("Press Ctrl+C to exit")
	}()

	// Wait for completion or interrupt
	select {
	case err := <-done:
		if err != nil {
			fmt.Printf("\n❌ Error: %v\n", err)
			fmt.Println("\nNote: Make sure the Streamlit server is running:")
			fmt.Println("  ./bin/streamlit-go")
		} else {
			fmt.Println("\n✅ App completed successfully")
		}
	case <-interrupt:
		fmt.Println("\n\n👋 Received interrupt, exiting...")
	}
}

func app() {
	st.Title("🧪 Test App")
	st.Write("This app demonstrates that Go apps now stay running!")
	st.Success("If you can see this in your terminal logs, the app is working!")
}
