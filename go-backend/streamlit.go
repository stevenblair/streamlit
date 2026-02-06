package streamlit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"time"

	"github.com/apache/arrow/go/v18/arrow"
	"github.com/apache/arrow/go/v18/arrow/array"
	"github.com/apache/arrow/go/v18/arrow/ipc"
	"github.com/apache/arrow/go/v18/arrow/memory"
	"github.com/gorilla/websocket"
	pb "github.com/streamlit/streamlit/go-backend/proto"
	"google.golang.org/protobuf/proto"
)

var (
	elements        []Element
	sidebarElements []Element
	mu              sync.Mutex
	sidebarMu       sync.Mutex
	appFunc         func()
	widgetState     map[string]interface{}
	widgetStateMu   sync.RWMutex
	inSidebar       bool
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Element represents a Streamlit UI element
type Element interface {
	Render() ([]byte, error)
}

// Initialize resets the element list for a new render
func Initialize() {
	mu.Lock()
	defer mu.Unlock()
	elements = nil

	sidebarMu.Lock()
	sidebarElements = nil
	sidebarMu.Unlock()

	inSidebar = false

	// Initialize widget state if needed
	widgetStateMu.Lock()
	if widgetState == nil {
		widgetState = make(map[string]interface{})
	}
	widgetStateMu.Unlock()
}

// addElement adds an element to the current render
func addElement(elem Element) {
	if inSidebar {
		sidebarMu.Lock()
		defer sidebarMu.Unlock()
		sidebarElements = append(sidebarElements, elem)
	} else {
		mu.Lock()
		defer mu.Unlock()
		elements = append(elements, elem)
	}
}

// handleWebSocket handles WebSocket connections from the frontend
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	log.Println("Client connected")

	// Send initial render
	sendAppRender(conn)

	// Keep connection alive and handle messages from client
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Parse BackMsg protobuf
		var backMsg pb.BackMsg
		if err := proto.Unmarshal(message, &backMsg); err != nil {
			log.Printf("Failed to parse BackMsg: %v", err)
			continue
		}

		// Handle RerunScript message which contains widget states
		if rerunScript := backMsg.GetRerunScript(); rerunScript != nil {
			if widgetStates := rerunScript.GetWidgetStates(); widgetStates != nil {
				log.Printf("Received %d widget states", len(widgetStates.GetWidgets()))
				widgetStateMu.Lock()
				for _, widget := range widgetStates.GetWidgets() {
					widgetId := widget.GetId()

					// Extract the typed value
					switch v := widget.Value.(type) {
					case *pb.WidgetState_DoubleValue:
						widgetState[widgetId] = v.DoubleValue
						log.Printf("Widget update: %s = %f (double)", widgetId, v.DoubleValue)
					case *pb.WidgetState_IntValue:
						widgetState[widgetId] = v.IntValue
						log.Printf("Widget update: %s = %d (int)", widgetId, v.IntValue)
					case *pb.WidgetState_BoolValue:
						widgetState[widgetId] = v.BoolValue
						log.Printf("Widget update: %s = %v (bool)", widgetId, v.BoolValue)
					case *pb.WidgetState_StringValue:
						widgetState[widgetId] = v.StringValue
						log.Printf("Widget update: %s = %s (string)", widgetId, v.StringValue)
					case *pb.WidgetState_TriggerValue:
						widgetState[widgetId] = true
						log.Printf("Widget update: %s = triggered", widgetId)
					case *pb.WidgetState_DoubleArrayValue:
						if len(v.DoubleArrayValue.Data) > 0 {
							widgetState[widgetId] = v.DoubleArrayValue.Data[0]
							log.Printf("Widget update: %s = %f (double array)", widgetId, v.DoubleArrayValue.Data[0])
						}
					case *pb.WidgetState_IntArrayValue:
						if len(v.IntArrayValue.Data) > 0 {
							widgetState[widgetId] = v.IntArrayValue.Data[0]
							log.Printf("Widget update: %s = %d (int array)", widgetId, v.IntArrayValue.Data[0])
						}
					default:
						log.Printf("Widget update: %s = unknown type %T", widgetId, widget.Value)
					}
				}
				widgetStateMu.Unlock()
			}
		}

		// Send updated render
		sendAppRender(conn)
	}

	log.Println("Client disconnected")
}

