// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: task_runs.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createTaskRun = `-- name: CreateTaskRun :one
INSERT INTO task_runs (job_run_id, task_id, status)
VALUES ($1, $2, $3)
RETURNING id, job_run_id, task_id, status, started_at, finished_at, exit_code, output, error_message, created_at
`

type CreateTaskRunParams struct {
	JobRunID pgtype.UUID `json:"job_run_id"`
	TaskID   pgtype.UUID `json:"task_id"`
	Status   pgtype.Text `json:"status"`
}

type CreateTaskRunRow struct {
	ID           pgtype.UUID        `json:"id"`
	JobRunID     pgtype.UUID        `json:"job_run_id"`
	TaskID       pgtype.UUID        `json:"task_id"`
	Status       pgtype.Text        `json:"status"`
	StartedAt    pgtype.Timestamp   `json:"started_at"`
	FinishedAt   pgtype.Timestamp   `json:"finished_at"`
	ExitCode     pgtype.Int4        `json:"exit_code"`
	Output       pgtype.Text        `json:"output"`
	ErrorMessage pgtype.Text        `json:"error_message"`
	CreatedAt    pgtype.Timestamptz `json:"created_at"`
}

func (q *Queries) CreateTaskRun(ctx context.Context, arg CreateTaskRunParams) (CreateTaskRunRow, error) {
	row := q.db.QueryRow(ctx, createTaskRun, arg.JobRunID, arg.TaskID, arg.Status)
	var i CreateTaskRunRow
	err := row.Scan(
		&i.ID,
		&i.JobRunID,
		&i.TaskID,
		&i.Status,
		&i.StartedAt,
		&i.FinishedAt,
		&i.ExitCode,
		&i.Output,
		&i.ErrorMessage,
		&i.CreatedAt,
	)
	return i, err
}

const deleteTaskRun = `-- name: DeleteTaskRun :exec
DELETE FROM task_runs
WHERE id = $1
`

func (q *Queries) DeleteTaskRun(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteTaskRun, id)
	return err
}

const getTaskRun = `-- name: GetTaskRun :one
SELECT id, job_run_id, task_id, status, started_at, finished_at, exit_code, output, error_message, created_at
FROM task_runs
WHERE id = $1 LIMIT 1
`

type GetTaskRunRow struct {
	ID           pgtype.UUID        `json:"id"`
	JobRunID     pgtype.UUID        `json:"job_run_id"`
	TaskID       pgtype.UUID        `json:"task_id"`
	Status       pgtype.Text        `json:"status"`
	StartedAt    pgtype.Timestamp   `json:"started_at"`
	FinishedAt   pgtype.Timestamp   `json:"finished_at"`
	ExitCode     pgtype.Int4        `json:"exit_code"`
	Output       pgtype.Text        `json:"output"`
	ErrorMessage pgtype.Text        `json:"error_message"`
	CreatedAt    pgtype.Timestamptz `json:"created_at"`
}

func (q *Queries) GetTaskRun(ctx context.Context, id pgtype.UUID) (GetTaskRunRow, error) {
	row := q.db.QueryRow(ctx, getTaskRun, id)
	var i GetTaskRunRow
	err := row.Scan(
		&i.ID,
		&i.JobRunID,
		&i.TaskID,
		&i.Status,
		&i.StartedAt,
		&i.FinishedAt,
		&i.ExitCode,
		&i.Output,
		&i.ErrorMessage,
		&i.CreatedAt,
	)
	return i, err
}

const getTaskRunByJobRunAndTaskID = `-- name: GetTaskRunByJobRunAndTaskID :one
SELECT id, job_run_id, task_id, status, started_at, finished_at, exit_code, output, error_message, created_at
FROM task_runs
WHERE job_run_id = $1 AND task_id = $2 LIMIT 1
`

type GetTaskRunByJobRunAndTaskIDParams struct {
	JobRunID pgtype.UUID `json:"job_run_id"`
	TaskID   pgtype.UUID `json:"task_id"`
}

type GetTaskRunByJobRunAndTaskIDRow struct {
	ID           pgtype.UUID        `json:"id"`
	JobRunID     pgtype.UUID        `json:"job_run_id"`
	TaskID       pgtype.UUID        `json:"task_id"`
	Status       pgtype.Text        `json:"status"`
	StartedAt    pgtype.Timestamp   `json:"started_at"`
	FinishedAt   pgtype.Timestamp   `json:"finished_at"`
	ExitCode     pgtype.Int4        `json:"exit_code"`
	Output       pgtype.Text        `json:"output"`
	ErrorMessage pgtype.Text        `json:"error_message"`
	CreatedAt    pgtype.Timestamptz `json:"created_at"`
}

func (q *Queries) GetTaskRunByJobRunAndTaskID(ctx context.Context, arg GetTaskRunByJobRunAndTaskIDParams) (GetTaskRunByJobRunAndTaskIDRow, error) {
	row := q.db.QueryRow(ctx, getTaskRunByJobRunAndTaskID, arg.JobRunID, arg.TaskID)
	var i GetTaskRunByJobRunAndTaskIDRow
	err := row.Scan(
		&i.ID,
		&i.JobRunID,
		&i.TaskID,
		&i.Status,
		&i.StartedAt,
		&i.FinishedAt,
		&i.ExitCode,
		&i.Output,
		&i.ErrorMessage,
		&i.CreatedAt,
	)
	return i, err
}

const listTaskRuns = `-- name: ListTaskRuns :many
SELECT id, job_run_id, task_id, status, started_at, finished_at, exit_code, output, error_message, created_at
FROM task_runs
WHERE job_run_id = $1
ORDER BY created_at DESC
`

type ListTaskRunsRow struct {
	ID           pgtype.UUID        `json:"id"`
	JobRunID     pgtype.UUID        `json:"job_run_id"`
	TaskID       pgtype.UUID        `json:"task_id"`
	Status       pgtype.Text        `json:"status"`
	StartedAt    pgtype.Timestamp   `json:"started_at"`
	FinishedAt   pgtype.Timestamp   `json:"finished_at"`
	ExitCode     pgtype.Int4        `json:"exit_code"`
	Output       pgtype.Text        `json:"output"`
	ErrorMessage pgtype.Text        `json:"error_message"`
	CreatedAt    pgtype.Timestamptz `json:"created_at"`
}

func (q *Queries) ListTaskRuns(ctx context.Context, jobRunID pgtype.UUID) ([]ListTaskRunsRow, error) {
	rows, err := q.db.Query(ctx, listTaskRuns, jobRunID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListTaskRunsRow{}
	for rows.Next() {
		var i ListTaskRunsRow
		if err := rows.Scan(
			&i.ID,
			&i.JobRunID,
			&i.TaskID,
			&i.Status,
			&i.StartedAt,
			&i.FinishedAt,
			&i.ExitCode,
			&i.Output,
			&i.ErrorMessage,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listTaskRunsByJob = `-- name: ListTaskRunsByJob :many
