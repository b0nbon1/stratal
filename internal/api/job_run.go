package api

import (
	"fmt"
	"log"
	"net/http"

	db "github.com/b0nbon1/stratal/internal/storage/db/sqlc"
	"github.com/b0nbon1/stratal/pkg/utils"
)

type JobRunBody struct {
	JobID       string `json:"job_id"`
	TriggeredBy string `json:"triggered_by"`
}

func (hs *HTTPServer) CreateJobRun(w http.ResponseWriter, r *http.Request) {
	var reqBodyJobRun JobRunBody
	err := parseJSON(r, &reqBodyJobRun)
	if err != nil {
		respondJSON(w, 500, fmt.Errorf("unable to parse body %w", err))
		return
	}

	parsedJobID, err := utils.ParseUUID(reqBodyJobRun.JobID)
	if err != nil {
		respondJSON(w, 400, fmt.Errorf("invalid Job UUID %w", err))
		return
	}

	// execute the transaction to create the jobs and also tasks, if fails rollback everthing
	data, err := hs.store.CreateJobRunTx(hs.ctx, parsedJobID, reqBodyJobRun.TriggeredBy)
	if err != nil {
		respondJSON(w, 500, fmt.Errorf("job Run failed, %w", err))
		return
	}

	err = hs.queue.Enqueue(data.JobRunId)
	if err != nil {
		respondJSON(w, 500, fmt.Errorf("job Run failed to be queued: %w", err))
		return
	}

	parsedJobRunID, _ := utils.ParseUUID(reqBodyJobRun.JobID)

	err = hs.store.UpdateJobRunStatus(hs.ctx, db.UpdateJobRunStatusParams{
		ID:     parsedJobRunID,
		Status: utils.ParseText("queued"),
	})
	if err != nil {
		log.Println("unable to update JobRun Status")
	}

	respondJSON(w, 200, map[string]interface{}{
		"message":  "Created Job_run successful, wait for it to be queued",
		"JobRunID": reqBodyJobRun.JobID,
	})

}