// sendAppRender runs the app and sends all elements as protobuf messages
func sendAppRender(conn *websocket.Conn) {
	// Run the app to collect elements
	Initialize()
	appFunc()

	// Send NewSession message first with all required fields
	sessionMsg := &pb.ForwardMsg{
		Type: &pb.ForwardMsg_NewSession{
			NewSession: &pb.NewSession{
				Initialize: &pb.Initialize{
					UserInfo: &pb.UserInfo{
						InstallationId:   "go-backend-id",
						InstallationIdV3: "go-backend-id-v3",
						InstallationIdV4: "go-backend-id-v4",
					},
					EnvironmentInfo: &pb.EnvironmentInfo{
						StreamlitVersion: "1.0.0-go",
						PythonVersion:    "Go Backend",
					},
					SessionStatus: &pb.SessionStatus{
						RunOnSave:       false,
						ScriptIsRunning: false,
					},
					SessionId: "go-session-1",
					IsHello:   false,
				},
				ScriptRunId:    "run-1",
				Name:           "app.go",
				MainScriptPath: "/app.go",
				Config: &pb.Config{
					GatherUsageStats:    false,
					MaxCachedMessageAge: 0,
					AllowRunOnSave:      true,
					HideTopBar:          false,
					HideSidebarNav:      false,
				},
				AppPages:       []*pb.AppPage{},
				PageScriptHash: "hash-1",
			},
		},
	}
	sendProtoMsg(conn, sessionMsg)

	// Send the main vertical block that will contain elements
	// Delta path [0] means main container
	mainBlockDelta := &pb.ForwardMsg{
		Metadata: &pb.ForwardMsgMetadata{
			DeltaPath:        []uint32{0},
			ActiveScriptHash: "main",
		},
		Type: &pb.ForwardMsg_Delta{
			Delta: &pb.Delta{
				Type: &pb.Delta_AddBlock{
					AddBlock: &pb.Block{
						Type: &pb.Block_Vertical_{
							Vertical: &pb.Block_Vertical{
								Border: false,
							},
						},
						AllowEmpty: true,
					},
				},
			},
		},
	}
	sendProtoMsg(conn, mainBlockDelta)

	// Send the sidebar vertical block
	// Delta path [1] means sidebar container (RootContainer.SIDEBAR)
	sidebarBlockDelta := &pb.ForwardMsg{
		Metadata: &pb.ForwardMsgMetadata{
			DeltaPath:        []uint32{1},
			ActiveScriptHash: "main",
		},
		Type: &pb.ForwardMsg_Delta{
			Delta: &pb.Delta{
				Type: &pb.Delta_AddBlock{
					AddBlock: &pb.Block{
						Type: &pb.Block_Vertical_{
							Vertical: &pb.Block_Vertical{
								Border: false,
							},
						},
						AllowEmpty: true,
					},
				},
			},
		},
	}
	sendProtoMsg(conn, sidebarBlockDelta)

	// Get all main elements
	mu.Lock()
	mainElements := make([]Element, len(elements))
	copy(mainElements, elements)
	mu.Unlock()

	// Get all sidebar elements
	sidebarMu.Lock()
	sidebarElems := make([]Element, len(sidebarElements))
	copy(sidebarElems, sidebarElements)
	sidebarMu.Unlock()

	// Send each element as a Delta
	// Delta path [0, i] means element i in main container
	elementIndex := uint32(0)
	for _, elem := range mainElements {
		// Convert element to JSON to get the data
		data, err := elem.Render()
		if err != nil {
			log.Printf("Failed to render element: %v", err)
			continue
		}

		var elemMap map[string]interface{}
		if err := json.Unmarshal(data, &elemMap); err != nil {
			log.Printf("Failed to parse element: %v", err)
			continue
		}

		// Special handling for columns
		if elemType, ok := elemMap["type"].(string); ok && elemType == "columns" {
			// Send the horizontal flex container
			flexContainer := &pb.ForwardMsg{
				Metadata: &pb.ForwardMsgMetadata{
					DeltaPath:        []uint32{0, elementIndex},
					ActiveScriptHash: "main",
				},
				Type: &pb.ForwardMsg_Delta{
					Delta: &pb.Delta{
						Type: &pb.Delta_AddBlock{
							AddBlock: &pb.Block{
								Type: &pb.Block_FlexContainer_{
									FlexContainer: &pb.Block_FlexContainer{
										Direction: pb.Block_FlexContainer_HORIZONTAL,
										Wrap:      true,
										GapConfig: &pb.GapConfig{
											GapSpec: &pb.GapConfig_GapSize{
												GapSize: pb.GapSize_SMALL,
											},
										},
										Scale: 1.0,
										Align: pb.Block_FlexContainer_STRETCH,
									},
								},
								AllowEmpty: true,
							},
						},
					},
				},
			}
			sendProtoMsg(conn, flexContainer)

			// Send individual column blocks
			if cols, ok := elemMap["columns"].([]interface{}); ok {
				numCols := len(cols)
				weight := 1.0 / float64(numCols)

				for colIdx, col := range cols {
					colBlock := &pb.ForwardMsg{
						Metadata: &pb.ForwardMsgMetadata{
							DeltaPath:        []uint32{0, elementIndex, uint32(colIdx)},
							ActiveScriptHash: "main",
						},
						Type: &pb.ForwardMsg_Delta{
							Delta: &pb.Delta{
								Type: &pb.Delta_AddBlock{
									AddBlock: &pb.Block{
										Type: &pb.Block_Column_{
											Column: &pb.Block_Column{
												Weight:            weight,
												VerticalAlignment: pb.Block_Column_TOP,
												ShowBorder:        false,
												GapConfig: &pb.GapConfig{
													GapSpec: &pb.GapConfig_GapSize{
														GapSize: pb.GapSize_SMALL,
													},
												},
											},
										},
										AllowEmpty: true,
									},
								},
							},
						},
					}
					sendProtoMsg(conn, colBlock)

					// Send elements inside this column
					if colMap, ok := col.(map[string]interface{}); ok {
						if elems, ok := colMap["elements"].([]interface{}); ok {
							for elemIdx, colElem := range elems {
								if colElemMap, ok := colElem.(map[string]interface{}); ok {
									protoElement := convertElementToProto(colElemMap)
									if protoElement != nil {
										delta := &pb.ForwardMsg{
											Metadata: &pb.ForwardMsgMetadata{
												DeltaPath:        []uint32{0, elementIndex, uint32(colIdx), uint32(elemIdx)},
												ActiveScriptHash: "main",
											},
											Type: &pb.ForwardMsg_Delta{
												Delta: &pb.Delta{
													Type: &pb.Delta_NewElement{
														NewElement: protoElement,
													},
												},
											},
										}
										sendProtoMsg(conn, delta)
									}
								}
							}
						}
					}
				}
			}
			elementIndex++
			continue
		}

		// Convert the element based on its type
		protoElement := convertElementToProto(elemMap)
		if protoElement == nil {
			log.Printf("Skipping unsupported element type: %v", elemMap["type"])
			continue
		}

		delta := &pb.ForwardMsg{
			Metadata: &pb.ForwardMsgMetadata{
				DeltaPath:        []uint32{0, elementIndex},
				ActiveScriptHash: "main",
			},
			Type: &pb.ForwardMsg_Delta{
				Delta: &pb.Delta{
					Type: &pb.Delta_NewElement{
						NewElement: protoElement,
					},
				},
			},
		}
		sendProtoMsg(conn, delta)
		elementIndex++
	}

	log.Printf("Sent %d main protobuf elements", elementIndex)

	// Send sidebar elements
	// Delta path [1, i] means element i in sidebar container
	sidebarIndex := uint32(0)
	log.Printf("Starting to send %d sidebar elements...", len(sidebarElems))
	for _, elem := range sidebarElems {
		// Convert element to JSON to get the data
		data, err := elem.Render()
		if err != nil {
			log.Printf("Failed to render sidebar element: %v", err)
			continue
		}

		var elemMap map[string]interface{}
		if err := json.Unmarshal(data, &elemMap); err != nil {
			log.Printf("Failed to parse sidebar element: %v", err)
			continue
		}

		// Convert the element based on its type
		protoElement := convertElementToProto(elemMap)
		if protoElement == nil {
			log.Printf("Skipping unsupported sidebar element type: %v", elemMap["type"])
			continue
		}

		delta := &pb.ForwardMsg{
			Metadata: &pb.ForwardMsgMetadata{
				DeltaPath:        []uint32{1, sidebarIndex},
				ActiveScriptHash: "main",
			},
			Type: &pb.ForwardMsg_Delta{
				Delta: &pb.Delta{
					Type: &pb.Delta_NewElement{
						NewElement: protoElement,
					},
				},
			},
		}
		sendProtoMsg(conn, delta)
		log.Printf("Sent sidebar element %d: type=%v", sidebarIndex, elemMap["type"])
		sidebarIndex++
	}

	log.Printf("Sent %d sidebar protobuf elements", sidebarIndex)
}

