package server

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stevenblair/streamlit/go-backend/internal/runtime"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// TODO: Implement proper origin checking based on config
		return true
	},
}

// handleWebSocket handles WebSocket connections from the browser
func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	// Create a new session for this connection
	session := NewWebSocketSession(conn, s.runtime)
	session.Start()
}

// WebSocketSession represents a single WebSocket connection/session
type WebSocketSession struct {
	conn    *websocket.Conn
	runtime *runtime.Runtime
	send    chan []byte
	done    chan struct{}
	mu      sync.Mutex
}

// NewWebSocketSession creates a new WebSocket session
func NewWebSocketSession(conn *websocket.Conn, rt *runtime.Runtime) *WebSocketSession {
	return &WebSocketSession{
		conn:    conn,
		runtime: rt,
		send:    make(chan []byte, 256),
		done:    make(chan struct{}),
	}
}

// Start begins handling the WebSocket session
func (s *WebSocketSession) Start() {
	log.Printf("New WebSocket connection from %s", s.conn.RemoteAddr())

	// Start reader and writer goroutines
	go s.writePump()
	go s.readPump()

	// Send initial session info
	s.sendInitialMessages()
}

// readPump reads messages from the WebSocket connection
func (s *WebSocketSession) readPump() {
	defer func() {
		s.Close()
	}()

	s.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	s.conn.SetPongHandler(func(string) error {
		s.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		messageType, message, err := s.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket read error: %v", err)
			}
			break
		}

		if messageType == websocket.BinaryMessage {
			s.handleBackMsg(message)
		}
	}
}

// writePump writes messages to the WebSocket connection
func (s *WebSocketSession) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		s.Close()
	}()

	for {
		select {
		case message, ok := <-s.send:
			s.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				// Channel closed
				s.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := s.conn.WriteMessage(websocket.BinaryMessage, message); err != nil {
				log.Printf("WebSocket write error: %v", err)
				return
			}

		case <-ticker.C:
			s.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := s.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}

		case <-s.done:
			return
		}
	}
}

// handleBackMsg processes a BackMsg from the client
func (s *WebSocketSession) handleBackMsg(data []byte) {
	// TODO: Unmarshal BackMsg protobuf and handle different message types
	// For now, just log that we received a message
	log.Printf("Received BackMsg (%d bytes)", len(data))

	// Example message types to handle:
	// - rerun_script: Trigger script re-execution
	// - stop_script: Stop current script execution
	// - clear_cache: Clear cached data
	// - widget state updates: Update widget values
}

// SendForwardMsg sends a ForwardMsg to the client
func (s *WebSocketSession) SendForwardMsg(data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	select {
	case s.send <- data:
		return nil
	case <-s.done:
		return fmt.Errorf("session closed")
	default:
		return fmt.Errorf("send buffer full")
	}
}

// sendInitialMessages sends the initial session setup messages
func (s *WebSocketSession) sendInitialMessages() {
	// TODO: Send NewSession message with configuration
	// TODO: Send SessionStatus message
	// TODO: Trigger initial script run

	log.Printf("Sending initial session messages")

	// For now, just trigger an app run
	go func() {
		time.Sleep(100 * time.Millisecond)
		s.runtime.RunApp(s)
	}()
}

// Close closes the WebSocket session
func (s *WebSocketSession) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()

	select {
	case <-s.done:
		// Already closed
		return
	default:
		close(s.done)
		close(s.send)
		s.conn.Close()
		log.Printf("WebSocket connection closed")
	}
}
