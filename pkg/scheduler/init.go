package scheduler

import (
	"context"
	"encoding/json"
	"fmt"

	// "fmt"
	"log"

	db "github.com/b0nbon1/temporal-lite/db/sqlc"
	"github.com/b0nbon1/temporal-lite/pkg/queue"
	"github.com/robfig/cron/v3"
)

func StartScheduler(q queue.TaskQueue, store db.Queries) {
	c := cron.New()
	c.AddFunc("@every 10s", func() {
		var jobs []db.Job
		jobs, err := store.ListPendingJobs(context.Background())
		fmt.Println("Listing pending jobs")
		if err != nil {
			log.Println("Error listing pending jobs:", err)
			return
		}
		fmt.Println("Pending jobs:", len(jobs))
		for _, job := range jobs {
			jobJSON, err := json.Marshal(job)
			if err != nil {
				panic(err)
			}

			q.Enqueue(jobJSON)
		}
	})
	c.Start()
}