// convertElementToProto converts a JSON element map to a protobuf Element
func convertElementToProto(elemMap map[string]interface{}) *pb.Element {
	elemType, ok := elemMap["type"].(string)
	if !ok {
		return nil
	}

	switch elemType {
	case "text":
		content, _ := elemMap["content"].(string)
		help, _ := elemMap["help"].(string)
		return &pb.Element{
			Type: &pb.Element_Text{
				Text: &pb.Text{
					Body: content,
					Help: help,
				},
			},
		}

	case "markdown":
		content, _ := elemMap["content"].(string)
		help, _ := elemMap["help"].(string)
		return &pb.Element{
			Type: &pb.Element_Markdown{
				Markdown: &pb.Markdown{
					Body: content,
					Help: help,
				},
			},
		}

	case "heading":
		content, _ := elemMap["content"].(string)
		help, _ := elemMap["help"].(string)
		anchor, _ := elemMap["anchor"].(string)

		// Convert level to tag
		var tag string
		if level, ok := elemMap["level"].(float64); ok {
			switch int(level) {
			case 1:
				tag = "h1"
			case 2:
				tag = "h2"
			case 3:
				tag = "h3"
			default:
				tag = "h2"
			}
		} else {
			tag = "h2"
		}

		return &pb.Element{
			Type: &pb.Element_Heading{
				Heading: &pb.Heading{
					Body:   content,
					Tag:    tag,
					Anchor: anchor,
					Help:   help,
				},
			},
		}

	case "code":
		content, _ := elemMap["content"].(string)
		language, _ := elemMap["language"].(string)
		return &pb.Element{
			Type: &pb.Element_Code{
				Code: &pb.Code{
					CodeText: content,
					Language: language,
				},
			},
		}

	case "divider":
		return &pb.Element{
			Type: &pb.Element_Markdown{
				Markdown: &pb.Markdown{
					Body:        "---",
					ElementType: pb.Markdown_DIVIDER,
				},
			},
		}

	case "alert":
		body, _ := elemMap["body"].(string)
		format, _ := elemMap["format"].(string)
		icon, _ := elemMap["icon"].(string)

		var alertFormat pb.Alert_Format
		switch format {
		case "success":
			alertFormat = pb.Alert_SUCCESS
		case "info":
			alertFormat = pb.Alert_INFO
		case "warning":
			alertFormat = pb.Alert_WARNING
		case "error":
			alertFormat = pb.Alert_ERROR
		default:
			alertFormat = pb.Alert_INFO
		}

		return &pb.Element{
			Type: &pb.Element_Alert{
				Alert: &pb.Alert{
					Body:   body,
					Format: alertFormat,
					Icon:   icon,
				},
			},
		}

	case "balloons":
		return &pb.Element{
			Type: &pb.Element_Balloons{
				Balloons: &pb.Balloons{},
			},
		}

	case "line_chart":
		// Get data array
		var dataPoints []float64
		if data, ok := elemMap["data"].([]interface{}); ok {
			for _, val := range data {
				if fval, ok := val.(float64); ok {
					dataPoints = append(dataPoints, fval)
				}
			}
		}

		// Create Arrow schema (index: int64, value: float64)
		pool := memory.NewGoAllocator()
		schema := arrow.NewSchema(
			[]arrow.Field{
				{Name: "index", Type: arrow.PrimitiveTypes.Int64},
				{Name: "value", Type: arrow.PrimitiveTypes.Float64},
			},
			nil,
		)

		// Build record batch
		indexBuilder := array.NewInt64Builder(pool)
		defer indexBuilder.Release()
		valueBuilder := array.NewFloat64Builder(pool)
		defer valueBuilder.Release()

		for i, v := range dataPoints {
			indexBuilder.Append(int64(i))
			valueBuilder.Append(v)
		}

		indexArray := indexBuilder.NewArray()
		defer indexArray.Release()
		valueArray := valueBuilder.NewArray()
		defer valueArray.Release()

		record := array.NewRecord(schema, []arrow.Array{indexArray, valueArray}, int64(len(dataPoints)))
		defer record.Release()

		// Serialize to Arrow IPC format using a bytes buffer
		var buf bytes.Buffer
		writer := ipc.NewWriter(&buf, ipc.WithSchema(schema))
		defer writer.Close()

		if err := writer.Write(record); err != nil {
			log.Printf("Error writing Arrow record: %v", err)
			return nil
		}

		// Create a simple Vega-Lite spec for a line chart
		vegaSpec := `{
			"mark": "line",
			"encoding": {
				"x": {"field": "index", "type": "quantitative"},
				"y": {"field": "value", "type": "quantitative"}
			}
		}`

		// Generate unique ID based on data to force chart updates when data changes
		// Use first, last, and middle values plus length for a lightweight hash
		var chartID string
		if len(dataPoints) > 0 {
			mid := dataPoints[len(dataPoints)/2]
			chartID = fmt.Sprintf("line_chart_%d_%.2f_%.2f_%.2f",
				len(dataPoints), dataPoints[0], mid, dataPoints[len(dataPoints)-1])
		} else {
			chartID = "line_chart_empty"
		}

		return &pb.Element{
			Type: &pb.Element_ArrowVegaLiteChart{
				ArrowVegaLiteChart: &pb.ArrowVegaLiteChart{
					Spec:              vegaSpec,
					Data:              &pb.ArrowData{Data: buf.Bytes()},
					UseContainerWidth: true,
					Theme:             "streamlit",
					Id:                chartID,
				},
			},
		}

	case "button":
		label, _ := elemMap["label"].(string)
		key, _ := elemMap["key"].(string)
		help, _ := elemMap["help"].(string)
		return &pb.Element{
			Type: &pb.Element_Button{
				Button: &pb.Button{
					Label: label,
					Id:    key,
					Help:  help,
				},
			},
		}

	case "checkbox":
		label, _ := elemMap["label"].(string)
		key, _ := elemMap["key"].(string)
		defaultValue, _ := elemMap["defaultValue"].(bool)
		help, _ := elemMap["help"].(string)

		// Check for stored value
		setValue := false
		value := defaultValue
		widgetStateMu.RLock()
		if storedValue, ok := widgetState[key].(bool); ok {
			value = storedValue
			setValue = true
		}
		widgetStateMu.RUnlock()

		return &pb.Element{
			Type: &pb.Element_Checkbox{
				Checkbox: &pb.Checkbox{
					Label:    label,
					Id:       key,
					Default:  defaultValue,
					Value:    value,
					SetValue: setValue,
					Help:     help,
				},
			},
		}

	case "text_input":
		label, _ := elemMap["label"].(string)
		key, _ := elemMap["key"].(string)
		defaultValue, _ := elemMap["defaultValue"].(string)
		help, _ := elemMap["help"].(string)

		// Check for stored value
		setValue := false
		value := defaultValue
		widgetStateMu.RLock()
		if storedValue, ok := widgetState[key].(string); ok {
			value = storedValue
			setValue = true
		}
		widgetStateMu.RUnlock()

		return &pb.Element{
			Type: &pb.Element_TextInput{
				TextInput: &pb.TextInput{
					Label:    label,
					Id:       key,
					Default:  &defaultValue,
					Value:    &value,
					SetValue: setValue,
					Help:     help,
				},
			},
		}

	case "slider":
		label, _ := elemMap["label"].(string)
		key, _ := elemMap["key"].(string)
		help, _ := elemMap["help"].(string)

		// Handle numeric conversions
		var min, max, defaultValue, step float64
		if v, ok := elemMap["min"].(float64); ok {
			min = v
		} else if v, ok := elemMap["min"].(int); ok {
			min = float64(v)
		}
		if v, ok := elemMap["max"].(float64); ok {
			max = v
		} else if v, ok := elemMap["max"].(int); ok {
			max = float64(v)
		}
		if v, ok := elemMap["defaultValue"].(float64); ok {
			defaultValue = v
		} else if v, ok := elemMap["defaultValue"].(int); ok {
			defaultValue = float64(v)
		}
		if v, ok := elemMap["step"].(float64); ok {
			step = v
		} else if v, ok := elemMap["step"].(int); ok {
			step = float64(v)
		} else {
			// Default step to 1.0 if not specified
			step = 1.0
		}

		// Check if we have a stored value from widget state
		setValue := false
		value := []float64{defaultValue}
		widgetStateMu.RLock()
		if storedValue, ok := widgetState[key].(float64); ok {
			value = []float64{storedValue}
			setValue = true
		} else if storedInt, ok := widgetState[key].(int64); ok {
			value = []float64{float64(storedInt)}
			setValue = true
		}
		widgetStateMu.RUnlock()

		return &pb.Element{
			Type: &pb.Element_Slider{
				Slider: &pb.Slider{
					Label:    label,
					Id:       key,
					Min:      min,
					Max:      max,
					Step:     step,
					Default:  []float64{defaultValue},
					Value:    value,
					SetValue: setValue,
					Help:     help,
				},
			},
		}

	case "selectbox":
		label, _ := elemMap["label"].(string)
		key, _ := elemMap["key"].(string)
		help, _ := elemMap["help"].(string)

		var defaultIndex int32
		if v, ok := elemMap["defaultIndex"].(float64); ok {
			defaultIndex = int32(v)
		} else if v, ok := elemMap["defaultIndex"].(int); ok {
			defaultIndex = int32(v)
		}

		var options []string
		if opts, ok := elemMap["options"].([]interface{}); ok {
			for _, opt := range opts {
				if str, ok := opt.(string); ok {
					options = append(options, str)
				}
			}
		}

		// Check for stored value (selectbox stores the selected string)
		setValue := false
		var rawValue *string
		widgetStateMu.RLock()
		if storedValue, ok := widgetState[key].(string); ok {
			rawValue = &storedValue
			setValue = true
		}
		widgetStateMu.RUnlock()

		return &pb.Element{
			Type: &pb.Element_Selectbox{
				Selectbox: &pb.Selectbox{
					Label:    label,
					Id:       key,
					Default:  &defaultIndex,
					RawValue: rawValue,
					SetValue: setValue,
					Options:  options,
					Help:     help,
				},
			},
		}

	case "metric":
		label, _ := elemMap["label"].(string)
		value, _ := elemMap["value"].(string)
		delta, _ := elemMap["delta"].(string)
		help, _ := elemMap["help"].(string)

		return &pb.Element{
			Type: &pb.Element_Metric{
				Metric: &pb.Metric{
					Label: label,
					Body:  value,
					Delta: delta,
					Help:  help,
				},
			},
		}

	default:
		// For unsupported types, render as text for debugging
		log.Printf("Unknown element type: %s, rendering as text", elemType)
		return &pb.Element{
			Type: &pb.Element_Text{
				Text: &pb.Text{
					Body: fmt.Sprintf("[Unsupported: %s] %v", elemType, elemMap),
				},
			},
		}
	}
}

