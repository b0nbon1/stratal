package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/b0nbon1/stratal/internal/storage/db/dto"
	db "github.com/b0nbon1/stratal/internal/storage/db/sqlc"
	"github.com/b0nbon1/stratal/pkg/router"
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
	RunImmediately bool          `json:"run_immediately"`
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
		RawPayload:  []byte(reqBodyJob.RawPayload),
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
		jobRunData, err := hs.store.CreateJobRunTx(hs.ctx, data.Job.ID, "api")
		if err != nil {
			response["warning"] = fmt.Sprintf("Job created but failed to create job run: %v", err)
		} else {
			err = hs.queue.Enqueue(jobRunData.JobRunId)
			if err != nil {
				response["warning"] = fmt.Sprintf("Job and job run created but failed to queue: %v", err)
			} else {
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

func (hs *HTTPServer) GetJob(w http.ResponseWriter, r *http.Request) {
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

func (hs *HTTPServer) ListJobs(w http.ResponseWriter, r *http.Request) {
	if jobID := r.URL.Query().Get("id"); jobID != "" {
		hs.GetJob(w, r)
		return
	}

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
