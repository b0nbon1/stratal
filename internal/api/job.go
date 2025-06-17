package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/b0nbon1/stratal/internal/storage/db/dto"
	db "github.com/b0nbon1/stratal/internal/storage/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type TaskJobBody struct {
	Name   string        
	Type   string         
	Config dto.TaskConfig 
	Order  int32          
}

type JobBodyParams struct {
	Name        string 
	Description string
	Source      string
	RawPayload  []byte
	Tasks       []TaskJobBody
}

func (hs *HTTPServer) CreateJob(w http.ResponseWriter, r *http.Request) {
	var reqBodyJob JobBodyParams
	err := parseJSON(r, &reqBodyJob)
	if err != nil {
		respondJSON(w, 500, nil)
		return
	}
	ctx := context.Background()

	// create the job, which will be a blue_print
	job, err := hs.store.CreateJob(ctx, db.CreateJobParams{
		Name: reqBodyJob.Name,
		Description: pgtype.Text{String: reqBodyJob.Description, Valid: true},
		Source: reqBodyJob.Source,
		RawPayload: reqBodyJob.RawPayload,
	})

	if err != nil {
		respondJSON(w, 500, nil)
		return
	}

	fmt.Println(job)
}