// sendProtoMsg sends a protobuf ForwardMsg over WebSocket
func sendProtoMsg(conn *websocket.Conn, msg *pb.ForwardMsg) {
	data, err := proto.Marshal(msg)
	if err != nil {
		log.Printf("Failed to marshal protobuf: %v", err)
		return
	}

	if err := conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
		log.Printf("Failed to send protobuf: %v", err)
	}
}

// findFrontendBuild looks for the frontend build directory in multiple locations
func findFrontendBuild() string {
	// Get the executable's directory
	exePath, err := os.Executable()
	if err != nil {
		return ""
	}
	exeDir := filepath.Dir(exePath)

	// Possible frontend locations relative to the executable
	candidates := []string{
		// UPSTREAM-COMPATIBLE: Use go-backend/static (not lib/streamlit/static)
		// From go-backend/examples/X/ -> go-backend/static/
		filepath.Join(exeDir, "..", "..", "static"),
		// From go-backend/bin/ -> go-backend/static/
		filepath.Join(exeDir, "..", "static"),
		// From go-backend/ -> go-backend/static/
		filepath.Join(exeDir, "static"),
		// Fallback to original Streamlit location (for backwards compatibility)
		filepath.Join(exeDir, "..", "..", "..", "lib", "streamlit", "static"),
		filepath.Join(exeDir, "..", "..", "lib", "streamlit", "static"),
		filepath.Join(exeDir, "lib", "streamlit", "static"),
	}

	for _, dir := range candidates {
		absDir, err := filepath.Abs(dir)
		if err != nil {
			continue
		}
		// Check if index.html exists
		indexPath := filepath.Join(absDir, "index.html")
		if _, err := os.Stat(indexPath); err == nil {
			return absDir
		}
	}

	return ""
}

