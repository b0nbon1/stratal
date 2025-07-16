package api

import (
	"net/http"
	"strconv"

	db "github.com/b0nbon1/stratal/internal/storage/db/sqlc"
	"github.com/b0nbon1/stratal/pkg/router"
	"github.com/b0nbon1/stratal/pkg/utils"
)

// GetJob retrieves a specific job by ID from URL path parameter
func (hs *HTTPServer) GetJob(w http.ResponseWriter, r *http.Request) {
	// Extract job ID from URL path parameter (e.g., /jobs/:id)
	jobID := router.GetParam(r, "id")
	if jobID == "" {
		respondError(w, 400, "Job ID is required in URL path")
		return
	}

	jobUUID, err := utils.ParseUUID(jobID)
	if err != nil {
		respondError(w, 400, "Invalid job UUID", err.Error())
		return
	}

	job, err := hs.store.GetJobWithTasks(hs.ctx, jobUUID)
	if err != nil {
		if utils.ContainsSubstring(err.Error(), "no rows") {
			respondError(w, 404, "Job not found")
		} else {
			respondError(w, 500, "Failed to fetch job", err.Error())
		}
		return
	}

	respondJSON(w, 200, job)
}

// ListJobs retrieves a list of jobs with pagination or a specific job by ID
func (hs *HTTPServer) ListJobs(w http.ResponseWriter, r *http.Request) {
	// Check if this is a request for a specific job
	if jobID := r.URL.Query().Get("id"); jobID != "" {
		hs.GetJob(w, r)
		return
	}

	// Parse pagination parameters
	limit := 20
	offset := 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	jobs, err := hs.store.ListJobs(hs.ctx, db.ListJobsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		respondError(w, 500, "Failed to list jobs", err.Error())
		return
	}

	respondJSON(w, 200, map[string]interface{}{
		"jobs":   jobs,
		"limit":  limit,
		"offset": offset,
		"count":  len(jobs),
	})
}

// GetJobRun retrieves a specific job run by ID from query parameter
func (hs *HTTPServer) GetJobRun(w http.ResponseWriter, r *http.Request) {
	jobRunID := r.URL.Query().Get("id")
	if jobRunID == "" {
		respondError(w, 400, "Job run ID is required as query parameter")
		return
	}

	jobRunUUID, err := utils.ParseUUID(jobRunID)
	if err != nil {
		respondError(w, 400, "Invalid job run UUID", err.Error())
		return
	}

	jobRun, err := hs.store.JobRunsWithTasks(hs.ctx, jobRunUUID)
	if err != nil {
		if utils.ContainsSubstring(err.Error(), "no rows") {
			respondError(w, 404, "Job run not found")
		} else {
			respondError(w, 500, "Failed to fetch job run", err.Error())
		}
		return
	}

	respondJSON(w, 200, jobRun)
}

