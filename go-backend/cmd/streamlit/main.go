package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/stevenblair/streamlit/go-backend/internal/server"
)

const (
	defaultPort    = 8501
	defaultHost    = "localhost"
	defaultBaseURL = ""
)

func main() {
	// Command line flags
	var (
		port      = flag.Int("port", defaultPort, "Port to run the server on")
		host      = flag.String("host", defaultHost, "Host to bind the server to")
		baseURL   = flag.String("base-url", defaultBaseURL, "Base URL for the app")
		serverURL = flag.String("server-url", "", "Full server URL (overrides host and port)")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Streamlit Go Backend - Native Go API for Streamlit apps\n\n")
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExample:\n")
		fmt.Fprintf(os.Stderr, "  %s --port 8502\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nNote: This server uses the native Go Streamlit API.\n")
		fmt.Fprintf(os.Stderr, "Build your app as a Go binary and it will connect to this server.\n")
	}

	flag.Parse()

	// Create server configuration
	config := &server.Config{
		Host:      *host,
		Port:      *port,
		BaseURL:   *baseURL,
		ServerURL: *serverURL,
	}

	// Create and start the server
	srv, err := server.New(config)
	if err != nil {
		log.Fatalf("Error creating server: %v", err)
	}

	log.Printf("Starting Streamlit Go backend...")
	log.Printf("Server: http://%s:%d", config.Host, config.Port)
	log.Printf("")
	log.Printf("You can now view your Streamlit app in your browser.")
	log.Printf("")
	log.Printf("  Local URL: http://%s:%d", config.Host, config.Port)
	if config.Host != "0.0.0.0" && config.Host != "localhost" {
		log.Printf("  Network URL: http://%s:%d", config.Host, config.Port)
	}
	log.Printf("")
	log.Printf("Run your Streamlit Go app binary to connect.")
	log.Printf("")

	if err := srv.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
