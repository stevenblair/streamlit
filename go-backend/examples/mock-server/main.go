package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade error: %v", err)
		return
	}
	defer conn.Close()

	log.Printf("✅ Client connected from %s", conn.RemoteAddr())

	// Read messages from client
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("❌ Connection error: %v", err)
			}
			break
		}

		log.Printf("📩 Received message (%d bytes): %s", len(message), string(message))

		// Echo back
		if err := conn.WriteMessage(messageType, []byte("OK")); err != nil {
			log.Printf("Write error: %v", err)
			break
		}
	}

	log.Printf("👋 Client disconnected")
}

func main() {
	fmt.Println("🧪 Starting mock Streamlit server...")
	fmt.Println("This server accepts WebSocket connections to test that apps stay alive.")
	fmt.Println()

	http.HandleFunc("/_stcore/stream", handleWebSocket)
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"ok"}`))
	})

	addr := ":8501"
	fmt.Printf("Mock server listening on http://localhost%s\n", addr)
	fmt.Println("WebSocket endpoint: ws://localhost:8501/_stcore/stream")
	fmt.Println()
	fmt.Println("Now run an example app in another terminal:")
	fmt.Println("  .\\bin\\hello.exe")
	fmt.Println("  .\\bin\\widgets.exe")
	fmt.Println()
	fmt.Println("Press Ctrl+C to stop")
	fmt.Println()

	server := &http.Server{
		Addr:         addr,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
