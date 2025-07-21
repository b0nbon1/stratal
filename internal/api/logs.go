package api

import (
	"net/http"
	"strconv"

	"github.com/b0nbon1/stratal/internal/logger"
	db "github.com/b0nbon1/stratal/internal/storage/db/sqlc"
	"github.com/b0nbon1/stratal/pkg/router"
	"github.com/b0nbon1/stratal/pkg/utils"
)

// LogsHandler handles log-related HTTP requests
type LogsHandler struct {
	store    *db.SQLStore
	streamer *logger.LogStreamer
}

// NewLogsHandler creates a new logs handler
func NewLogsHandler(store *db.SQLStore, streamer *logger.LogStreamer) *LogsHandler {
	return &LogsHandler{
		store:    store,
		streamer: streamer,
	}
}

// RegisterLogRoutes registers log-related routes
func (h *LogsHandler) RegisterLogRoutes(v1 *router.Router) {
	// Streaming endpoints
	v1.Get("/logs/stream/ws", h.HandleWebSocketStream)  // WebSocket streaming
	v1.Get("/logs/stream/sse", h.HandleSSEStream)       // Server-Sent Events streaming
	v1.Get("/logs/stream/status", h.GetStreamingStatus) // Get streaming connection status

	// REST API endpoints
	v1.Get("/logs/job-runs/:job_run_id", h.GetJobRunLogs)               // Get all logs for a job run
	v1.Get("/logs/task-runs/:task_run_id", h.GetTaskRunLogs)            // Get logs for a specific task run
	v1.Get("/logs/job-runs/:job_run_id/download", h.DownloadJobRunLogs) // Download log file
	v1.Get("/logs/system", h.GetSystemLogs)                             // Get system logs
	v1.Get("/logs/type/:type", h.GetLogsByType)                         // Get logs by type (system, job, task)
}

// HandleWebSocketStream handles WebSocket log streaming
func (h *LogsHandler) HandleWebSocketStream(w http.ResponseWriter, r *http.Request) {
	h.streamer.HandleWebSocket(w, r)
}

// HandleSSEStream handles Server-Sent Events log streaming
func (h *LogsHandler) HandleSSEStream(w http.ResponseWriter, r *http.Request) {
	h.streamer.HandleSSE(w, r)
}

// GetStreamingStatus returns the current streaming connection status
func (h *LogsHandler) GetStreamingStatus(w http.ResponseWriter, r *http.Request) {
	status := h.streamer.GetActiveConnections()
	respondJSON(w, 200, map[string]interface{}{
		"status":      "active",
		"connections": status,
	})
}

// GetJobRunLogs retrieves all logs for a specific job run
func (h *LogsHandler) GetJobRunLogs(w http.ResponseWriter, r *http.Request) {
	jobRunIDStr := r.PathValue("job_run_id")
	if jobRunIDStr == "" {
		respondJSON(w, 400, map[string]string{"error": "job_run_id is required"})
		return
	}

	jobRunID, err := utils.ParseUUID(jobRunIDStr)
	if err != nil {
		respondJSON(w, 400, map[string]string{"error": "Invalid job_run_id format"})
		return
	}

	// Get pagination parameters
	limit, offset := getPaginationParams(r)

	logs, err := h.store.ListLogsByJobRunPaginated(r.Context(), db.ListLogsByJobRunPaginatedParams{
		JobRunID: jobRunID,
		Limit:    int32(limit),
		Offset:   int32(offset),
	})
	if err != nil {
		respondJSON(w, 500, map[string]string{"error": "Failed to retrieve logs"})
		return
	}

	// Get total count for pagination
	totalCount := len(logs) // This is approximate since we're using LIMIT/OFFSET
	if len(logs) == limit {
		// There might be more logs, but we don't have exact count
		totalCount = offset + limit + 1
	} else {
		totalCount = offset + len(logs)
	}

	respondJSON(w, 200, map[string]interface{}{
		"logs": logs,
		"pagination": map[string]interface{}{
			"total":  totalCount,
			"limit":  limit,
			"offset": offset,
			"count":  len(logs),
		},
	})
}

