package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/b0nbon1/stratal/internal/logger"
	db "github.com/b0nbon1/stratal/internal/storage/db/sqlc"
	"github.com/b0nbon1/stratal/pkg/router"
	"github.com/b0nbon1/stratal/pkg/utils"
)

type LogsHandler struct {
	store    *db.SQLStore
	streamer *logger.LogStreamer
}

func NewLogsHandler(store *db.SQLStore, streamer *logger.LogStreamer) *LogsHandler {
	return &LogsHandler{
		store:    store,
		streamer: streamer,
	}
}

func (h *LogsHandler) RegisterLogRoutes(v1 *router.Router) {
	v1.Get("/logs/stream/ws", h.HandleWebSocketStream)
	v1.Get("/logs/stream/sse", h.HandleSSEStream)
	v1.Get("/logs/stream/status", h.GetStreamingStatus)

	v1.Get("/logs/job-runs/:job_run_id", h.GetJobRunLogs)
	v1.Get("/logs/task-runs/:task_run_id", h.GetTaskRunLogs)
	v1.Get("/logs/job-runs/:job_run_id/download", h.DownloadJobRunLogs)
	v1.Get("/logs/system", h.GetSystemLogs)
	v1.Get("/logs/type/:type", h.GetLogsByType)
}

func (h *LogsHandler) HandleWebSocketStream(w http.ResponseWriter, r *http.Request) {
	h.streamer.HandleWebSocket(w, r)
}

func (h *LogsHandler) HandleSSEStream(w http.ResponseWriter, r *http.Request) {
	h.streamer.HandleSSE(w, r)
}

func (h *LogsHandler) GetStreamingStatus(w http.ResponseWriter, r *http.Request) {
	status := h.streamer.GetActiveConnections()
	respondJSON(w, 200, map[string]interface{}{
		"status":      "active",
		"connections": status,
	})
}

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

	totalCount := len(logs)
	if len(logs) == limit {
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

func (h *LogsHandler) GetLogsByType(w http.ResponseWriter, r *http.Request) {
	logType := r.PathValue("type")
	if logType == "" {
		respondJSON(w, 400, map[string]string{"error": "type is required"})
		return
	}

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

func (h *LogsHandler) DownloadJobRunLogs(w http.ResponseWriter, r *http.Request) {
	jobRunIDStr := r.PathValue("job_run_id")
	if jobRunIDStr == "" {
		respondJSON(w, 400, map[string]string{"error": "job_run_id is required"})
		return
	}

	_, err := utils.ParseUUID(jobRunIDStr)
	if err != nil {
		respondJSON(w, 400, map[string]string{"error": "Invalid job_run_id format"})
		return
	}

	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		dateStr = r.URL.Query().Get("date")
	}
	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
	}

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

	limit = 100
	offset = 0

	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
			if limit > 1000 {
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
