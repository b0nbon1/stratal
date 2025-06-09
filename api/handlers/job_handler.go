package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/b0nbon1/stratal/db/dto"
	db "github.com/b0nbon1/stratal/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/b0nbon1/stratal/pkg/utils"
)

type JobHandler struct {
	store *db.Queries
}

type CreateJobRequest struct {
	Name       string          `json:"name" binding:"required"`
	Schedule   string          `json:"schedule" binding:"required"`
	Type       string          `json:"type" binding:"required"`
	Config     json.RawMessage `json:"config" binding:"required"`
	Status     string          `json:"status" binding:"omitempty,oneof=pending running success failed"`
	Retries    int             `json:"retries" binding:"omitempty,gte=0"`
	MaxRetries int             `json:"max_retries" binding:"required,gte=0"`
}

type ListJobsParams struct {
    Limit  int32 `form:"limit" json:"limit"`
    Offset int32 `form:"offset" json:"offset"`
}

type JobIDRequestBind struct {
	ID pgtype.UUID `uri:"id" binding:"required"`
}

func NewJobHandler(store *db.Queries) *JobHandler {
	return &JobHandler{store: store}
}

func (h *JobHandler) CreateJob(ctx *gin.Context) {
	var req CreateJobRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
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

	job, err := h.store.CreateJob(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, job)
}

func (h *JobHandler) GetJobRequest(ctx *gin.Context) {
	var req JobIDRequestBind
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	job, err := h.store.GetJob(ctx, req.ID)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, utils.ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, job)
}

func (h *JobHandler) ListJobs(ctx *gin.Context) {
	var req ListJobsParams
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}
	fmt.Println("ListJobsParams:", req)
	dbParams := db.ListJobsParams{
		Limit:  req.Limit,
		Offset: req.Offset,
	}
	jobs, err := h.store.ListJobs(ctx, dbParams)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, jobs)
}