// GetTaskRunLogs retrieves logs for a specific task run
func (h *LogsHandler) GetTaskRunLogs(w http.ResponseWriter, r *http.Request) {
	taskRunIDStr := r.PathValue("task_run_id")
	if taskRunIDStr == "" {
		respondJSON(w, 400, map[string]string{"error": "task_run_id is required"})
		return
	}

	taskRunID, err := utils.ParseUUID(taskRunIDStr)
	if err != nil {
		respondJSON(w, 400, map[string]string{"error": "Invalid task_run_id format"})
		return
	}

	logs, err := h.store.ListLogsByTaskRun(r.Context(), taskRunID)
	if err != nil {
		respondJSON(w, 500, map[string]string{"error": "Failed to retrieve logs"})
		return
	}

	respondJSON(w, 200, map[string]interface{}{
		"logs": logs,
	})
}

// GetSystemLogs retrieves system logs
func (h *LogsHandler) GetSystemLogs(w http.ResponseWriter, r *http.Request) {
	limit, offset := getPaginationParams(r)

	logs, err := h.store.ListSystemLogs(r.Context(), db.ListSystemLogsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		respondJSON(w, 500, map[string]string{"error": "Failed to retrieve system logs"})
		return
	}

	respondJSON(w, 200, map[string]interface{}{
		"logs": logs,
		"pagination": map[string]interface{}{
			"limit":  limit,
			"offset": offset,
			"count":  len(logs),
		},
	})
}

// GetLogsByType retrieves logs by type (system, job, task)
func (h *LogsHandler) GetLogsByType(w http.ResponseWriter, r *http.Request) {
	logType := r.PathValue("type")
	if logType == "" {
		respondJSON(w, 400, map[string]string{"error": "type is required"})
		return
	}

	// Validate log type
	if logType != "system" && logType != "job" && logType != "task" {
		respondJSON(w, 400, map[string]string{"error": "Invalid log type. Must be 'system', 'job', or 'task'"})
		return
	}

	limit, offset := getPaginationParams(r)

	logs, err := h.store.ListLogsByType(r.Context(), db.ListLogsByTypeParams{
		Type:   logType,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		respondJSON(w, 500, map[string]string{"error": "Failed to retrieve logs"})
		return
	}

	respondJSON(w, 200, map[string]interface{}{
		"logs": logs,
		"pagination": map[string]interface{}{
			"limit":  limit,
			"offset": offset,
			"count":  len(logs),
		},
	})
}

// DownloadJobRunLogs allows downloading the log file for a job run
func (h *LogsHandler) DownloadJobRunLogs(w http.ResponseWriter, r *http.Request) {
	jobRunIDStr := r.PathValue("job_run_id")
	if jobRunIDStr == "" {
		respondJSON(w, 400, map[string]string{"error": "job_run_id is required"})
		return
	}

	// Validate UUID format
	_, err := utils.ParseUUID(jobRunIDStr)
	if err != nil {
		respondJSON(w, 400, map[string]string{"error": "Invalid job_run_id format"})
		return
	}

	// Get the date parameter (optional, defaults to today)
	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		dateStr = r.URL.Query().Get("date")
	}
	if dateStr == "" {
		// Default to today's date
		dateStr = "2024-01-15" // You might want to use time.Now().Format("2006-01-02")
	}

	// Construct file path
	logFilePath := h.constructLogFilePath(jobRunIDStr, dateStr)

	// Set headers for file download
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+jobRunIDStr+"-"+dateStr+".txt\"")

	// Serve the file
	http.ServeFile(w, r, logFilePath)
}

// Helper functions

func (h *LogsHandler) constructLogFilePath(jobRunID, date string) string {
	// This should match the path used in the logger
	return "internal/storage/files/logs/" + jobRunID + "-" + date + ".txt"
}

func getPaginationParams(r *http.Request) (limit, offset int) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit = 100 // default limit
	offset = 0  // default offset

	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
			if limit > 1000 { // max limit
				limit = 1000
			}
		}
	}

	if offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	return limit, offset
}