// Run executes the Streamlit app with embedded server
func Run(app func()) error {
	appFunc = app
	port := 8501

	// Set up HTTP server
	mux := http.NewServeMux()

	// WebSocket endpoint for Streamlit protocol
	mux.HandleFunc("/_stcore/stream", handleWebSocket)

	// Health check endpoint
	mux.HandleFunc("/_stcore/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "ok")
	})

	// Host config endpoint - required by React frontend
	mux.HandleFunc("/_stcore/host-config", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		config := map[string]interface{}{
			"allowedOrigins":             []string{"*"},
			"useExternalAuthToken":       false,
			"enableCustomParentMessages": false,
			"enforceDownloadInNewTab":    false,
			"metricsUrl":                 "",
			"blockErrorDialogs":          false,
			"resourceCrossOriginMode":    nil,
		}
		json.NewEncoder(w).Encode(config)
	})

	// Try to serve the built frontend, fall back to simple HTML if not available
	staticDir := findFrontendBuild()
	if staticDir != "" {
		log.Printf("Serving frontend from: %s", staticDir)

		// Create a file server for the static directory
		fs := http.FileServer(http.Dir(staticDir))

		// Serve the root path - serve index.html or files from the static directory
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			// Serve the request with the file server
			fs.ServeHTTP(w, r)
		})
	} else {
		log.Println("Frontend not found, serving simple HTML page")
		log.Println("To use the full UI, build the frontend with: cd ../frontend && make frontend-fast")

		// Serve a simple HTML page (placeholder until frontend is built)
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			html := `<!DOCTYPE html>
<html>
<head>
    <title>Streamlit Go App</title>
    <style>
        body { font-family: sans-serif; max-width: 800px; margin: 50px auto; padding: 20px; }
        .status { padding: 20px; background: #f0f2f6; border-radius: 8px; margin: 20px 0; }
        code { background: #e8eaed; padding: 2px 6px; border-radius: 3px; }
    </style>
</head>
<body>
    <h1>🎈 Streamlit Go App</h1>
    <div class="status">
        <h3>✅ Server is running</h3>
        <p>WebSocket: <code>ws://localhost:8501/_stcore/stream</code></p>
        <p><strong>Note:</strong> Frontend UI not available. Build it with:</p>
        <code>cd ../frontend && make frontend-fast</code>
    </div>
    <script>
        const ws = new WebSocket("ws://localhost:8501/_stcore/stream");
        ws.onopen = () => console.log("Connected");
        ws.onmessage = (e) => console.log("Element:", e.data);
    </script>
</body>
</html>`
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprint(w, html)
		})
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	// Start server in goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait a moment for server to start
	time.Sleep(100 * time.Millisecond)

	log.Printf("🎈 Streamlit app is running!")
	log.Printf("   Local URL: http://localhost:%d", port)
	log.Printf("   Press Ctrl+C to stop")

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("\nShutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	log.Println("Server stopped")
	return nil
}

