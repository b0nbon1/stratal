package api

import (
	"fmt"
	"net/http"

	"github.com/b0nbon1/stratal/pkg/router"
	"github.com/b0nbon1/stratal/pkg/utils"
)

// PauseJobRun pauses a running or queued job run
func (hs *HTTPServer) PauseJobRun(w http.ResponseWriter, r *http.Request) {
	jobRunID := router.GetParam(r, "id")

	if jobRunID == "" {
		respondError(w, 400, "Job run ID is required")
		return
	}

	jobRunUUID, err := utils.ParseUUID(jobRunID)
	if err != nil {
		respondError(w, 400, "Invalid job run UUID", err.Error())
		return
	}

	// Check if job run exists and is in a pausable state
	jobRun, err := hs.store.GetJobRunWithPauseInfo(hs.ctx, jobRunUUID)
	if err != nil {
		if utils.ContainsSubstring(err.Error(), "no rows") {
			respondError(w, 404, "Job run not found")
		} else {
			respondError(w, 500, "Failed to fetch job run", err.Error())
		}
		return
	}

	// Check if job run is in a pausable state
	status := jobRun.Status.String
	if status != "running" && status != "queued" {
		respondError(w, 400, fmt.Sprintf("Cannot pause job run in '%s' status. Only running or queued jobs can be paused.", status))
		return
	}

	// Pause the job run
	err = hs.store.PauseJobRun(hs.ctx, jobRunUUID)
	if err != nil {
		respondError(w, 500, "Failed to pause job run", err.Error())
		return
	}

	respondJSON(w, 200, map[string]interface{}{
		"message":         "Job run paused successfully",
		"job_run_id":      jobRunID,
		"job_name":        jobRun.JobName,
		"previous_status": status,
		"new_status":      "paused",
	})
}

// ResumeJobRun resumes a paused job run
func (hs *HTTPServer) ResumeJobRun(w http.ResponseWriter, r *http.Request) {
	jobRunID := router.GetParam(r, "id")

	if jobRunID == "" {
		respondError(w, 400, "Job run ID is required")
		return
	}

	jobRunUUID, err := utils.ParseUUID(jobRunID)
	if err != nil {
		respondError(w, 400, "Invalid job run UUID", err.Error())
		return
	}

	// Check if job run exists and is paused
	jobRun, err := hs.store.GetJobRunWithPauseInfo(hs.ctx, jobRunUUID)
	if err != nil {
		if utils.ContainsSubstring(err.Error(), "no rows") {
			respondError(w, 404, "Job run not found")
		} else {
			respondError(w, 500, "Failed to fetch job run", err.Error())
		}
		return
	}

	// Check if job run is paused
	if jobRun.Status.String != "paused" {
		respondError(w, 400, fmt.Sprintf("Cannot resume job run in '%s' status. Only paused jobs can be resumed.", jobRun.Status.String))
		return
	}

	// Resume the job run
	err = hs.store.ResumeJobRun(hs.ctx, jobRunUUID)
	if err != nil {
		respondError(w, 500, "Failed to resume job run", err.Error())
		return
	}

	// If job was previously running, we need to re-queue it
	var newStatus string
	if jobRun.StartedAt.Valid {
		// Job was running when paused, queue it to continue
		err = hs.queue.Enqueue(jobRunID)
		if err != nil {
			respondError(w, 500, "Failed to re-queue resumed job run", err.Error())
			return
		}
		newStatus = "queued"
	} else {
		// Job was queued when paused, just change status back to queued
		newStatus = "queued"
	}

	respondJSON(w, 200, map[string]interface{}{
		"message":         "Job run resumed successfully",
		"job_run_id":      jobRunID,
		"job_name":        jobRun.JobName,
		"previous_status": "paused",
		"new_status":      newStatus,
	})
}

// GetPausedJobRuns returns all paused job runs
func (hs *HTTPServer) GetPausedJobRuns(w http.ResponseWriter, r *http.Request) {
	pausedJobRuns, err := hs.store.GetPausedJobRuns(hs.ctx)
	if err != nil {
		respondError(w, 500, "Failed to fetch paused job runs", err.Error())
		return
	}

	respondJSON(w, 200, map[string]interface{}{
		"paused_job_runs": pausedJobRuns,
		"count":           len(pausedJobRuns),
	})
}