SELECT tr.id, tr.job_run_id, tr.task_id, tr.status, tr.started_at, tr.finished_at, tr.exit_code, tr.output, tr.error_message, tr.created_at
FROM task_runs tr
JOIN jobs j ON tr.job_run_id = j.id
WHERE j.id = $1
ORDER BY tr.created_at DESC
`

type ListTaskRunsByJobRow struct {
	ID           pgtype.UUID        `json:"id"`
	JobRunID     pgtype.UUID        `json:"job_run_id"`
	TaskID       pgtype.UUID        `json:"task_id"`
	Status       pgtype.Text        `json:"status"`
	StartedAt    pgtype.Timestamp   `json:"started_at"`
	FinishedAt   pgtype.Timestamp   `json:"finished_at"`
	ExitCode     pgtype.Int4        `json:"exit_code"`
	Output       pgtype.Text        `json:"output"`
	ErrorMessage pgtype.Text        `json:"error_message"`
	CreatedAt    pgtype.Timestamptz `json:"created_at"`
}

func (q *Queries) ListTaskRunsByJob(ctx context.Context, id pgtype.UUID) ([]ListTaskRunsByJobRow, error) {
	rows, err := q.db.Query(ctx, listTaskRunsByJob, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListTaskRunsByJobRow{}
	for rows.Next() {
		var i ListTaskRunsByJobRow
		if err := rows.Scan(
			&i.ID,
			&i.JobRunID,
			&i.TaskID,
			&i.Status,
			&i.StartedAt,
			&i.FinishedAt,
			&i.ExitCode,
			&i.Output,
			&i.ErrorMessage,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateTaskRun = `-- name: UpdateTaskRun :exec
UPDATE task_runs
SET status = $2, started_at = $3, finished_at = $4, exit_code = $5, output = $6, error_message = $7, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
`

type UpdateTaskRunParams struct {
	ID           pgtype.UUID      `json:"id"`
	Status       pgtype.Text      `json:"status"`
	StartedAt    pgtype.Timestamp `json:"started_at"`
	FinishedAt   pgtype.Timestamp `json:"finished_at"`
	ExitCode     pgtype.Int4      `json:"exit_code"`
	Output       pgtype.Text      `json:"output"`
	ErrorMessage pgtype.Text      `json:"error_message"`
}

func (q *Queries) UpdateTaskRun(ctx context.Context, arg UpdateTaskRunParams) error {
	_, err := q.db.Exec(ctx, updateTaskRun,
		arg.ID,
		arg.Status,
		arg.StartedAt,
		arg.FinishedAt,
		arg.ExitCode,
		arg.Output,
		arg.ErrorMessage,
	)
	return err
}

const updateTaskRunError = `-- name: UpdateTaskRunError :exec
UPDATE task_runs
SET error_message = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
`

type UpdateTaskRunErrorParams struct {
	ID           pgtype.UUID `json:"id"`
	ErrorMessage pgtype.Text `json:"error_message"`
}

func (q *Queries) UpdateTaskRunError(ctx context.Context, arg UpdateTaskRunErrorParams) error {
	_, err := q.db.Exec(ctx, updateTaskRunError, arg.ID, arg.ErrorMessage)
	return err
}

const updateTaskRunOutput = `-- name: UpdateTaskRunOutput :exec
UPDATE task_runs
SET output = $2, updated_at = CURRENT_TIMESTAMP, finished_at = CURRENT_TIMESTAMP, status = 'completed'
WHERE id = $1
`

type UpdateTaskRunOutputParams struct {
	ID     pgtype.UUID `json:"id"`
	Output pgtype.Text `json:"output"`
}

func (q *Queries) UpdateTaskRunOutput(ctx context.Context, arg UpdateTaskRunOutputParams) error {
	_, err := q.db.Exec(ctx, updateTaskRunOutput, arg.ID, arg.Output)
	return err
}

const updateTaskRunStatus = `-- name: UpdateTaskRunStatus :exec
UPDATE task_runs
SET status = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
`

type UpdateTaskRunStatusParams struct {
	ID     pgtype.UUID `json:"id"`
	Status pgtype.Text `json:"status"`
}

func (q *Queries) UpdateTaskRunStatus(ctx context.Context, arg UpdateTaskRunStatusParams) error {
	_, err := q.db.Exec(ctx, updateTaskRunStatus, arg.ID, arg.Status)
	return err
}
