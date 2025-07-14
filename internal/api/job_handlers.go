package api

import (
	"net/http"
	"strconv"
	"strings"

	db "github.com/b0nbon1/stratal/internal/storage/db/sqlc"
	"github.com/b0nbon1/stratal/pkg/utils"
)

// GetJob retrieves a specific job by ID
func (hs *HTTPServer) GetJob(w http.ResponseWriter, r *http.Request) {
	// Extract job ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 {
		respondError(w, 400, "Invalid job ID")
		return
	}

	jobID := pathParts[len(pathParts)-1]
	jobUUID, err := utils.ParseUUID(jobID)
	if err != nil {
		respondError(w, 400, "Invalid job UUID", err.Error())
		return
	}

	job, err := hs.store.GetJobWithTasks(hs.ctx, jobUUID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			respondError(w, 404, "Job not found")
		} else {
			respondError(w, 500, "Failed to fetch job", err.Error())
		}
		return
	}

	respondJSON(w, 200, job)
}

// ListJobs retrieves a list of jobs with pagination
func (hs *HTTPServer) ListJobs(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
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

// GetJobRun retrieves a specific job run by ID
func (hs *HTTPServer) GetJobRun(w http.ResponseWriter, r *http.Request) {
	// Extract job run ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 {
		respondError(w, 400, "Invalid job run ID")
		return
	}

	jobRunID := pathParts[len(pathParts)-1]
	jobRunUUID, err := utils.ParseUUID(jobRunID)
	if err != nil {
		respondError(w, 400, "Invalid job run UUID", err.Error())
		return
	}

	jobRun, err := hs.store.JobRunsWithTasks(hs.ctx, jobRunUUID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			respondError(w, 404, "Job run not found")
		} else {
			respondError(w, 500, "Failed to fetch job run", err.Error())
		}
		return
	}

	respondJSON(w, 200, jobRun)
}

// ListJobRuns retrieves all runs for a specific job
func (hs *HTTPServer) ListJobRuns(w http.ResponseWriter, r *http.Request) {
	// Extract job ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 6 {
		respondError(w, 400, "Invalid job ID")
		return
	}

	jobID := pathParts[len(pathParts)-2]
	jobUUID, err := utils.ParseUUID(jobID)
	if err != nil {
		respondError(w, 400, "Invalid job UUID", err.Error())
		return
	}

	jobRuns, err := hs.store.ListJobRuns(hs.ctx, jobUUID)
	if err != nil {
		respondError(w, 500, "Failed to list job runs", err.Error())
		return
	}

	respondJSON(w, 200, map[string]interface{}{
		"job_id":   jobID,
		"job_runs": jobRuns,
		"count":    len(jobRuns),
	})
}

// ListTaskRunsForJobRun retrieves all task runs for a specific job run
func (hs *HTTPServer) ListTaskRunsForJobRun(w http.ResponseWriter, r *http.Request) {
	// Extract job run ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 6 {
		respondError(w, 400, "Invalid job run ID")
		return
	}

	jobRunID := pathParts[len(pathParts)-2]
	jobRunUUID, err := utils.ParseUUID(jobRunID)
	if err != nil {
		respondError(w, 400, "Invalid job run UUID", err.Error())
		return
	}

	taskRuns, err := hs.store.ListTaskRuns(hs.ctx, jobRunUUID)
	if err != nil {
		respondError(w, 500, "Failed to list task runs", err.Error())
		return
	}

	// Get the job run details too
	jobRun, _ := hs.store.GetJobRun(hs.ctx, jobRunUUID)

	respondJSON(w, 200, map[string]interface{}{
		"job_run_id": jobRunID,
		"job_run":    jobRun,
		"task_runs":  taskRuns,
		"count":      len(taskRuns),
	})
}