// Write writes text to the app
func Write(args ...interface{}) {
	addElement(&TextElement{Content: fmt.Sprint(args...)})
}

// Writef writes formatted text to the app
func Writef(format string, args ...interface{}) {
	addElement(&TextElement{Content: fmt.Sprintf(format, args...)})
}

// Title displays a title
func Title(text string) {
	addElement(&HeadingElement{Content: text, Level: 1})
}

// Header displays a header
func Header(text string) {
	addElement(&HeadingElement{Content: text, Level: 2})
}

// Subheader displays a subheader
func Subheader(text string) {
	addElement(&HeadingElement{Content: text, Level: 3})
}

// Markdown displays markdown content
func Markdown(content string) {
	addElement(&MarkdownElement{Content: content})
}

// Code displays a code block
func Code(code string, language string) {
	addElement(&CodeElement{Content: code, Language: language})
}

// Divider displays a horizontal divider
func Divider() {
	addElement(&DividerElement{})
}

// LineChart displays a line chart
func LineChart(data []float64) {
	addElement(&LineChartElement{Data: data})
}

// Success displays a success message
func Success(message string) {
	addElement(&AlertElement{Message: message, Type: "success"})
}

// Info displays an info message
func Info(message string) {
	addElement(&AlertElement{Message: message, Type: "info"})
}

