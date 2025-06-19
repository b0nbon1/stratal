package api

import (
	"fmt"
	"net/http"

	"github.com/b0nbon1/stratal/internal/storage/db/dto"
	db "github.com/b0nbon1/stratal/internal/storage/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type TaskJobBody struct {
	Name   string   `json:"name"`     
	Type   string       `json:"type"`  
	Config dto.TaskConfig `json:"config"` 
	Order  int32      `json:"order"`     
}

type JobBodyParams struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Source      string `json:"source"`
	RawPayload  []byte `json:"raw_payload"`
	Tasks       []TaskJobBody `json:"tasks"`
}

func (hs *HTTPServer) CreateJob(w http.ResponseWriter, r *http.Request) {
	var reqBodyJob JobBodyParams
	err := parseJSON(r, &reqBodyJob)
	if err != nil {
		respondJSON(w, 500, nil)
		return
	}

	// create the job params, which will be a blue_print
	jobParam := db.CreateJobParams{
		Name: reqBodyJob.Name,
		Description: pgtype.Text{String: reqBodyJob.Description, Valid: true},
		Source: reqBodyJob.Source,
		RawPayload: reqBodyJob.RawPayload,
	}
	
	// create Tasks that will be linke to the Job
	var taskParams []db.CreateTaskParams
	for _, task := range reqBodyJob.Tasks {
		taskParams = append(taskParams, db.CreateTaskParams{
			Name:   task.Name,
			Type:   task.Type,
			Order:  task.Order,
			Config: task.Config,
		})
	}

	// execute the transaction to create the jobs and also tasks, if fails rollback everthing
	data, err := hs.store.CreateJobWithTasksTx(hs.ctx, jobParam, taskParams)
	if err != nil {
		respondJSON(w, 500, map[string]interface{}{
			"message": fmt.Errorf("create job with tasks failed: %w", err),
		})
		return
	}

	// create a Job_run


	respondJSON(w, 201, data)	
}

