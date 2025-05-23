package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/b0nbon1/stratal/db/dto"
	db "github.com/b0nbon1/stratal/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateJobRequest struct {
	Name       string          `json:"name" binding:"required"`
	Schedule   string          `json:"schedule" binding:"required"`
	Type       string          `json:"type" binding:"required"`
	Config     json.RawMessage `json:"config" binding:"required"`
	Status     string          `json:"status" binding:"omitempty,oneof=pending running success failed"`
	Retries    int             `json:"retries" binding:"omitempty,gte=0"`
	MaxRetries int             `json:"max_retries" binding:"required,gte=0"`
}

type JobIDRequestBind struct {
	ID pgtype.UUID `uri:"id" binding:"required"`
}

func (server *Server) createJobRequest(ctx *gin.Context) {
	var req CreateJobRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var config dto.AutomationConfig
	err := json.Unmarshal(req.Config, &config)
	if err != nil {
		log.Printf("failed to decode config: %v", err)
		return
	}

	arg := db.CreateJobParams{
		Name:       req.Name,
		Schedule:   pgtype.Text{String: req.Schedule, Valid: true},
		Type:       pgtype.Text{String: req.Type, Valid: true},
		Config:     config,
		Status:     db.NullJobStatus{JobStatus: db.JobStatusPending, Valid: true},
		Retries:    pgtype.Int4{Int32: int32(req.Retries), Valid: true},
		MaxRetries: pgtype.Int4{Int32: int32(req.MaxRetries), Valid: true},
	}

	job, err := server.store.CreateJob(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, job)
}

func (server *Server) getJobRequest(ctx *gin.Context) {
	var req JobIDRequestBind
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	job, err := server.store.GetJob(ctx, req.ID)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, job)
}

func (server *Server) listJobsRequest(ctx *gin.Context) {
	var req db.ListJobsParams
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	jobs, err := server.store.ListJobs(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, jobs)
}
