package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// LogMessage represents a streamed log message
type LogMessage struct {
	Type      string                 `json:"type"`
	JobRunID  string                 `json:"job_run_id"`
	TaskRunID string                 `json:"task_run_id,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	Level     LogLevel               `json:"level"`
	Stream    string                 `json:"stream"`
	Message   string                 `json:"message"`
	TaskName  string                 `json:"task_name,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// LogStreamer manages real-time log streaming to clients
type LogStreamer struct {
	// WebSocket connections for each job run
	wsConnections map[string]map[*websocket.Conn]bool
	wsLock        sync.RWMutex

	// SSE connections for each job run
	sseConnections map[string]map[chan LogMessage]bool
	sseLock        sync.RWMutex

	// Global connections (all job runs)
	globalWSConnections  map[*websocket.Conn]bool
	globalSSEConnections map[chan LogMessage]bool
	globalLock           sync.RWMutex

	upgrader websocket.Upgrader
}

// NewLogStreamer creates a new log streamer
func NewLogStreamer() *LogStreamer {
	return &LogStreamer{
		wsConnections:        make(map[string]map[*websocket.Conn]bool),
		sseConnections:       make(map[string]map[chan LogMessage]bool),
		globalWSConnections:  make(map[*websocket.Conn]bool),
		globalSSEConnections: make(map[chan LogMessage]bool),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// In production, implement proper CORS validation
				return true
			},
		},
	}
}

// BroadcastLog sends a log message to all connected clients
func (ls *LogStreamer) BroadcastLog(logMsg LogMessage) {
	// Broadcast to job-specific connections
	ls.broadcastToJobConnections(logMsg.JobRunID, logMsg)

	// Broadcast to global connections
	ls.broadcastToGlobalConnections(logMsg)
}

// broadcastToJobConnections sends log to job-specific connections
func (ls *LogStreamer) broadcastToJobConnections(jobRunID string, logMsg LogMessage) {
	// WebSocket connections
	ls.wsLock.RLock()
	if connections, exists := ls.wsConnections[jobRunID]; exists {
		for conn := range connections {
			if err := conn.WriteJSON(logMsg); err != nil {
				log.Printf("Failed to write to websocket: %v", err)
				conn.Close()
				delete(connections, conn)
			}
		}
	}
	ls.wsLock.RUnlock()

	// SSE connections
	ls.sseLock.RLock()
	if connections, exists := ls.sseConnections[jobRunID]; exists {
		for ch := range connections {
			select {
			case ch <- logMsg:
			default:
				// Channel is full, remove it
				close(ch)
				delete(connections, ch)
			}
		}
	}
	ls.sseLock.RUnlock()
}

// broadcastToGlobalConnections sends log to global connections
func (ls *LogStreamer) broadcastToGlobalConnections(logMsg LogMessage) {
	ls.globalLock.RLock()
	defer ls.globalLock.RUnlock()

	// Global WebSocket connections
	for conn := range ls.globalWSConnections {
		if err := conn.WriteJSON(logMsg); err != nil {
			log.Printf("Failed to write to global websocket: %v", err)
			conn.Close()
			delete(ls.globalWSConnections, conn)
		}
	}

	// Global SSE connections
	for ch := range ls.globalSSEConnections {
		select {
		case ch <- logMsg:
		default:
			// Channel is full, remove it
			close(ch)
			delete(ls.globalSSEConnections, ch)
		}
	}
}

// HandleWebSocket handles WebSocket connections for log streaming
func (ls *LogStreamer) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := ls.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade websocket: %v", err)
		return
	}
	defer conn.Close()

	// Get job run ID from query parameter
	jobRunID := r.URL.Query().Get("job_run_id")

	if jobRunID != "" {
		// Job-specific connection
		ls.addJobWSConnection(jobRunID, conn)
		defer ls.removeJobWSConnection(jobRunID, conn)
	} else {
		// Global connection (all jobs)
		ls.addGlobalWSConnection(conn)
		defer ls.removeGlobalWSConnection(conn)
	}

	// Send initial connection confirmation
	confirmMsg := LogMessage{
		JobRunID:  jobRunID,
		Timestamp: time.Now(),
		Level:     InfoLevel,
		Stream:    "system",
		Message:   "Connected to log stream",
	}

	if err := conn.WriteJSON(confirmMsg); err != nil {
		log.Printf("Failed to send confirmation: %v", err)
		return
	}

	// Keep connection alive and handle client messages
	for {
		var msg map[string]interface{}
		if err := conn.ReadJSON(&msg); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}
		// Handle any client messages if needed (ping/pong, etc.)
	}
}

