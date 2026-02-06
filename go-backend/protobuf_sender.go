package streamlit

import (
	"log"

	"github.com/gorilla/websocket"
	pb "github.com/stevenblair/streamlit/go-backend/proto"
	"google.golang.org/protobuf/proto"
)

// sendAppRenderProtobuf sends the app render using protobuf ForwardMsg
func sendAppRenderProtobuf(conn *websocket.Conn) {
	// Get all elements
	mu.Lock()
	mainElements := make([]Element, len(elements))
	copy(mainElements, elements)
	mu.Unlock()

	sidebarMu.Lock()
	sidebarEls := make([]Element, len(sidebarElements))
	copy(sidebarEls, sidebarElements)
	sidebarMu.Unlock()

	// Convert elements to protobuf Delta messages
	// For now, we'll send each element as a separate Delta
	// In a full implementation, you'd batch these into a single ForwardMsg

	// Send NewSession message first (required by Streamlit protocol)
	sessionMsg := &pb.ForwardMsg{
		Type: &pb.ForwardMsg_NewSession{
			NewSession: &pb.NewSession{
				Name: "Streamlit Go App",
			},
		},
	}

	data, err := proto.Marshal(sessionMsg)
	if err != nil {
		log.Printf("Failed to marshal NewSession: %v", err)
		return
	}

	if err := conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
		log.Printf("Failed to send NewSession: %v", err)
		return
	}

	// Send main elements
	for _, elem := range mainElements {
		delta := convertElementToDelta(elem)
		if delta == nil {
			continue
		}

		msg := &pb.ForwardMsg{
			Type: &pb.ForwardMsg_Delta{
				Delta: delta,
			},
		}

		data, err := proto.Marshal(msg)
		if err != nil {
			log.Printf("Failed to marshal element: %v", err)
			continue
		}

		if err := conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
			log.Printf("Failed to send element: %v", err)
			return
		}
	}

	// TODO: Send sidebar elements similarly

	log.Printf("Sent %d protobuf elements", len(mainElements))
}

// convertElementToDelta converts a Go Element to a protobuf Delta
func convertElementToDelta(elem Element) *pb.Delta {
	// This is a simplified conversion - you'd need to implement
	// proper conversion for each element type

	switch e := elem.(type) {
	case *TextElement:
		return &pb.Delta{
			Type: &pb.Delta_NewElement{
				NewElement: &pb.Element{
					Type: &pb.Element_Text{
						Text: &pb.Text{
							Body: e.Content,
						},
					},
				},
			},
		}

	case *ButtonElement:
		return &pb.Delta{
			Type: &pb.Delta_NewElement{
				NewElement: &pb.Element{
					Type: &pb.Element_Button{
						Button: &pb.Button{
							Id:    e.Key,
							Label: e.Label,
						},
					},
				},
			},
		}

	// Add more element types here as needed
	default:
		log.Printf("Unknown element type: %T", elem)
		return nil
	}
}

// EnableProtobuf switches the backend to use protobuf messages instead of JSON
func EnableProtobuf() {
	// This would be set as a global flag to use sendAppRenderProtobuf
	// instead of sendAppRender
	log.Println("Protobuf messaging enabled")
}
