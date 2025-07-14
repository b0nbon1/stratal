package api

import (
	"fmt"
	"net/http"

	"github.com/b0nbon1/stratal/internal/storage/db/dto"
	db "github.com/b0nbon1/stratal/internal/storage/db/sqlc"
	"github.com/b0nbon1/stratal/pkg/utils"
	"github.com/jackc/pgx/v5/pgtype"
)

type TaskJobBody struct {
	Name   string         `json:"name"`
	Type   string         `json:"type"`
	Config dto.TaskConfig `json:"config"`
	Order  int32          `json:"order"`
}

type JobBodyParams struct {
	Name           string        `json:"name"`
	Description    string        `json:"description"`
	Source         string        `json:"source"`
	RawPayload     []byte        `json:"raw_payload"`
	Tasks          []TaskJobBody `json:"tasks"`
	RunImmediately bool          `json:"run_immediately"` // Optional: create and queue a job run
}

func (hs *HTTPServer) CreateJob(w http.ResponseWriter, r *http.Request) {
	var reqBodyJob JobBodyParams
	err := parseJSON(r, &reqBodyJob)
	if err != nil {
		respondJSON(w, 400, map[string]interface{}{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Validate required fields
	if reqBodyJob.Name == "" {
		respondJSON(w, 400, map[string]interface{}{
			"error": "Job name is required",
		})
		return
	}

	if reqBodyJob.Source == "" {
		reqBodyJob.Source = "api" // Default source
	}

	// create the job params, which will be a blue_print
	jobParam := db.CreateJobParams{
		Name:        reqBodyJob.Name,
		Description: pgtype.Text{String: reqBodyJob.Description, Valid: true},
		Source:      reqBodyJob.Source,
		RawPayload:  reqBodyJob.RawPayload,
	}

	// create Tasks that will be linked to the Job
	var taskParams []db.CreateTaskParams
	for _, task := range reqBodyJob.Tasks {
		taskParams = append(taskParams, db.CreateTaskParams{
			Name:   task.Name,
			Type:   task.Type,
			Order:  task.Order,
			Config: task.Config,
		})
	}

	// execute the transaction to create the jobs and also tasks, if fails rollback everything
	data, err := hs.store.CreateJobWithTasksTx(hs.ctx, jobParam, taskParams)
	if err != nil {
		respondJSON(w, 500, map[string]interface{}{
			"error":   "Failed to create job with tasks",
			"details": err.Error(),
		})
		return
	}

	response := map[string]interface{}{
		"job":     data.Job,
		"tasks":   data.Tasks,
		"message": "Job created successfully",
	}

	// If requested, create and queue a job run immediately
	if reqBodyJob.RunImmediately && data.Job != nil {
		// Create job run
		jobRunData, err := hs.store.CreateJobRunTx(hs.ctx, data.Job.ID, "api")
		if err != nil {
			// Don't fail the whole request, just add warning
			response["warning"] = fmt.Sprintf("Job created but failed to create job run: %v", err)
		} else {
			// Queue the job run
			err = hs.queue.Enqueue(jobRunData.JobRunId)
			if err != nil {
				response["warning"] = fmt.Sprintf("Job and job run created but failed to queue: %v", err)
			} else {
				// Update job run status to queued
				jobRunID, _ := utils.ParseUUID(jobRunData.JobRunId)
				err = hs.store.UpdateJobRunStatus(hs.ctx, db.UpdateJobRunStatusParams{
					ID:     jobRunID,
					Status: utils.ParseText("queued"),
				})
				if err != nil {
					fmt.Printf("Failed to update job run status to queued: %v\n", err)
				}

				response["job_run_id"] = jobRunData.JobRunId
				response["status"] = "Job created and queued for execution"
			}
		}
	}

	respondJSON(w, 201, response)
}