// HandleSSE handles Server-Sent Events for log streaming
func (ls *LogStreamer) HandleSSE(w http.ResponseWriter, r *http.Request) {
	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Create a channel for this connection
	logChan := make(chan LogMessage, 100)

	// Get job run ID from query parameter
	jobRunID := r.URL.Query().Get("job_run_id")

	if jobRunID != "" {
		// Job-specific connection
		ls.addJobSSEConnection(jobRunID, logChan)
		defer ls.removeJobSSEConnection(jobRunID, logChan)
	} else {
		// Global connection
		ls.addGlobalSSEConnection(logChan)
		defer ls.removeGlobalSSEConnection(logChan)
	}

	// Send initial connection event
	fmt.Fprintf(w, "event: connected\n")
	fmt.Fprintf(w, "data: {\"message\":\"Connected to log stream\",\"job_run_id\":\"%s\"}\n\n", jobRunID)

	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}

	// Stream logs
	ctx := r.Context()
	for {
		select {
		case <-ctx.Done():
			return
		case logMsg := <-logChan:
			data, err := json.Marshal(logMsg)
			if err != nil {
				log.Printf("Failed to marshal log message: %v", err)
				continue
			}

			fmt.Fprintf(w, "event: log\n")
			fmt.Fprintf(w, "data: %s\n\n", data)

			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
		case <-time.After(30 * time.Second):
			// Send keep-alive ping
			fmt.Fprintf(w, "event: ping\n")
			fmt.Fprintf(w, "data: {\"timestamp\":\"%s\"}\n\n", time.Now().Format(time.RFC3339))

			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
		}
	}
}

// Connection management methods
func (ls *LogStreamer) addJobWSConnection(jobRunID string, conn *websocket.Conn) {
	ls.wsLock.Lock()
	defer ls.wsLock.Unlock()

	if ls.wsConnections[jobRunID] == nil {
		ls.wsConnections[jobRunID] = make(map[*websocket.Conn]bool)
	}
	ls.wsConnections[jobRunID][conn] = true
}

func (ls *LogStreamer) removeJobWSConnection(jobRunID string, conn *websocket.Conn) {
	ls.wsLock.Lock()
	defer ls.wsLock.Unlock()

	if connections, exists := ls.wsConnections[jobRunID]; exists {
		delete(connections, conn)
		if len(connections) == 0 {
			delete(ls.wsConnections, jobRunID)
		}
	}
}

func (ls *LogStreamer) addJobSSEConnection(jobRunID string, ch chan LogMessage) {
	ls.sseLock.Lock()
	defer ls.sseLock.Unlock()

	if ls.sseConnections[jobRunID] == nil {
		ls.sseConnections[jobRunID] = make(map[chan LogMessage]bool)
	}
	ls.sseConnections[jobRunID][ch] = true
}

func (ls *LogStreamer) removeJobSSEConnection(jobRunID string, ch chan LogMessage) {
	ls.sseLock.Lock()
	defer ls.sseLock.Unlock()

	if connections, exists := ls.sseConnections[jobRunID]; exists {
		delete(connections, ch)
		if len(connections) == 0 {
			delete(ls.sseConnections, jobRunID)
		}
	}
	close(ch)
}

func (ls *LogStreamer) addGlobalWSConnection(conn *websocket.Conn) {
	ls.globalLock.Lock()
	defer ls.globalLock.Unlock()
	ls.globalWSConnections[conn] = true
}

func (ls *LogStreamer) removeGlobalWSConnection(conn *websocket.Conn) {
	ls.globalLock.Lock()
	defer ls.globalLock.Unlock()
	delete(ls.globalWSConnections, conn)
}

func (ls *LogStreamer) addGlobalSSEConnection(ch chan LogMessage) {
	ls.globalLock.Lock()
	defer ls.globalLock.Unlock()
	ls.globalSSEConnections[ch] = true
}

func (ls *LogStreamer) removeGlobalSSEConnection(ch chan LogMessage) {
	ls.globalLock.Lock()
	defer ls.globalLock.Unlock()
	delete(ls.globalSSEConnections, ch)
	close(ch)
}

// GetActiveConnections returns the number of active connections
func (ls *LogStreamer) GetActiveConnections() map[string]interface{} {
	ls.wsLock.RLock()
	ls.sseLock.RLock()
	ls.globalLock.RLock()
	defer ls.wsLock.RUnlock()
	defer ls.sseLock.RUnlock()
	defer ls.globalLock.RUnlock()

	jobConnections := make(map[string]map[string]int)

	// Count job-specific connections
	for jobRunID, wsConns := range ls.wsConnections {
		if jobConnections[jobRunID] == nil {
			jobConnections[jobRunID] = make(map[string]int)
		}
		jobConnections[jobRunID]["websocket"] = len(wsConns)
	}

	for jobRunID, sseConns := range ls.sseConnections {
		if jobConnections[jobRunID] == nil {
			jobConnections[jobRunID] = make(map[string]int)
		}
		jobConnections[jobRunID]["sse"] = len(sseConns)
	}

	return map[string]interface{}{
		"global": map[string]int{
			"websocket": len(ls.globalWSConnections),
			"sse":       len(ls.globalSSEConnections),
		},
		"jobs": jobConnections,
	}
}