// Warning displays a warning message
func Warning(message string) {
	addElement(&AlertElement{Message: message, Type: "warning"})
}

// Error displays an error message
func Error(message string) {
	addElement(&AlertElement{Message: message, Type: "error"})
}

// Button displays a button widget
func Button(label string) bool {
	widgetStateMu.RLock()
	clicked, _ := widgetState[label].(bool)
	widgetStateMu.RUnlock()

	// Clear the clicked state after reading
	if clicked {
		widgetStateMu.Lock()
		delete(widgetState, label)
		widgetStateMu.Unlock()
	}

	addElement(&ButtonElement{Label: label, Key: label})
	return clicked
}

// TextInput displays a text input widget
func TextInput(label string, value string) string {
	widgetStateMu.RLock()
	if storedValue, ok := widgetState[label].(string); ok {
		value = storedValue
	}
	widgetStateMu.RUnlock()

	addElement(&TextInputElement{Label: label, DefaultValue: value, Key: label})
	return value
}

// Slider displays a slider widget
func Slider(label string, min float64, max float64, value float64) float64 {
	widgetStateMu.RLock()
	if storedValue, ok := widgetState[label].(float64); ok {
		value = storedValue
	} else if storedInt, ok := widgetState[label].(int); ok {
		value = float64(storedInt)
	}
	widgetStateMu.RUnlock()

	addElement(&SliderElement{Label: label, Min: int(min), Max: int(max), DefaultValue: int(value), Key: label})
	return value
}

// Checkbox displays a checkbox widget
func Checkbox(label string, value bool) bool {
	widgetStateMu.RLock()
	if storedValue, ok := widgetState[label].(bool); ok {
		value = storedValue
	}
	widgetStateMu.RUnlock()

	addElement(&CheckboxElement{Label: label, DefaultValue: value, Key: label})
	return value
}

// Selectbox displays a selectbox widget
func Selectbox(label string, options []string, index int) string {
	widgetStateMu.RLock()
	if storedIndex, ok := widgetState[label].(float64); ok {
		index = int(storedIndex)
	} else if storedInt, ok := widgetState[label].(int); ok {
		index = storedInt
	} else if storedStr, ok := widgetState[label].(string); ok {
		// String value from widget state - find the index
		for i, opt := range options {
			if opt == storedStr {
				index = i
				break
			}
		}
	}
	widgetStateMu.RUnlock()

	if index < 0 || index >= len(options) {
		index = 0
	}
	addElement(&SelectboxElement{Label: label, Options: options, DefaultIndex: index, Key: label})
	return options[index]
}

// Columns creates a multi-column layout
func Columns(count int) []*Column {
	cols := make([]*Column, count)
	for i := 0; i < count; i++ {
		cols[i] = &Column{index: i}
	}
	addElement(&ColumnsElement{Columns: cols})
	return cols
}

