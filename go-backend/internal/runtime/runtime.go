package runtime

import (
	"context"
	"log"
	"sync"
)

// Runtime manages the execution of Streamlit apps and sessions
type Runtime struct {
	mu       sync.RWMutex
	sessions map[string]*Session
	ctx      context.Context
	cancel   context.CancelFunc
}

// New creates a new Runtime instance
func New() *Runtime {
	ctx, cancel := context.WithCancel(context.Background())

	return &Runtime{
		sessions: make(map[string]*Session),
		ctx:      ctx,
		cancel:   cancel,
	}
}

// RunApp executes a Streamlit Go app for a session
func (r *Runtime) RunApp(session SessionClient) {
	log.Printf("Running Streamlit Go app")

	// TODO: Implement actual app execution
	// This should:
	// 1. Execute the Go app function
	// 2. Collect st.* API calls
	// 3. Convert them to ForwardMsg protobuf messages
	// 4. Send messages to the session client

	// For now, just send a placeholder message
	go func() {
		// Simulate app execution
		log.Printf("App execution started")

		// TODO: Send actual ForwardMsgs with rendered elements
		// session.SendForwardMsg(protoMessage)

		log.Printf("App execution completed")
	}()
}

// Stop stops the runtime
func (r *Runtime) Stop() {
	log.Printf("Stopping runtime...")
	r.cancel()

	r.mu.Lock()
	defer r.mu.Unlock()

	// Close all sessions
	for _, session := range r.sessions {
		session.Close()
	}
	r.sessions = make(map[string]*Session)
}

// SessionClient represents a client that can receive ForwardMsg messages
type SessionClient interface {
	SendForwardMsg(data []byte) error
	Close()
}

// Session represents a user session
type Session struct {
	id     string
	client SessionClient
	state  map[string]interface{}
	mu     sync.RWMutex
}

// Close closes the session
func (s *Session) Close() {
	s.client.Close()
}
