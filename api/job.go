package api

import (
	"encoding/json"
	"net/http"

	db "github.com/b0nbon1/temporal-lite/db/sqlc"
	"github.com/gin-gonic/gin"
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

func (server *Server) createJobRequest(ctx *gin.Context) {
	var req CreateJobRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateJobParams{
		Name:       req.Name,
		Schedule:   pgtype.Text{String: req.Schedule, Valid: true},
		Type:       pgtype.Text{String: req.Type, Valid: true},
		Config:     req.Config,
		Status:     db.NullJobStatus{JobStatus: db.JobStatusPending, Valid: true},
		Retries:    pgtype.Int4{Int32: int32(req.Retries), Valid: true},
		MaxRetries: pgtype.Int4{Int32:int32(req.MaxRetries), Valid: true},
	}

	job, err := server.store.CreateJob(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, job)
}