// Metric displays a metric with optional delta
func Metric(label string, value string, delta string) {
	addElement(&MetricElement{Label: label, Value: value, Delta: delta})
}

// Balloons triggers a balloons animation
func Balloons() {
	addElement(&BalloonsElement{})
}

// Snow triggers a snow animation
func Snow() {
	addElement(&SnowElement{})
}

// Sidebar type for sidebar context
type SidebarContext struct{}

var sidebarInstance = &SidebarContext{}

// Sidebar returns the sidebar context
func Sidebar() *SidebarContext {
	return sidebarInstance
}

// Write writes content to the sidebar
func (s *SidebarContext) Write(args ...interface{}) {
	inSidebar = true
	addElement(&TextElement{Content: fmt.Sprint(args...)})
	inSidebar = false
}

// Writef writes formatted content to the sidebar
func (s *SidebarContext) Writef(format string, args ...interface{}) {
	inSidebar = true
	addElement(&TextElement{Content: fmt.Sprintf(format, args...)})
	inSidebar = false
}

// Title displays a title in the sidebar
func (s *SidebarContext) Title(text string) {
	inSidebar = true
	addElement(&HeadingElement{Content: text, Level: 1})
	inSidebar = false
}

// Header displays a header in the sidebar
func (s *SidebarContext) Header(text string) {
	inSidebar = true
	addElement(&HeadingElement{Content: text, Level: 2})
	inSidebar = false
}

// Subheader displays a subheader in the sidebar
func (s *SidebarContext) Subheader(text string) {
	inSidebar = true
	addElement(&HeadingElement{Content: text, Level: 3})
	inSidebar = false
}

// Divider displays a divider in the sidebar
func (s *SidebarContext) Divider() {
	inSidebar = true
	addElement(&DividerElement{})
	inSidebar = false
}

// Button displays a button in the sidebar
func (s *SidebarContext) Button(label string) bool {
	inSidebar = true
	defer func() { inSidebar = false }()

	widgetStateMu.Lock()
	clicked := false
	if val, ok := widgetState[label]; ok && val == true {
		clicked = true
		widgetState[label] = false
	}
	widgetStateMu.Unlock()

	addElement(&ButtonElement{Label: label, Key: label})
	return clicked
}

// TextInput displays a text input in the sidebar
func (s *SidebarContext) TextInput(label string, defaultValue string) string {
	inSidebar = true
	defer func() { inSidebar = false }()

	value := defaultValue
	widgetStateMu.RLock()
	if storedValue, ok := widgetState[label].(string); ok {
		value = storedValue
	}
	widgetStateMu.RUnlock()

	addElement(&TextInputElement{Label: label, DefaultValue: value, Key: label})
	return value
}

// Slider displays a slider in the sidebar
func (s *SidebarContext) Slider(label string, min float64, max float64, value float64) float64 {
	inSidebar = true
	defer func() { inSidebar = false }()

	widgetStateMu.RLock()
	if storedValue, ok := widgetState[label].(float64); ok {
		value = storedValue
		log.Printf("Slider %s: using stored value %f", label, value)
	} else if storedInt, ok := widgetState[label].(int); ok {
		value = float64(storedInt)
		log.Printf("Slider %s: using stored int value %d", label, storedInt)
	} else {
		log.Printf("Slider %s: using default value %f (stored value type: %T)", label, value, widgetState[label])
	}
	widgetStateMu.RUnlock()

	addElement(&SliderElement{Label: label, Min: int(min), Max: int(max), DefaultValue: int(value), Key: label})
	return value
}

// Checkbox displays a checkbox in the sidebar
func (s *SidebarContext) Checkbox(label string, value bool) bool {
	inSidebar = true
	defer func() { inSidebar = false }()

	widgetStateMu.RLock()
	if storedValue, ok := widgetState[label].(bool); ok {
		value = storedValue
	}
	widgetStateMu.RUnlock()

	addElement(&CheckboxElement{Label: label, DefaultValue: value, Key: label})
	return value
}

// Selectbox displays a selectbox in the sidebar
func (s *SidebarContext) Selectbox(label string, options []string, index int) string {
	inSidebar = true
	defer func() { inSidebar = false }()

	widgetStateMu.RLock()
	if storedIndex, ok := widgetState[label].(float64); ok {
		index = int(storedIndex)
	} else if storedInt, ok := widgetState[label].(int); ok {
		index = storedInt
	} else if storedStr, ok := widgetState[label].(string); ok {
		// String value from widget state - find the index
		for i, opt := range options {
			if opt == storedStr {
				index = i
				break
			}
		}
	}
	widgetStateMu.RUnlock()

	if index < 0 || index >= len(options) {
		index = 0
	}
	addElement(&SelectboxElement{Label: label, Options: options, DefaultIndex: index, Key: label})
	return options[index]
}
