package scheduler

import (
	"context"
	"fmt"

	"log"

	"github.com/b0nbon1/stratal/internal/queue"
	db "github.com/b0nbon1/stratal/internal/storage/db/sqlc"
	"github.com/robfig/cron/v3"
)

func FetchCronJobsScheduled(q queue.TaskQueue, store *db.SQLStore, ctx context.Context) *cron.Cron {
	c := cron.New()
	c.AddFunc("@every 1m", func() {
		log.Println("Running scheduled job to enqueue scheduled jobs...")
		rows, err := store.ListPendingJobRuns(ctx)
		if err != nil {
			log.Println("Error listing pending jobs:", err)
			return
		}
		fmt.Println("Found", len(rows), "scheduled jobs")

		for _, job := range rows {
			err = q.Enqueue(job.ID.String())
			if err != nil {
				log.Println("Error queueing scheduled jobs:", err)
			}
		}
	})
	c.Start()

	return c
}
