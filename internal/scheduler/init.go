package scheduler

import (
	"context"
	"encoding/json"

	"log"

	db "github.com/b0nbon1/stratal/db/sqlc"
	"github.com/b0nbon1/stratal/internal/queue"
	"github.com/robfig/cron/v3"
)

func StartScheduler(q queue.TaskQueue, store db.Queries) {
	c := cron.New()
	c.AddFunc("@every 10s", func() {
		var jobs []db.Job
		jobs, err := store.ListPendingJobs(context.Background())
		if err != nil {
			log.Println("Error listing pending jobs:", err)
			return
		}

		for _, job := range jobs {
			jobJSON, err := json.Marshal(job)
			if err != nil {
				log.Println(err)
			}
			q.Enqueue(jobJSON)
		}
	})
	c.Start()
}
