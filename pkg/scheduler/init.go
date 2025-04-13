package scheduler

import (
	"context"
	"log"

	db "github.com/b0nbon1/temporal-lite/db/sqlc"
	"github.com/b0nbon1/temporal-lite/pkg/queue"
	"github.com/robfig/cron/v3"
)

func StartScheduler(q queue.TaskQueue, store db.Querier) {
	c := cron.New()
	c.AddFunc("@every 10s", func() {
		var jobs []db.Job
		jobs, err := store.ListPendingJobs(context.Background())
		if err != nil {
			log.Println("Error listing pending jobs:", err)
			return
		}
		for _, job := range jobs {
			q.Enqueue(job)
		}
	})
	c.Start()
}